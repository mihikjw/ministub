package main

import (
	"fmt"
	"os"

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

	if cfg != nil {
		log.Info("Config Loaded")
	}
}
