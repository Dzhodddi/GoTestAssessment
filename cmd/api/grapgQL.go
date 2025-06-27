package main

import (
	"WorkAssigment/graph"
	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/handler/extension"
	"github.com/99designs/gqlgen/graphql/handler/lru"
	"github.com/99designs/gqlgen/graphql/handler/transport"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/labstack/echo/v4"
	"github.com/vektah/gqlparser/v2/ast"
)

// playgroundHandler godoc
//
//	@Summary		GraphQL Playground
//	@Description	Serves the GraphQL Playground UI for testing queries
//	@Tags			GraphQL
//	@Accept			*/*
//	@Produce		text/html
//	@Success		200	{string}	string	"HTML page"
//	@Router			/playground [get]
func (app *application) playgroundHandler() echo.HandlerFunc {
	return echo.WrapHandler(playground.Handler("GraphQL Playground", "/v1/graphql"))
}

// playgroundHandler godoc
//
//	@Summary		GraphQL Playground
//	@Description	Serves the GraphQL Playground UI for testing queries
//	@Tags			GraphQL
//	@Accept			*/*
//	@Param			query	query		string	false	"Query"
//	@Success		200		{string}	string	"HTML page"
//	@Router			/graphql [get]
func (app *application) registerGraphQL() echo.HandlerFunc {
	srv := handler.New(graph.NewExecutableSchema(app.config.graphQLConfig.config))
	srv.AddTransport(transport.Options{})
	srv.AddTransport(transport.GET{})
	srv.AddTransport(transport.POST{})
	srv.AddTransport(transport.MultipartForm{})
	srv.SetQueryCache(lru.New[*ast.QueryDocument](1000))

	srv.Use(extension.Introspection{})
	srv.Use(extension.AutomaticPersistedQuery{
		Cache: lru.New[string](100),
	})
	return echo.WrapHandler(srv)
}
