package main

import (
	"log"
	"net/http"
	"os"

	"github.com/opensaucerer/barf"
	"github.com/opensaucerer/barf/app/database"
	"github.com/opensaucerer/barf/app/global"
	"github.com/opensaucerer/barf/app/version"
)

func main() {

	// load environment variables
	if err := barf.Env(global.ENV, os.Getenv("ENV_PATH")); err != nil {
		log.Fatal(err)
	}

	// configure barf
	allow := true
	if err := barf.Stark(barf.Augment{
		Port:     global.ENV.Port,
		Logging:  &allow, // enable request logging
		Recovery: &allow, // enable panic recovery
		CORS: &barf.CORS{
			AllowedOrigins: []string{"https://*.onrender.com"},
			MaxAge:         3600,
			AllowedMethods: []string{
				http.MethodGet,
				http.MethodPost,
				http.MethodPut,
				http.MethodPatch,
				http.MethodDelete,
			},
		},
	}); err != nil {
		log.Fatal(err)
	}

	// apply global barf middleware
	barf.Hippocampus().Hijack(func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

			// a custom security feature. only allow requests containing a valid app token in the header key "zeina-mfi"
			// this can probably be improved by using asymmetric encryption
			if r.Header.Get("zeina-mfi") != global.ENV.AppToken {
				barf.Response(w).Status(http.StatusUnauthorized).JSON(nil)
				return
			}

			h.ServeHTTP(w, r)
		})
	})

	if err := database.NewPostgreSQLConnection(global.ENV.PostgreSQLURI, global.ENV.PostgreSQLConnections); err != nil {
		log.Fatal(err)
	}

	if err := database.ReadFileAndExecuteQueries(global.ENV.SQLFilePath); err != nil {
		log.Fatal(err)
	}

	// preload v1 routes
	version.V1()

	// call upon barf to listen and serve
	if err := barf.Beck(); err != nil {
		log.Fatal(err)
	}
}
