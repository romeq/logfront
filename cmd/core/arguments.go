package main

import (
	"flag"
	"os"
	"strings"
)

const (
	EnvRunApi = "LOGFRONT_RUN_API"
)

type arguments struct {
	runApi     bool
	configFile string
}

func parseArguments() arguments {
	defaultEnvBool := func(env string) bool {
		envValueInLower := strings.ToLower(os.Getenv(env))
		return envValueInLower != "" && envValueInLower != "false"
	}

	var args arguments
	flag.BoolVar(&args.runApi, "api", defaultEnvBool(EnvRunApi), "Run the REST API server")
	flag.StringVar(&args.configFile, "f", "config.yaml", "Config file")
	flag.Parse()

	return args
}
