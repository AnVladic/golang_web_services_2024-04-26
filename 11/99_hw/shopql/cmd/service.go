package cmd

import (
	"encoding/json"
	"github.com/99designs/gqlgen/graphql/handler"
	"hw11_shopql/internal"
	"hw11_shopql/internal/graphql"
	"hw11_shopql/internal/user"
	"io"
	"net/http"
	"os"
)

func LoadData(path string) *internal.MarketplaceData {
	jsonFile, err := os.Open(path)
	if err != nil {
		panic(err)
	}
	defer func(jsonFile *os.File) {
		err := jsonFile.Close()
		if err != nil {
			panic(err)
		}
	}(jsonFile)
	byteValue, _ := io.ReadAll(jsonFile)
	var data internal.MarketplaceData
	err = json.Unmarshal(byteValue, &data)
	if err != nil {
		panic(err)
	}
	return &data
}

func GetApp() http.Handler {
	marketplaceData := LoadData("testdata.json")
	resolver := graphql.Resolver{
		Catalogs: marketplaceData.Catalog.ToMap(),
		Sellers:  marketplaceData.SellersToMap(),
	}
	userHandler := user.Handler{
		SessionManager: *internal.CreateSessionManager(),
	}

	config := graphql.Config{Resolvers: &resolver}
	config.Directives.Authorized = graphql.AuthorizedDirective

	mux := http.NewServeMux()
	srv := handler.NewDefaultServer(graphql.NewExecutableSchema(config))
	mux.Handle("/query", user.AuthMiddleware(userHandler, srv))
	mux.HandleFunc("/register", userHandler.Register)
	return mux
}
