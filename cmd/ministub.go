package main

import (
	"fmt"
	"os"
	"strconv"
	"time"

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

	log.Info("Config Loaded, Executing Startup Actions")

	if api := api.NewHTTPAPI(log, cfg); api != nil {
		go executeStartupActions(log, cfg)

		log.Fatal(
			fmt.Sprintf("Fatal Error: %s",
				api.ListenAndServe(bindHost, port).Error(),
			),
		)
	}
}

func executeStartupActions(log logger.Logger, cfg *config.Config) {
	for _, actionEntry := range cfg.StartupActions {
		for name, action := range actionEntry {
			switch {
			case name == "delay":
				if period, valid := action.(int); valid {
					log.Info(fmt.Sprintf("Delay Requested For %d Seconds", period))
					time.Sleep(time.Duration(period) * time.Second)
				}
			case name == "request":
				var target *config.Service
				var request *config.Request

				if actionEntry, valid := action.(map[interface{}]interface{}); valid {
					for k, v := range actionEntry {
						switch {
						case k.(string) == "target":
							target = cfg.Services[v.(string)]
						case k.(string) == "id":
							request = cfg.Requests[v.(string)]
						}
					}
				}

				if target != nil && request != nil {
					log.Info("Startup Request Requested, Currently Not Written")
				} else {
					log.Error("Invalid Startup Request Requested")
				}
			}
		}
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
			if i > 0 && os.Args[i-1] != "-p" && os.Args[i-1] != "-h" {
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
