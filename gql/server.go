package gql

import (
	"net/http"

	"github.com/spf13/viper"

	"github.com/99designs/gqlgen/handler"
	bam "github.com/cosmos/cosmos-sdk/baseapp"
	"github.com/wirelineio/wirechain/x/registry"
)

const defaultPort = "8080"

// Server configures and starts the GQL server.
func Server(baseApp *bam.BaseApp, keeper registry.Keeper) {
	if viper.GetBool("gql-server") {
		port := viper.GetString("gql-port")
		if port == "" {
			port = defaultPort
		}

		http.Handle("/", handler.Playground("GraphQL playground", "/query"))
		http.Handle("/query", handler.GraphQL(NewExecutableSchema(Config{Resolvers: &Resolver{baseApp: baseApp, keeper: keeper}})))
		http.ListenAndServe(":"+port, nil)
	}
}
