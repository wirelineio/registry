//
// Copyright 2019 Wireline, Inc.
//

package gql

import (
	"net/http"

	"github.com/spf13/viper"

	"github.com/99designs/gqlgen/handler"
	bam "github.com/cosmos/cosmos-sdk/baseapp"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/x/auth"
	"github.com/wirelineio/registry/x/registry"

	"github.com/go-chi/chi"
	"github.com/rs/cors"
)

const defaultPort = "9473"

// Server configures and starts the GQL server.
func Server(baseApp *bam.BaseApp, cdc *codec.Codec, keeper registry.Keeper, accountKeeper auth.AccountKeeper) {
	if viper.GetBool("gql-server") {
		port := viper.GetString("gql-port")
		if port == "" {
			port = defaultPort
		}

		router := chi.NewRouter()

		// Add CORS middleware around every request
		// See https://github.com/rs/cors for full option listing
		router.Use(cors.New(cors.Options{
			AllowedOrigins: []string{"*"},
			Debug:          true,
		}).Handler)

		if viper.GetBool("gql-playground") {
			router.Handle("/", handler.Playground("Wireline Registry", "/query"))
		}

		router.Handle("/query", handler.GraphQL(NewExecutableSchema(Config{Resolvers: &Resolver{
			baseApp:       baseApp,
			codec:         cdc,
			keeper:        keeper,
			accountKeeper: accountKeeper,
		}})))

		err := http.ListenAndServe(":"+port, router)
		if err != nil {
			panic(err)
		}
	}
}
