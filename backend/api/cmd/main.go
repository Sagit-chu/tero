// backend/api/cmd/main.go
package main

import (
	"log"
	"net/http"
	"time"

	"github.com/sagit-chu/flvx-monitor/backend/api"
	"github.com/sagit-chu/flvx-monitor/backend/cloudflare"
	"github.com/sagit-chu/flvx-monitor/backend/db"
	"github.com/sagit-chu/flvx-monitor/backend/flvx"
	"github.com/sagit-chu/flvx-monitor/backend/monitor"
	"github.com/sagit-chu/flvx-monitor/backend/repository"
)

func main() {
	database, err := db.InitDB("flvx.db")
	if err != nil {
		log.Fatal(err)
	}
	defer database.Close()

	nodeRepo := repository.NewNodeRepository(database)
	configRepo := repository.NewConfigRepository(database)

	// Fetch required configuration for the monitor loop
	flvxURL, _ := configRepo.GetConfig("flvx_api_url")
	flvxAcc, _ := configRepo.GetConfig("flvx_account")
	flvxPass, _ := configRepo.GetConfig("flvx_password")
	cfToken, _ := configRepo.GetConfig("cf_token")
	domainName, _ := configRepo.GetConfig("domain_name")
	checkAPIURL, _ := configRepo.GetConfig("check_api_url")

	// Instantiate clients and monitor
	pinger := monitor.NewTCPPinger(5 * time.Second)
	gfw := monitor.NewPingPeScraper(checkAPIURL)
	flvxClient := flvx.NewClient(flvxURL, flvxAcc, flvxPass)
	cfClient := cloudflare.NewClient("", cfToken)

	monitorSvc := monitor.NewMonitorService(nodeRepo, configRepo, pinger, gfw, flvxClient, cfClient, domainName)

	// Run monitor loop in background
	go func() {
		ticker := time.NewTicker(5 * time.Minute)
		defer ticker.Stop()
		for {
			status := monitorSvc.RunCycle()
			log.Printf("[Monitor] Cycle completed. Status: %s", status)
			<-ticker.C
		}
	}()

	srv := api.NewServer(nodeRepo, configRepo, monitorSvc)
	log.Println("Starting server on :8080")
	if err := http.ListenAndServe(":8080", srv.Router); err != nil {
		log.Fatal(err)
	}
}
