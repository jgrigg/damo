package main

import (
	"adv-caja-x-api/pkg/auth"
	"adv-caja-x-api/pkg/health"
	"adv-caja-x-api/pkg/user"
	"errors"
	"flag"
	"github.com/akrylysov/algnhsa"
	"github.com/julienschmidt/httprouter"
	"github.com/rs/cors"
	"log"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"
)

var version string

const defaultPort = 3000

type config struct {
	Auth0Domain   string
	Auth0Audience string
	RunLocal      bool
	Port          int
	TalapiBase    url.URL
	AllowOrigin   string
}

func getEnv(name string) (string, error) {
	val := os.Getenv(name)
	if val == "" {
		error := "Missing env var " + name
		return "", errors.New(error)
	}
	log.Printf("Environment : %s:%s", name, val)
	return val, nil
}

func buildConfig() (*config, error) {

	runLocal := flag.Bool("l", false, "Serve the api locally rather than as a lambda handler")
	port := flag.Int("p", defaultPort, "Port to bind local server to")
	flag.Parse()

	auth0Domain, err := getEnv("AUTH0_DOMAIN")
	if err != nil {
		return nil, err
	}

	auth0Audience, err := getEnv("AUTH0_AUDIENCE")
	if err != nil {
		return nil, err
	}

	talapiBase, err := getEnv("TALAPI_BASE")
	if err != nil {
		return nil, err
	}

	talapiBaseURI, err := url.Parse(talapiBase)
	log.Println("base " + talapiBaseURI.RequestURI())
	if err != nil {
		return nil, err
	}

	origin, err := getEnv("ALLOW_ORIGIN")
	if err != nil {
		return nil, err
	}

	return &config{
		Auth0Audience: auth0Audience,
		Auth0Domain:   auth0Domain,
		RunLocal:      *runLocal,
		Port:          *port,
		TalapiBase:    *talapiBaseURI,
		AllowOrigin:   origin,
	}, nil
}

func main() {
	config, err := buildConfig()
	if err != nil {
		log.Fatal("Failed to load config: " + err.Error())
	}
	corsHandler := cors.New(cors.Options{
		AllowedOrigins:   strings.Split(config.AllowOrigin, ","),
		AllowCredentials: true,
	})
	router := httprouter.New()
	// Not authenticated
	router.Handler("GET", "/health", corsHandler.Handler(health.CheckHandler(config.TalapiBase, version)))
	var userfactory auth.UserFactory

	// Authenticated
	authHandler := auth.AuthHandlerBuilder(auth.NewSSOValidator(config.Auth0Domain, config.Auth0Audience), userfactory)

	router.Handler("GET", "/user/me", corsHandler.Handler(authHandler(user.MeHandler())))

	if config.RunLocal {
		addr := ":" + strconv.Itoa(config.Port)
		log.Println("Listening on", addr)
		log.Fatal(http.ListenAndServe(addr, router))
	} else {
		// Lambda
		algnhsa.ListenAndServe(router, nil)
	}
}
