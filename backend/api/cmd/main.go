// backend/api/cmd/main.go
package main

import (
	"log"
	"net/http"
	"github.com/sagit-chu/flvx-monitor/backend/api"
	"github.com/sagit-chu/flvx-monitor/backend/db"
	"github.com/sagit-chu/flvx-monitor/backend/repository"
)

func main() {
	database, err := db.InitDB("flvx.db")
	if err != nil {
		log.Fatal(err)
	}
	defer database.Close()
	
	nodeRepo := repository.NewNodeRepository(database)

	srv := api.NewServer(nodeRepo)
	log.Println("Starting server on :8080")
	if err := http.ListenAndServe(":8080", srv.Router); err != nil {
		log.Fatal(err)
	}
}
