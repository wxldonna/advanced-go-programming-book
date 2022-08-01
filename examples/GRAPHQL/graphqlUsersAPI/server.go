package main

import (
	"log"
	"net/http"
	"os"

	"github.xiaoliang.graphql.users/tasks"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.xiaoliang.graphql.users/graph"
	"github.xiaoliang.graphql.users/graph/generated"
)

const defaultPort = "8090"

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = defaultPort
	}
	resolver := &graph.Resolver{
		Tasks: tasks.New(),
		Atts:  tasks.NewAttachments(),
	}
	srv := handler.NewDefaultServer(generated.NewExecutableSchema(generated.Config{Resolvers: resolver}))

	http.Handle("/", playground.Handler("GraphQL playground", "/query"))
	http.Handle("/query", srv)

	log.Printf("connect to http://localhost:%s/ for GraphQL playground", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
