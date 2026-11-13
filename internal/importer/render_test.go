package importer

import (
	"reflect"
	"strings"
	"testing"

	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/hclsyntax"
	"github.com/zclconf/go-cty/cty"

	"terraform-provider-verity/internal/spec"
)

func TestRecursiveRendererAndPruningUseAliasesAtEveryDepth(t *testing.T) {
	fields := []spec.FieldSpec{
		{APIName: "Parent", TerraformName: "configuration", Kind: spec.FieldKindObject, Fields: []spec.FieldSpec{
			{APIName: "Entries", TerraformName: "members", Kind: spec.FieldKindList, ElementKind: spec.FieldKindObject, Collection: &spec.CollectionSpec{IdentityField: "index"}, Fields: []spec.FieldSpec{
				{APIName: "Order", TerraformName: "index", Kind: spec.FieldKindInt64},
				{APIName: "Config", TerraformName: "options", Kind: spec.FieldKindObject, Fields: []spec.FieldSpec{
					{APIName: "Nested", TerraformName: "children", Kind: spec.FieldKindList, ElementKind: spec.FieldKindObject, Fields: []spec.FieldSpec{
						{APIName: "Title", TerraformName: "label", Kind: spec.FieldKindString},
						{APIName: "Scores", TerraformName: "values", Kind: spec.FieldKindList, ElementKind: spec.FieldKindNumber},
					}},
				}},
				{APIName: "Values", TerraformName: "flags", Kind: spec.FieldKindList, ElementKind: spec.FieldKindBool},
			}},
		}},
		{APIName: "Tags", TerraformName: "labels", Kind: spec.FieldKindList, ElementKind: spec.FieldKindString},
	}
	supported := &SchemaFields{Attributes: map[string]bool{"labels": true}, Blocks: map[string]*SchemaFields{
		"configuration": {Blocks: map[string]*SchemaFields{
			"members": {Attributes: map[string]bool{"index": true, "flags": true}, Blocks: map[string]*SchemaFields{
				"options": {Blocks: map[string]*SchemaFields{
					"children": {Attributes: map[string]bool{"label": true, "values": true}},
				}},
			}},
		}},
	}}
	object := map[string]interface{}{
		"Parent": map[string]interface{}{"Entries": []interface{}{
			map[string]interface{}{"Order": float64(2), "Values": []interface{}{true, false, nil}, "Config": map[string]interface{}{"Nested": []interface{}{
				map[string]interface{}{"Title": "child", "Scores": []interface{}{float64(1), float64(2.5)}, "unknown": "omit"},
			}}},
		}},
		"Tags": []interface{}{"a", "b"},
	}
	imp := &Importer{}
	imp.pruneObject("verity_example", supported, object, "", nil, fields)
	if got, want := imp.UnsupportedFields(), map[string][]string{"verity_example": {"configuration.members.options.children.unknown"}}; !reflect.DeepEqual(got, want) {
		t.Fatalf("unsupported fields = %v, want %v", got, want)
	}
	got, _, err := imp.generateResourceTF(map[string]map[string]interface{}{"example": object}, ResourceConfig{ResourceType: "example", StageName: "example_stage", Fields: fields})
	if err != nil {
		t.Fatal(err)
	}
	want := "\nresource \"verity_example\" \"example\" {\n" +
		"    name = \"example\"\n    depends_on = [verity_operation_stage.example_stage]\n" +
		"\tconfiguration {\n\t\tmembers {\n\t\t\tindex = 2\n" +
		"\t\t\toptions {\n\t\t\t\tchildren {\n\t\t\t\t\tvalues = [\n\t\t\t\t\t\t1,\n\t\t\t\t\t\t2.5,\n\t\t\t\t\t]\n" +
		"\t\t\t\t\tlabel = \"child\"\n\t\t\t\t}\n\t\t\t}\n" +
		"\t\t\tflags = [\n\t\t\t\ttrue,\n\t\t\t\tfalse,\n\t\t\t\tnull,\n\t\t\t]\n\t\t}\n\t}\n" +
		"\tlabels = [\n\t\t\"a\",\n\t\t\"b\",\n\t]\n}\n\n"
	if got != want {
		t.Fatalf("generated HCL:\n%s\nwant:\n%s", got, want)
	}
	if _, diags := hclsyntax.ParseConfig([]byte(got), "import.tf", hcl.InitialPos); diags.HasErrors() {
		t.Fatal(diags)
	}
}

