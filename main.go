package main

import (
	"damo/pkg/config"
	"damo/pkg/health"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/akrylysov/algnhsa"
	"github.com/julienschmidt/httprouter"
	"github.com/rs/cors"
)

func main() {
	conf := config.MustBuildConfig()

	mux := httprouter.New()
	mux.Handler("GET", "/health", health.CheckHandler(conf.Version))

	corsHandler := cors.New(cors.Options{
		AllowedOrigins:   strings.Split(conf.AllowOrigin, ","),
		AllowCredentials: true,
	})

	handler := corsHandler.Handler(mux)

	if conf.RunLocal {
		addr := ":" + strconv.Itoa(conf.Port)
		log.Println("Listening on", addr)
		log.Fatal(http.ListenAndServe(addr, handler))
	} else {
		// Lambda
		algnhsa.ListenAndServe(handler, nil)
	}
}
