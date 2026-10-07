package main

import (
	"bytes"
	"fmt"
	"go/ast"
	"go/format"
	"go/parser"
	"go/token"
	"os"
	"sort"
	"strconv"
	"strings"

	"terraform-provider-verity/internal/spec"
)

type bulkOptions struct {
	Registry   string
	OpenAPIDir string
	Output     string
	Check      bool
}

type bulkSDK struct {
	services map[string]string
	methods  map[string]*ast.FuncDecl
	structs  map[string]goStruct
}

type bulkOperation struct {
	Service     string
	Method      string
	Request     string
	BodyField   string
	ValueType   string
	BodyType    string
	Setter      string
	SplitSetter string
}

type bulkBinding struct {
	Key        string
	SplitKey   string
	Variants   []spec.ResourceSpec
	Operations map[string]bulkOperation
}

func discoverBulkSDK(dir string) (bulkSDK, error) {
	packages, err := parser.ParseDir(token.NewFileSet(), dir, nil, 0)
	if err != nil {
		return bulkSDK{}, err
	}
	sdk := bulkSDK{services: map[string]string{}, methods: map[string]*ast.FuncDecl{}, structs: map[string]goStruct{}}
	for _, pkg := range packages {
		for _, file := range pkg.Files {
			for _, decl := range file.Decls {
				if method, ok := decl.(*ast.FuncDecl); ok && method.Recv != nil {
					receiver := strings.TrimPrefix(typeString(method.Recv.List[0].Type), "*")
					sdk.methods[receiver+"."+method.Name.Name] = method
				}
				declaration, ok := decl.(*ast.GenDecl)
				if !ok || declaration.Tok != token.TYPE {
					continue
				}
				for _, definition := range declaration.Specs {
					typeSpec := definition.(*ast.TypeSpec)
					structure, ok := typeSpec.Type.(*ast.StructType)
					if !ok {
						continue
					}
					sdk.structs[typeSpec.Name.Name] = goStruct{Name: typeSpec.Name.Name, Fields: structFields(structure)}
					if typeSpec.Name.Name == "APIClient" {
						for _, field := range structure.Fields.List {
							service := strings.TrimPrefix(typeString(field.Type), "*")
							if len(field.Names) == 1 && strings.HasSuffix(service, "APIService") {
								sdk.services[service] = field.Names[0].Name
							}
						}
					}
				}
			}
		}
	}
	return sdk, nil
}

func sdkRoute(method *ast.FuncDecl) (string, string) {
	path, operation := "", ""
	ast.Inspect(method.Body, func(node ast.Node) bool {
		if assignment, ok := node.(*ast.AssignStmt); ok && len(assignment.Lhs) == 1 && len(assignment.Rhs) == 1 {
			if name, ok := assignment.Lhs[0].(*ast.Ident); ok && name.Name == "localVarPath" {
				if expression, ok := assignment.Rhs[0].(*ast.BinaryExpr); ok && expression.Op == token.ADD {
					if literal, ok := expression.Y.(*ast.BasicLit); ok && literal.Kind == token.STRING {
						path, _ = strconv.Unquote(literal.Value)
					}
				}
			}
		}
		if value, ok := node.(*ast.ValueSpec); ok && len(value.Names) == 1 && len(value.Values) == 1 && value.Names[0].Name == "localVarHTTPMethod" {
			if selector, ok := value.Values[0].(*ast.SelectorExpr); ok && typeString(selector.X) == "http" {
				operation = strings.ToUpper(strings.TrimPrefix(selector.Sel.Name, "Method"))
			}
		}
		return true
	})
	return path, operation
}

func (sdk bulkSDK) setter(request, parameter, valueType string) (string, error) {
	var found []string
	for key, method := range sdk.methods {
		if !strings.HasPrefix(key, request+".") || len(method.Type.Params.List) != 1 || method.Type.Results == nil || len(method.Type.Results.List) != 1 || typeString(method.Type.Results.List[0].Type) != request {
			continue
		}
		if typeString(method.Type.Params.List[0].Type) != valueType {
			continue
		}
		if parameter != "" && !strings.EqualFold(method.Name.Name, strings.ReplaceAll(parameter, "_", "")) {
			continue
		}
		found = append(found, method.Name.Name)
	}
	if len(found) != 1 {
		sort.Strings(found)
		return "", fmt.Errorf("%s needs one setter for %q (%s), found %v", request, parameter, valueType, found)
	}
	return found[0], nil
}

