package config

import (
	"errors"
	"flag"
	"log"
	"os"
)

// Version represents the version of the binary at build time as set by ldflags
type Version string

var version string

const defaultPort = 3000

// Config holds configurations variables
type Config struct {
	Version     Version
	RunLocal    bool
	Port        int
	AllowOrigin string
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

// MustBuildConfig attempts to load various config from environment, flags, and baked-in variables.
func MustBuildConfig() *Config {
	config, err := buildConfig()
	if err != nil {
		log.Fatal("Failed to load config: " + err.Error())
	}
	return config
}

func buildConfig() (*Config, error) {

	runLocal := flag.Bool("l", false, "Serve the api locally rather than as a lambda handler")
	port := flag.Int("p", defaultPort, "Port to bind local server to")
	flag.Parse()

	origin, err := getEnv("ALLOW_ORIGIN")
	if err != nil {
		return nil, err
	}

	return &Config{
		Version:     Version(version),
		RunLocal:    *runLocal,
		Port:        *port,
		AllowOrigin: origin,
	}, nil
}
