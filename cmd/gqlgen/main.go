package main

import (
	"fmt"
	"go/types"
	"os"
	"strings"

	"github.com/vektah/gqlparser/v2/ast"

	"github.com/99designs/gqlgen/api"
	"github.com/99designs/gqlgen/codegen/config"
	"github.com/99designs/gqlgen/plugin/modelgen"
)

func DirectiveStr(fd ast.DirectiveList) string {
	names := []string{}
	for _, d := range fd {
		names = append(names, d.Name)
	}
	return "[" + strings.Join(names, ",") + "]"
}

// Defining mutation function
func toStrArrFieldHook(td *ast.Definition, fd *ast.FieldDefinition, f *modelgen.Field) (*modelgen.Field, error) {
	if f, err := modelgen.DefaultFieldMutateHook(td, fd, f); err != nil {
		return f, err
	}
	fmt.Printf("field hook %s.%s, %s@%s \n", td.Name, f.Name, f.GoName, DirectiveStr(fd.Directives))
	c := fd.Directives.ForName("strArr")
	if c != nil {
		fmt.Println(f.GoName, "->", "[]string")
		f.Type = types.NewSlice(types.Typ[types.String])
		f.GoName = "[]string"
		// f.Type = types.String
	}

	c = fd.Directives.ForName("str")
	if c != nil {
		fmt.Println(f.GoName, "->", "string")
		f.Type = types.Typ[types.String]
		f.GoName = "string"
	}
	return f, nil
}

func main() {
	cfg, err := config.LoadConfigFromDefaultLocations()
	if err != nil {
		fmt.Fprintln(os.Stderr, "failed to load config", err.Error())
		os.Exit(2)
	}

	// Attaching the mutation function onto modelgen plugin
	p := modelgen.Plugin{
		FieldHook: toStrArrFieldHook,
	}

	err = api.Generate(cfg, api.ReplacePlugin(&p))

	if err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		os.Exit(3)
	}
}
