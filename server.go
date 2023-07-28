package main

import (
	"net/http"

	"git.sr.ht/~mpldr/uniview/graph"

	"git.sr.ht/~poldi1405/glog"
	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/handler/transport"
	"github.com/99designs/gqlgen/graphql/playground"
)

func startServer() error {
	srv := handler.NewDefaultServer(graph.NewExecutableSchema(graph.Config{Resolvers: &graph.Resolver{}}))

	srv.AddTransport(&transport.Websocket{})

	http.Handle("/", srv)
	http.Handle("/play", playground.Handler("GraphQL playground", "/"))

	glog.Info("connect to http://localhost:8080/play for GraphQL playground")
	return http.ListenAndServe(":8080", nil)
}