func TestImportTasksPreserveACLFileNames(t *testing.T) {
	for _, mode := range []string{"campus", "datacenter"} {
		tasks, err := (&Importer{Mode: mode}).importTasks()
		if err != nil {
			t.Fatal(err)
		}
		found := map[string]string{}
		for _, task := range tasks {
			found[task.terraformResourceType] = task.name
		}
		if found["verity_acl_v4"] != "acls_ipv4" || found["verity_acl_v6"] != "acls_ipv6" {
			t.Fatalf("%s ACL filenames changed: %v", mode, found)
		}
	}
}

func TestRendererAndPruningHonorAliasedRootSkipFields(t *testing.T) {
	fields := []spec.FieldSpec{
		{APIName: "ObjectName", TerraformName: "name", Kind: spec.FieldKindString},
		{APIName: "RowIndex", TerraformName: "index", Kind: spec.FieldKindInt64},
	}
	skip := map[string]bool{"name": true, "index": true}
	object := map[string]interface{}{"ObjectName": "x", "RowIndex": float64(2)}
	imp := &Importer{}
	imp.pruneObject("verity_example", &SchemaFields{}, object, "", skip, fields)
	if got := imp.UnsupportedFields(); len(got) != 0 {
		t.Fatalf("root identity/index aliases reported as unsupported: %v", got)
	}
	got, _, err := imp.generateResourceTF(map[string]map[string]interface{}{"x": object}, ResourceConfig{ResourceType: "example", StageName: "example_stage", Fields: fields, SkipTopLevelKeys: skip})
	if err != nil {
		t.Fatal(err)
	}
	if strings.Count(got, "name =") != 1 || strings.Contains(got, "index =") {
		t.Fatalf("root aliases ignored skip fields:\n%s", got)
	}
}

func TestRendererUsesDeclaredAutoAssignmentFlagWithAliases(t *testing.T) {
	fields := []spec.FieldSpec{
		{APIName: "Value", TerraformName: "address", Kind: spec.FieldKindString, AutoAssignment: &spec.AutoAssignmentSpec{FlagField: "automatic"}},
		{APIName: "Flag", TerraformName: "automatic", Kind: spec.FieldKindBool},
	}
	for _, automatic := range []bool{true, false} {
		got, _, err := (&Importer{}).generateResourceTF(map[string]map[string]interface{}{"x": {"Value": "server-assigned", "Flag": automatic}}, ResourceConfig{ResourceType: "example", StageName: "example_stage", Fields: fields})
		if err != nil {
			t.Fatal(err)
		}
		if strings.Contains(got, "address =") == automatic || !strings.Contains(got, "automatic =") {
			t.Fatalf("automatic %t:\n%s", automatic, got)
		}
	}
}

func TestRendererHandlesEmptyAndNullCollections(t *testing.T) {
	fields := []spec.FieldSpec{
		{APIName: "properties", TerraformName: "properties", Kind: spec.FieldKindObject},
		{APIName: "entries", TerraformName: "entries", Kind: spec.FieldKindList, ElementKind: spec.FieldKindObject},
		{APIName: "labels", TerraformName: "labels", Kind: spec.FieldKindList, ElementKind: spec.FieldKindString},
	}
	for _, test := range []struct {
		name   string
		object map[string]interface{}
		labels string
	}{
		{name: "empty", object: map[string]interface{}{"properties": map[string]interface{}{}, "entries": []interface{}{}, "labels": []interface{}{}}, labels: "labels = []"},
		{name: "null", object: map[string]interface{}{"properties": nil, "entries": nil, "labels": nil}, labels: "labels = null"},
	} {
		t.Run(test.name, func(t *testing.T) {
			got, _, err := (&Importer{}).generateResourceTF(map[string]map[string]interface{}{"x": test.object}, ResourceConfig{ResourceType: "example", StageName: "example_stage", Fields: fields})
			if err != nil {
				t.Fatal(err)
			}
			if !strings.Contains(got, "properties {}") || !strings.Contains(got, test.labels) || strings.Contains(got, "entries") {
				t.Fatalf("empty/null collection HCL:\n%s", got)
			}
		})
	}
}

