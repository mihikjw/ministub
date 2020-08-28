package main

import (
	"fmt"
	"os"
	"time"

	"github.com/MichaelWittgreffe/ministub/pkg/api"
	"github.com/MichaelWittgreffe/ministub/pkg/config"
	"github.com/MichaelWittgreffe/ministub/pkg/logger"
)

func main() {
	log := logger.NewLogger("std")
	log.Info("Loading Definition")

	var cfgPath string
	if len(os.Args) >= 2 {
		cfgPath = os.Args[1]
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
				api.ListenAndServe("localhost", 8080).Error(),
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
