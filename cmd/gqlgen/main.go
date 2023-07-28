package main

import (
	"fmt"
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

func main() {
	cfg, err := config.LoadConfigFromDefaultLocations()
	if err != nil {
		fmt.Fprintln(os.Stderr, "failed to load config", err.Error())
		os.Exit(2)
	}

	// Attaching the mutation function onto modelgen plugin
	p := modelgen.Plugin{
		MutateHook: func(b *modelgen.ModelBuild) *modelgen.ModelBuild {
			return b
		},
	}

	err = api.Generate(cfg, api.ReplacePlugin(&p))

	if err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		os.Exit(3)
	}
}
