package evedt

import (
	"log"
	"net/http"
	"time"

	"github.com/trichner/evedt/tracker"
)

var (
	repo *tracker.Repo
)

func Run() {

	config, err := LoadConfig("config.toml")
	if err != nil {
		log.Fatalf("Cannot read config file: %s", err)
	}

	repo = &tracker.Repo{}
	if err := repo.Open(); err != nil {
		log.Fatalf("Cannot open DB: %s", err)
	}
	defer repo.Close()

	// init replicator
	replicator := tracker.Replicator{}
	replicator.Init(repo, config.ApiCredentials.ApiKey, config.ApiCredentials.VCode, 1000)

	// schedule replicator. Note: This runs forever
	ticker := time.NewTicker(15 * time.Minute)
	go func() {
		for {
			if err := replicator.Run(); err != nil {
				log.Printf("Eve API failed: %s\n", err)
			}
			<-ticker.C
		}
	}()

	// setup and start our REST api
	prefix := config.ServerConfig.Prefix
	router := NewRouter(prefix)

	port := config.ServerConfig.Port
	log.Printf("Docking at localhost:%s%s\n", port, prefix)

	log.Fatal(http.ListenAndServe(":"+port, router))
}
