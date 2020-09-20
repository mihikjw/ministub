package api

import (
	"fmt"
	"time"

	"github.com/MichaelWittgreffe/ministub/pkg/config"
	"github.com/MichaelWittgreffe/ministub/pkg/logger"
)

// ExecuteActions runs any actions requested, performs own logging and is designed to be run in its own goroutine
func ExecuteActions(actions []map[string]interface{}, requester string, cfg *config.Config, log logger.Logger) {
	for _, actionEntry := range actions {
		for name, action := range actionEntry {
			switch {
			case name == "delay":
				if period, valid := action.(int); valid {
					log.Info(fmt.Sprintf("%s: Delay Requested For %d Seconds", requester, period))
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
					log.Info(fmt.Sprintf("%s: Request Requested, Currently Not Written", requester))
				} else {
					log.Error("Invalid Request Requested")
				}
			}
		}
	}
}
