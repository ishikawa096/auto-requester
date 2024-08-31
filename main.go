package main

import (
	"net/http"

	"github.com/ishikawa096/auto-requester/handlers"
	"github.com/ishikawa096/auto-requester/scheduler"
	"github.com/ishikawa096/auto-requester/utils"
)

func main() {
	scheduler.StartJob()

	// register handlers
	http.HandleFunc("/stop", handlers.StopHandler)
	http.HandleFunc("/start", handlers.StartHandler)

	// start the HTTP server in a goroutine
	go func() {
		utils.Logger("HTTP server is running on port 8080")
		if err := http.ListenAndServe(":8080", nil); err != nil {
			utils.Logger("Error starting HTTP server:", err)
		}
	}()

	// block main goroutine forever
	select {}
}