func sdkQueryParameter(method *ast.FuncDecl, name string) bool {
	found := false
	ast.Inspect(method.Body, func(node ast.Node) bool {
		call, ok := node.(*ast.CallExpr)
		if !ok || typeString(call.Fun) != "parameterAddToHeaderOrQuery" || len(call.Args) < 2 || typeString(call.Args[0]) != "localVarQueryParams" {
			return true
		}
		literal, ok := call.Args[1].(*ast.BasicLit)
		if ok && literal.Kind == token.STRING {
			value, _ := strconv.Unquote(literal.Value)
			found = found || value == name
		}
		return true
	})
	return found
}

func (sdk bulkSDK) operation(resource spec.ResourceSpec, operation, splitKey string) (bulkOperation, error) {
	var matches []string
	for key, method := range sdk.methods {
		if !strings.HasSuffix(method.Name.Name, "Execute") || len(method.Type.Params.List) != 1 {
			continue
		}
		service := strings.SplitN(key, ".", 2)[0]
		if sdk.services[service] == "" {
			continue
		}
		path, verb := sdkRoute(method)
		if path == resource.API.EndpointPath && verb == operation {
			matches = append(matches, key)
		}
	}
	if len(matches) != 1 {
		sort.Strings(matches)
		return bulkOperation{}, fmt.Errorf("%s %s needs one SDK binding, found %v", operation, resource.API.EndpointPath, matches)
	}
	execute := sdk.methods[matches[0]]
	service := strings.TrimPrefix(typeString(execute.Recv.List[0].Type), "*")
	methodName := strings.TrimSuffix(execute.Name.Name, "Execute")
	request := typeString(execute.Type.Params.List[0].Type)
	constructor := sdk.methods[service+"."+methodName]
	if constructor == nil || len(constructor.Type.Params.List) != 1 || typeString(constructor.Type.Params.List[0].Type) != "context.Context" || constructor.Type.Results == nil || len(constructor.Type.Results.List) != 1 || typeString(constructor.Type.Results.List[0].Type) != request {
		return bulkOperation{}, fmt.Errorf("%s has no compatible SDK request constructor", matches[0])
	}
	requestExecute := sdk.methods[request+".Execute"]
	if requestExecute == nil || len(requestExecute.Type.Params.List) != 0 || requestExecute.Type.Results == nil || len(requestExecute.Type.Results.List) != 2 || typeString(requestExecute.Type.Results.List[0].Type) != "*http.Response" || typeString(requestExecute.Type.Results.List[1].Type) != "error" {
		return bulkOperation{}, fmt.Errorf("%s has no compatible SDK Execute method", request)
	}
	result := bulkOperation{Service: sdk.services[service], Method: methodName, Request: request}
	if operation == "PUT" || operation == "PATCH" {
		var bodies []string
		for name, structure := range sdk.structs {
			if wrapper, ok := structure.Fields[resource.API.RequestWrapperKey]; ok && strings.HasPrefix(wrapper.GoType, "*map[string]") {
				if _, err := sdk.setter(request, "", name); err == nil {
					bodies = append(bodies, name)
				}
			}
		}
		if len(bodies) != 1 {
			sort.Strings(bodies)
			return bulkOperation{}, fmt.Errorf("%s needs one SDK body with wrapper %q, found %v", request, resource.API.RequestWrapperKey, bodies)
		}
		result.BodyType = bodies[0]
		wrapper := sdk.structs[result.BodyType].Fields[resource.API.RequestWrapperKey]
		result.BodyField = wrapper.GoName
		result.ValueType = strings.TrimPrefix(wrapper.GoType, "*map[string]")
		if _, exists := sdk.structs[result.ValueType]; !exists {
			return bulkOperation{}, fmt.Errorf("%s wrapper has no SDK value struct %s", result.BodyType, result.ValueType)
		}
		result.Setter, _ = sdk.setter(request, "", result.BodyType)
	}
	if operation == "DELETE" {
		if !sdkQueryParameter(execute, resource.API.DeleteParameter) {
			return bulkOperation{}, fmt.Errorf("%s has no query parameter %q", matches[0], resource.API.DeleteParameter)
		}
		var err error
		result.Setter, err = sdk.setter(request, resource.API.DeleteParameter, "[]string")
		if err != nil {
			return bulkOperation{}, err
		}
	}
	if splitKey != "" {
		if !sdkQueryParameter(execute, splitKey) {
			return bulkOperation{}, fmt.Errorf("%s has no split query parameter %q", matches[0], splitKey)
		}
		var err error
		result.SplitSetter, err = sdk.setter(request, splitKey, "string")
		if err != nil {
			return bulkOperation{}, err
		}
	}
	return result, nil
}