func TestRendererRejectsMalformedValuesInsteadOfLosingData(t *testing.T) {
	for _, test := range []struct {
		kind    spec.FieldKind
		element spec.FieldKind
		value   interface{}
		want    string
	}{
		{spec.FieldKindObject, "", "bad", "expected object"},
		{spec.FieldKindList, spec.FieldKindObject, true, "expected list of objects"},
		{spec.FieldKindList, spec.FieldKindObject, []interface{}{"bad"}, "entry 0: expected object"},
		{spec.FieldKindList, spec.FieldKindString, []interface{}{map[string]interface{}{}}, "entry 0: expected JSON scalar"},
	} {
		got, _, err := (&Importer{}).generateResourceTF(map[string]map[string]interface{}{"x": {"bad": test.value}}, ResourceConfig{ResourceType: "example", StageName: "example_stage", Fields: []spec.FieldSpec{{APIName: "bad", TerraformName: "bad", Kind: test.kind, ElementKind: test.element}}})
		if got != "" || err == nil || !strings.Contains(err.Error(), "verity_example \"x\": bad:") || !strings.Contains(err.Error(), test.want) {
			t.Fatalf("HCL %q, error %v, want %q", got, err, test.want)
		}
	}
}

func TestRendererWritesStringsAsLiteralHCL(t *testing.T) {
	value := "${not_a_reference} %{not_a_directive}\n\"quoted\"\\path"
	formatted, err := scalarHCL(value)
	if err != nil {
		t.Fatal(err)
	}
	expr, diags := hclsyntax.ParseExpression([]byte(formatted), "value.tf", hcl.InitialPos)
	if diags.HasErrors() {
		t.Fatal(diags)
	}
	got, diags := expr.Value(nil)
	if diags.HasErrors() || !got.RawEquals(cty.StringVal(value)) {
		t.Fatalf("literal value = %v, diagnostics = %v", got, diags)
	}
}

func TestRendererNaturalOrderIsDeterministicForNumericTies(t *testing.T) {
	objects := map[string]map[string]interface{}{"node10": {}, "node2": {}, "node1": {}, "node01": {}}
	var previous string
	for run := 0; run < 30; run++ {
		got, resources, err := (&Importer{}).generateResourceTF(objects, ResourceConfig{ResourceType: "example", StageName: "example_stage"})
		if err != nil {
			t.Fatal(err)
		}
		if run > 0 && got != previous {
			t.Fatal("natural sorting depends on map iteration order")
		}
		previous = got
		last := -1
		for index, name := range []string{"node01", "node1", "node2", "node10"} {
			if want := (ImportedResource{TerraformType: "verity_example", TerraformName: name, ID: name}); resources[index] != want {
				t.Fatalf("resource %d = %#v, want %#v", index, resources[index], want)
			}
			position := strings.Index(got, "name = \""+name+"\"")
			if position <= last {
				t.Fatalf("incorrect natural order:\n%s", got)
			}
			last = position
		}
	}
}

func TestRendererReturnsTopLevelImportIdentity(t *testing.T) {
	fields := []spec.FieldSpec{{APIName: "properties", TerraformName: "properties", Kind: spec.FieldKindObject, Fields: []spec.FieldSpec{
		{APIName: "name", TerraformName: "name", Kind: spec.FieldKindString},
	}}}
	config, resources, err := (&Importer{}).generateResourceTF(map[string]map[string]interface{}{
		"edge\"node": {"name": "different-body-name", "properties": map[string]interface{}{"name": "nested-name"}},
	}, ResourceConfig{ResourceType: "example", StageName: "example_stage", Fields: fields})
	if err != nil {
		t.Fatal(err)
	}
	want := []ImportedResource{{TerraformType: "verity_example", TerraformName: "edge_node", ID: "edge\"node"}}
	if !reflect.DeepEqual(resources, want) {
		t.Fatalf("resources = %#v, want %#v", resources, want)
	}
	file, diags := hclsyntax.ParseConfig([]byte(config), "resource.tf", hcl.InitialPos)
	if diags.HasErrors() {
		t.Fatal(diags)
	}
	block := file.Body.(*hclsyntax.Body).Blocks[0]
	if !reflect.DeepEqual(block.Labels, []string{want[0].TerraformType, want[0].TerraformName}) {
		t.Fatalf("resource labels = %v, want %#v", block.Labels, want[0])
	}
	id, diags := block.Body.Attributes["name"].Expr.Value(nil)
	if diags.HasErrors() || !id.RawEquals(cty.StringVal(want[0].ID)) {
		t.Fatalf("resource name = %v, diagnostics = %v", id, diags)
	}
}

