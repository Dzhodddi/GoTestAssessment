package main

import (
	"WorkAssigment/graph"
	db "WorkAssigment/internal/db"
	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/handler/extension"
	"github.com/99designs/gqlgen/graphql/handler/lru"
	"github.com/99designs/gqlgen/graphql/handler/transport"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/vektah/gqlparser/v2/ast"
	"log"
	"net/http"
)

const graphQL = "42069"

func main() {
	database, err := db.New("postgresql://postgres:adminadmin@localhost:5432/testGo?sslmode=disable", 5, 5, "15m")
	if err != nil {
		log.Fatal(err)
	}
	resolver := graph.NewResolver(database)

	srv := handler.New(graph.NewExecutableSchema(graph.Config{Resolvers: resolver}))
	srv.AddTransport(transport.Options{})
	srv.AddTransport(transport.GET{})
	srv.AddTransport(transport.POST{})

	srv.SetQueryCache(lru.New[*ast.QueryDocument](1000))

	srv.Use(extension.Introspection{})
	srv.Use(extension.AutomaticPersistedQuery{
		Cache: lru.New[string](100),
	})

	http.Handle("/", playground.Handler("GraphQL playground", "/query"))
	http.Handle("/query", srv)

	log.Printf("connect to http://localhost:%s/ for GraphQL playground", graphQL)
	_ = http.ListenAndServe(":"+graphQL, nil)
}
