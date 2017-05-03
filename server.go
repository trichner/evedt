package evedt

import (
	"log"
	"net/http"
	"time"

	"github.com/trichner/evedt/tracker"
)

const (
	confFile = "config.toml"
)

func Start() error {

	config, err := loadConfig(confFile)
	if err != nil {
		return err
	}

	// init replicator
	apiKey := config.ApiCredentials.ApiKey
	vCode := config.ApiCredentials.VCode
	replicator, err := tracker.NewReplicator(tracker.ApiCredentials(apiKey, vCode))

	// schedule replicator. Note: This runs forever
	quit := make(chan struct{})
	go func() {

		ticker := time.NewTicker(15 * time.Minute)
		for {
			if err := replicator.Run(); err != nil {
				log.Printf("Eve API failed: %s\n", err)
			}

			select {
			case <-quit:
				ticker.Stop()
				return
			case <-ticker.C:
				break;
			}
		}
	}()

	// setup and start our REST api
	prefix := config.ServerConfig.Prefix
	router := NewRouter(prefix, &replicator)

	port := config.ServerConfig.Port
	log.Printf("Docking at :%s%s\n", port, prefix)

	return http.ListenAndServe(":"+port, router)
}