func planBulkBindings(registry spec.Registry, sdk bulkSDK) ([]bulkBinding, error) {
	if err := registry.ValidateBulkOrders(); err != nil {
		return nil, err
	}
	groups := map[string][]spec.ResourceSpec{}
	for _, resource := range registry {
		groups[resource.API.BulkKey] = append(groups[resource.API.BulkKey], resource)
	}
	keys := make([]string, 0, len(groups))
	for key := range groups {
		keys = append(keys, key)
	}
	sort.Strings(keys)
	var bindings []bulkBinding
	for _, key := range keys {
		variants := groups[key]
		sort.Slice(variants, func(i, j int) bool { return variants[i].TerraformType < variants[j].TerraformType })
		resource := variants[0]
		binding := bulkBinding{Key: key, Variants: variants, Operations: map[string]bulkOperation{}}
		for _, variant := range variants {
			if variant.API.EndpointPath != resource.API.EndpointPath || variant.API.RequestWrapperKey != resource.API.RequestWrapperKey || variant.API.DeleteParameter != resource.API.DeleteParameter || variant.Operations != resource.Operations {
				return nil, fmt.Errorf("bulk variants for %s disagree on API route or operations", key)
			}
			if len(variant.API.FixedHeaders) > 1 || (len(variants) > 1 && len(variant.API.FixedHeaders) != 1) {
				return nil, fmt.Errorf("bulk key %s requires one fixed split parameter per variant", key)
			}
			for splitKey, value := range variant.API.FixedHeaders {
				if value == "" || (binding.SplitKey != "" && binding.SplitKey != splitKey) {
					return nil, fmt.Errorf("bulk key %s has inconsistent split parameters", key)
				}
				binding.SplitKey = splitKey
			}
		}
		for _, operation := range []struct {
			verb    string
			enabled bool
		}{
			{"PUT", resource.Operations.Create}, {"PATCH", resource.Operations.Update}, {"DELETE", resource.Operations.Delete}, {"GET", resource.Operations.Read},
		} {
			if !operation.enabled {
				continue
			}
			planned, err := sdk.operation(resource, operation.verb, binding.SplitKey)
			if err != nil {
				return nil, fmt.Errorf("bulk key %s: %w", key, err)
			}
			binding.Operations[operation.verb] = planned
		}
		bindings = append(bindings, binding)
	}
	return bindings, nil
}

func writeBulkCallback(buf *bytes.Buffer, binding bulkBinding, verb string) {
	op, exists := binding.Operations[verb]
	if !exists {
		return
	}
	field := strings.ToUpper(verb[:1]) + strings.ToLower(verb[1:]) + "Func"
	arguments := "c *openapi.APIClient, ctx context.Context"
	request := fmt.Sprintf("c.%s.%s(ctx)", op.Service, op.Method)
	switch verb {
	case "PUT", "PATCH":
		arguments += ", req interface{}"
		request += fmt.Sprintf(".%s(*req.(*openapi.%s))", op.Setter, op.BodyType)
	case "DELETE":
		arguments += ", names []string"
		request += fmt.Sprintf(".%s(names)", op.Setter)
	}
	if binding.SplitKey != "" {
		field = "Header" + field
		arguments += ", headers map[string]string"
	}
	fmt.Fprintf(buf, "%s: func(%s) (*http.Response, error) {\n", field, arguments)
	if binding.SplitKey == "" {
		fmt.Fprintf(buf, "return %s.Execute()\n", request)
	} else {
		fmt.Fprintf(buf, "request := %s\nif value, ok := headers[%q]; ok { request = request.%s(value) }\nreturn request.Execute()\n", request, binding.SplitKey, op.SplitSetter)
	}
	buf.WriteString("},\n")
}

