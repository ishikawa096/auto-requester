package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/ishikawa096/auto-requester/handlers"
	"github.com/ishikawa096/auto-requester/scheduler"
	"github.com/ishikawa096/auto-requester/utils"
)

const (
	defaultPort = "8080"
)

func main() {
	scheduler.StartJob()

	// register handlers
	http.HandleFunc("/stop", handlers.StopHandler)
	http.HandleFunc("/start", handlers.StartHandler)

	port := defaultPort
	if value, exists := os.LookupEnv("PORT"); exists {
		port = value
	}

	// start the HTTP server in a goroutine
	go func() {
		utils.Logger("HTTP server is running on port", port)
		utils.Logger(fmt.Sprintf("To STOP 👉 http://localhost:%s/stop", port))
		utils.Logger(fmt.Sprintf("To START 👉 http://localhost:%s/start", port))
		if err := http.ListenAndServe(":"+port, nil); err != nil {
			utils.Logger("Error starting HTTP server:", err)
		}
	}()

	// block main goroutine forever
	select {}
}
