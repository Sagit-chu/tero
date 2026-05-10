// backend/api/cmd/main.go
package main

import (
	"log"
	"net/http"
	"github.com/sagit-chu/flvx-monitor/backend/api"
)

func main() {
	srv := api.NewServer()
	log.Println("Starting server on :8080")
	if err := http.ListenAndServe(":8080", srv.Router); err != nil {
		log.Fatal(err)
	}
}