func TestRendererNeverWritesSecrets(t *testing.T) {
	encrypted := spec.FieldSpec{APIName: "ssh_key_or_password_encrypted", TerraformName: "ssh_key_or_password_encrypted", Kind: spec.FieldKindString, Sensitive: true, Secret: spec.SecretMasked,
		AutoAssignment: &spec.AutoAssignmentSpec{FlagField: "ssh_key_or_password_encrypted_auto_assigned_"}}
	fields := []spec.FieldSpec{
		{APIName: "enable", TerraformName: "enable", Kind: spec.FieldKindBool},
		{APIName: "ssh_key_or_password", TerraformName: "ssh_key_or_password", Kind: spec.FieldKindString, Sensitive: true, Secret: spec.SecretWriteOnly},
		encrypted,
		{APIName: "ssh_key_or_password_encrypted_auto_assigned_", TerraformName: "ssh_key_or_password_encrypted_auto_assigned_", Kind: spec.FieldKindBool},
		{APIName: "md5_password_encrypted", TerraformName: "md5_password_encrypted", Kind: spec.FieldKindString, Sensitive: true, Secret: spec.SecretMasked},
	}
	objects := map[string]map[string]interface{}{
		"with-password": {"enable": true, "ssh_key_or_password": "", "ssh_key_or_password_encrypted": "this field has been set",
			"ssh_key_or_password_encrypted_auto_assigned_": false, "md5_password_encrypted": "This field has been set"},
		"no-password": {"enable": true, "ssh_key_or_password": "", "ssh_key_or_password_encrypted": "",
			"ssh_key_or_password_encrypted_auto_assigned_": false, "md5_password_encrypted": ""},
		"auto-assigned": {"enable": true, "ssh_key_or_password_encrypted": "this field has been set",
			"ssh_key_or_password_encrypted_auto_assigned_": true},
	}
	got, _, err := (&Importer{}).generateResourceTF(objects, ResourceConfig{ResourceType: "switchpoint", StageName: "switchpoint_stage", Fields: fields})
	if err != nil {
		t.Fatal(err)
	}
	if strings.Contains(got, "ssh_key_or_password =") || strings.Contains(got, "_encrypted =") || strings.Contains(strings.ToLower(got), "= \"this field has been set\"") {
		t.Fatalf("a secret value was written:\n%s", got)
	}
	blocks := map[string]string{}
	for _, block := range strings.Split(got, "\nresource ")[1:] {
		blocks[strings.Fields(block)[1]] = block
	}
	for _, note := range []string{"# md5_password_encrypted: a password is set on the system", "# ssh_key_or_password_encrypted: a password is set on the system"} {
		if !strings.Contains(blocks[`"with-password"`], note) {
			t.Errorf("with-password: missing note %q:\n%s", note, blocks[`"with-password"`])
		}
	}
	for _, name := range []string{`"no-password"`, `"auto-assigned"`} {
		if strings.Contains(blocks[name], "a password is set") {
			t.Errorf("%s: unexpected password note:\n%s", name, blocks[name])
		}
	}
	file, diags := hclsyntax.ParseConfig([]byte(got), "switchpoint.tf", hcl.InitialPos)
	if diags.HasErrors() || file == nil {
		t.Fatalf("generated HCL does not parse: %v", diags)
	}
}

func TestRendererNeverWritesReadOnlyFields(t *testing.T) {
	fields := []spec.FieldSpec{
		{APIName: "eth_port_paths", TerraformName: "eth_port_paths", Kind: spec.FieldKindList, ElementKind: spec.FieldKindObject,
			Collection: &spec.CollectionSpec{IdentityField: "index"}, Fields: []spec.FieldSpec{
				{APIName: "index", TerraformName: "index", Kind: spec.FieldKindInt64},
				{APIName: "port_name", TerraformName: "port_name", Kind: spec.FieldKindString, Access: spec.AccessComputed},
				{APIName: "eth_port_num_ee_profile", TerraformName: "eth_port_num_ee_profile", Kind: spec.FieldKindString},
			}},
	}
	objects := map[string]map[string]interface{}{"b1": {"eth_port_paths": []interface{}{
		map[string]interface{}{"index": float64(1), "port_name": "Eth1/1", "eth_port_num_ee_profile": "p1"},
	}}}
	got, _, err := (&Importer{}).generateResourceTF(objects, ResourceConfig{ResourceType: "bundle", StageName: "bundle_stage", Fields: fields})
	if err != nil {
		t.Fatal(err)
	}
	if strings.Contains(got, "port_name") {
		t.Fatalf("a read-only field was written:\n%s", got)
	}
	if !strings.Contains(got, `eth_port_num_ee_profile = "p1"`) || !strings.Contains(got, "index = 1") {
		t.Fatalf("other members are missing:\n%s", got)
	}
}