func renderBulkBindings(registry spec.Registry, bindings []bulkBinding) ([]byte, error) {
	var buf bytes.Buffer
	buf.WriteString("package bulkops\n\nimport (\"context\"; \"net/http\"; \"reflect\"; \"terraform-provider-verity/openapi\")\n\nvar resourceRegistry = map[string]ResourceConfig{\n")
	for _, binding := range bindings {
		put, exists := binding.Operations["PUT"]
		if !exists {
			put = binding.Operations["PATCH"]
		}
		patch := binding.Operations["PATCH"]
		fmt.Fprintf(&buf, "%q: {ResourceType: %q, HeaderSplitKey: %q,\n", binding.Key, binding.Key, binding.SplitKey)
		fmt.Fprintf(&buf, "ResponseCollectionKey: %q,\n", binding.Variants[0].API.ResponseCollectionKey)
		for _, verb := range []string{"PUT", "PATCH"} {
			if op, exists := binding.Operations[verb]; exists {
				name := strings.ToUpper(verb[:1]) + strings.ToLower(verb[1:])
				fmt.Fprintf(&buf, "Prepare%s: func(data map[string]interface{}) (interface{}, error) { return prepareTypedBulkRequest(data, func(request *openapi.%s, values map[string]openapi.%s) { request.%s = &values }) },\n", name, op.BodyType, op.ValueType, op.BodyField)
			}
		}
		if put.BodyType != "" {
			fmt.Fprintf(&buf, "PutRequestType: reflect.TypeOf(openapi.%s{}),\n", put.BodyType)
		}
		if patch.BodyType != "" {
			fmt.Fprintf(&buf, "PatchRequestType: reflect.TypeOf(openapi.%s{}),\n", patch.BodyType)
		}
		fmt.Fprintf(&buf, "APIClientGetter: func(c *openapi.APIClient) ResourceAPIClient { return &GenericAPIClient{client: c, resourceType: %q} },\n", binding.Key)
		for _, verb := range []string{"PUT", "PATCH", "DELETE", "GET"} {
			writeBulkCallback(&buf, binding, verb)
		}
		if binding.SplitKey != "" {
			variants := append([]spec.ResourceSpec(nil), binding.Variants...)
			sort.Slice(variants, func(i, j int) bool {
				return variants[i].API.FixedHeaders[binding.SplitKey] < variants[j].API.FixedHeaders[binding.SplitKey]
			})
			fmt.Fprintf(&buf, "HeaderResponseExtractor: func(raw map[string]interface{}, headers map[string]string) (map[string]interface{}, error) {\nkey := %q\nswitch headers[%q] {\n", variants[0].API.ResponseCollectionKey, binding.SplitKey)
			for _, variant := range variants {
				fmt.Fprintf(&buf, "case %q: key = %q\n", variant.API.FixedHeaders[binding.SplitKey], variant.API.ResponseCollectionKey)
			}
			buf.WriteString("}\nif data, ok := raw[key].(map[string]interface{}); ok { return data, nil }\nreturn make(map[string]interface{}), nil\n},\n")
		}
		buf.WriteString("},\n")
	}
	buf.WriteString("}\n")
	cacheKeys := make(map[string]bool)
	for _, resource := range registry {
		cacheKeys[resource.API.CacheKey] = true
	}
	keys := make([]string, 0, len(cacheKeys))
	for key := range cacheKeys {
		keys = append(keys, key)
	}
	sort.Strings(keys)
	buf.WriteString("\nvar finalCacheRefreshKeys = []string{\n")
	for _, key := range keys {
		fmt.Fprintf(&buf, "%q,\n", key)
	}
	buf.WriteString("}\n")
	for _, mode := range []spec.Mode{spec.ModeDatacenter, spec.ModeCampus} {
		for _, verb := range []string{"PUT", "PATCH", "DELETE"} {
			fmt.Fprintf(&buf, "\nvar %s%sOrder = []string{\n", mode, strings.ToUpper(verb[:1])+strings.ToLower(verb[1:]))
			for _, key := range registry.BulkOperationOrder(mode, verb) {
				fmt.Fprintf(&buf, "%q,\n", key)
			}
			buf.WriteString("}\n")
		}
	}
	return format.Source(buf.Bytes())
}

func generateBulk(opts bulkOptions) error {
	if opts.Registry == "" || opts.OpenAPIDir == "" || opts.Output == "" {
		return fmt.Errorf("--registry, --openapi-dir and --output are required")
	}
	registry, err := readRegistry(opts.Registry)
	if err != nil {
		return err
	}
	sdk, err := discoverBulkSDK(opts.OpenAPIDir)
	if err != nil {
		return err
	}
	bindings, err := planBulkBindings(registry, sdk)
	if err != nil {
		return err
	}
	output, err := renderBulkBindings(registry, bindings)
	if err != nil {
		return fmt.Errorf("format bulk bindings: %w", err)
	}
	if opts.Check {
		existing, err := os.ReadFile(opts.Output)
		if err != nil {
			return err
		}
		if !bytes.Equal(existing, output) {
			return fmt.Errorf("%s differs: run tools/generate_provider.sh --write", opts.Output)
		}
		return nil
	}
	return writeFile(opts.Output, output)
}
