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

	cfgPath, bindHost, port, err := parseArgs()
	if err != nil {
		logFatal(log, fmt.Sprintf("Startup Error: %s", err.Error()))
	}

	log.Info("Loading Config...")

	cfg, err := config.LoadFromFile(cfgPath)
	if err != nil {
		logFatal(log, fmt.Sprintf("Unable To Load From From Path %s: %s", cfgPath, err.Error()))
	}

	if err = config.Validate(cfg); err != nil {
		logFatal(log, fmt.Sprintf("Config Validation Error: %s", err.Error()))
	}

	log.Info(fmt.Sprintf("Config Loaded From Path: %s", cfgPath))

	requester := api.NewRequester("http")

	if server := api.NewHTTPAPI(log, cfg, requester); server != nil {
		if cfg.StartupActions != nil && len(cfg.StartupActions) > 0 {
			log.Info("Executing Startup Actions...")
			go api.ExecuteActions(cfg.StartupActions, "Startup", cfg, log, requester)
		}

		logFatal(
			log,
			fmt.Sprintf(
				"Fatal Error: %s",
				server.ListenAndServe(bindHost, port).Error(),
			),
		)
	}
}

// logFatal prints the given message to the logger error stream then os.Exit(1)
func logFatal(log logger.Logger, msg string) {
	log.Error(msg)
	os.Exit(1)
}

// parseArgs parses the cmd args and returns
func parseArgs() (cfgPath string, bind string, port int, err error) {
	for i, data := range os.Args {
		switch {
		case data == "-h":
			fmt.Printf("ministub is an API stubbing tool allowing follow-on actions from an incoming request\n\nUsage:\nministub [path]\n\t-h: Help\n\t-p: Port\n\t-b: Accept Host\n")
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
			if string(cwd[len(cwd)-1]) != "/" {
				cwd += "/"
			}
			cfgPath = fmt.Sprintf("%sministub.yml", cwd)
		} else {
			return "", "", -1, err
		}
	}

	return cfgPath, bind, port, err
}
