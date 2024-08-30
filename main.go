package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/go-co-op/gocron/v2"
)

var scheduler gocron.Scheduler

func init() {
	var err error
	// initialize the scheduler
	scheduler, err = gocron.NewScheduler()
	if err != nil {
		fmt.Println("Error creating scheduler:", err)
		os.Exit(1)
	}
}

// Stop job handler
func StopHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	response := map[string]string{"message": "job stop! To restart, please access /start"}
	json.NewEncoder(w).Encode(response)

	err := scheduler.StopJobs()
	if err != nil {
		fmt.Println("Error stoping jobs:", err)
		os.Exit(1)
	} else {
		fmt.Println(time.Now(), "Scheduler Stop Jobs!")
	}
}

// Start job handler
func StartHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	response := map[string]string{"message": "Scheduler Start! To stop, please access /stop"}
	json.NewEncoder(w).Encode(response)

	scheduler.Start()
	fmt.Println(time.Now(), "Scheduler Restart!")
}

func main() {
	_, err := scheduler.NewJob(
		gocron.DurationRandomJob(
			// TODO: 指定した時間間隔でジョブを実行する
			4*time.Second,
			6*time.Second,
		),
		gocron.NewTask(
			func(a string, b time.Time) {
				// TODO: httpリクエストを送信する

				fmt.Println(time.Now(), "executing task with params:", a, b)
			},
			"hello",
			time.Now(),
		),
	)
	if err != nil {
		fmt.Println("Error creating job:", err)
		os.Exit(1)
	}

	// start the scheduler
	scheduler.Start()
	fmt.Println(time.Now(), "Scheduler Start!")

	// register handlers
	http.HandleFunc("/stop", StopHandler)
	http.HandleFunc("/start", StartHandler)

	// start the HTTP server in a goroutine
	go func() {
		fmt.Println("HTTP server is running on port 8080")
		if err := http.ListenAndServe(":8080", nil); err != nil {
			fmt.Println("Error starting HTTP server:", err)
		}
	}()

	// block main goroutine forever
	select {}
}
