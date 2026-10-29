package main

import (
	"bytes"
	_ "embed"
	"flag"
	"fmt"
	"go/ast"
	"go/format"
	"go/parser"
	"go/token"
	"os"
	"path/filepath"
	"strings"
)

//go:embed call_api.go.tmpl
var clientCustomization []byte

func main() {
	dir := flag.String("sdk-dir", "", "temporary generated SDK directory")
	flag.Parse()
	if *dir == "" {
		fmt.Fprintln(os.Stderr, "--sdk-dir is required")
		os.Exit(2)
	}
	if err := prepare(*dir); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func prepare(dir string) error {
	clientPath := filepath.Join(dir, "client.go")
	source, err := os.ReadFile(clientPath)
	if err != nil {
		return err
	}
	source, err = customizeClient(source)
	if err != nil {
		return err
	}
	if err := os.WriteFile(clientPath, source, 0644); err != nil {
		return err
	}
	entries, err := os.ReadDir(dir)
	if err != nil {
		return err
	}
	for _, entry := range entries {
		path := filepath.Join(dir, entry.Name())
		if entry.IsDir() || !strings.HasSuffix(entry.Name(), ".go") || strings.HasSuffix(entry.Name(), "_test.go") {
			if err := os.RemoveAll(path); err != nil {
				return err
			}
			continue
		}
		data, err := os.ReadFile(path)
		if err != nil {
			return err
		}
		data, err = format.Source(data)
		if err != nil {
			return fmt.Errorf("format %s: %w", entry.Name(), err)
		}
		if err := os.WriteFile(path, data, 0644); err != nil {
			return err
		}
	}
	return nil
}

func callAPI(source []byte) (*token.FileSet, *ast.FuncDecl, error) {
	fset := token.NewFileSet()
	file, err := parser.ParseFile(fset, "client.go", source, 0)
	if err != nil {
		return nil, nil, err
	}
	var found *ast.FuncDecl
	for _, decl := range file.Decls {
		fn, ok := decl.(*ast.FuncDecl)
		if !ok || fn.Name.Name != "callAPI" {
			continue
		}
		if found != nil {
			return nil, nil, fmt.Errorf("multiple callAPI declarations")
		}
		found = fn
	}
	if found == nil || found.Recv == nil || len(found.Recv.List) != 1 || found.Body == nil {
		return nil, nil, fmt.Errorf("generated client has no callAPI method")
	}
	return fset, found, nil
}

func customizeClient(source []byte) ([]byte, error) {
	fset, original, err := callAPI(source)
	if err != nil {
		return nil, err
	}
	customSet, custom, err := callAPI(clientCustomization)
	if err != nil {
		return nil, err
	}
	signature := func(set *token.FileSet, fn *ast.FuncDecl) (string, error) {
		var buf bytes.Buffer
		for _, name := range fn.Recv.List[0].Names {
			fmt.Fprintf(&buf, "%s ", name.Name)
		}
		if err := format.Node(&buf, set, fn.Recv.List[0].Type); err != nil {
			return "", err
		}
		if err := format.Node(&buf, set, fn.Type); err != nil {
			return "", err
		}
		return buf.String(), nil
	}
	originalSignature, err := signature(fset, original)
	if err != nil {
		return nil, err
	}
	customSignature, err := signature(customSet, custom)
	if err != nil {
		return nil, err
	}
	if originalSignature != customSignature {
		return nil, fmt.Errorf("generated callAPI signature changed; review credential-redaction customization")
	}
	start, end := fset.Position(original.Pos()).Offset, fset.Position(original.End()).Offset
	customStart, customEnd := customSet.Position(custom.Pos()).Offset, customSet.Position(custom.End()).Offset
	result := append([]byte{}, source[:start]...)
	result = append(result, clientCustomization[customStart:customEnd]...)
	result = append(result, source[end:]...)
	return format.Source(result)
}
