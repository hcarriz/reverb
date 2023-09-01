//go:build ignore
// +build ignore

package main

import (
	"log/slog"
	"os"

	"entgo.io/contrib/entgql"
	"entgo.io/ent/entc"
	"entgo.io/ent/entc/gen"
)

func main() {

	log := slog.Default()

	ex, err := entgql.NewExtension(
		entgql.WithWhereInputs(true),
		entgql.WithConfigPath("./gqlgen.yml"),
		entgql.WithSchemaGenerator(),
		entgql.WithSchemaPath("schema/ent.graphqls"),
	)
	if err != nil {
		log.Error("unable to create extension", slog.Any("error", err.Error()))
		os.Exit(1)
	}

	config := &gen.Config{
		Target:    "ent",
		Package:   "reverb/generated/ent",
		Schema:    "reverb/generated/schema",
		Templates: entgql.AllTemplates,
		Features:  gen.AllFeatures,
	}

	opts := []entc.Option{
		entc.Extensions(ex),
	}

	if err := entc.Generate("./schema", config, opts...); err != nil {
		log.Error("unable to generate ent", slog.Any("error", err.Error()))
		os.Exit(1)
	}

	log.Info("generated schema")

}
