package main

import (
	"fmt"
	"os"
	"strconv"

	"github.com/MichaelWittgreffe/ministub/pkg/api"
	"github.com/MichaelWittgreffe/ministub/pkg/config"
	"github.com/MichaelWittgreffe/ministub/pkg/logger"
)

func main() {
	log := logger.NewLogger("std")
	log.Info("Loading Definition")

	cfgPath, bindHost, port, err := parseArgs()
	if err != nil {
		log.Fatal(fmt.Sprintf("Startup Error: %s", err.Error()))
	}

	cfg, err := config.LoadFromFile(cfgPath)
	if err != nil {
		log.Fatal(fmt.Sprintf("Unable To Load From From Path %s: %s", cfgPath, err.Error()))
	}

	if err = config.Validate(cfg); err != nil {
		log.Fatal(fmt.Sprintf("Config Validation Error: %s", err.Error()))
	}

	log.Info("Config Loaded, Executing Startup Actions")

	if server := api.NewHTTPAPI(log, cfg); server != nil {
		go api.ExecuteActions(cfg.StartupActions, "Startup", cfg, log) // run startup actions before starting server

		log.Fatal(
			fmt.Sprintf(
				"Fatal Error: %s",
				server.ListenAndServe(bindHost, port).Error(),
			),
		)
	}
}

// parseArgs parses the cmd args and returns
func parseArgs() (cfgPath string, bind string, port int, err error) {
	for i, data := range os.Args {
		switch {
		case data == "-h":
			fmt.Printf("\n")
			os.Exit(0)
		case data == "-p":
			port, err = strconv.Atoi(os.Args[i+1])
		case data == "-b":
			bind = os.Args[i+1]
		default:
			if i > 0 && os.Args[i-1] != "-p" && os.Args[i-1] != "-b" {
				cfgPath = os.Args[i]
			}
		}
	}

	if port == 0 {
		port = 8080
	}
	if len(bind) == 0 {
		bind = "0.0.0.0"
	}
	if len(cfgPath) == 0 {
		if cwd, err := os.Getwd(); err == nil {
			cfgPath = fmt.Sprintf("%s/ministub.yml", cwd)
		} else {
			return "", "", -1, err
		}
	}

	return cfgPath, bind, port, err
}
