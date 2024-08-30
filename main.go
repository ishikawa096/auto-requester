package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/go-co-op/gocron/v2"
)

var scheduler gocron.Scheduler

func init() {
	var err error
	// initialize the scheduler
	scheduler, err = gocron.NewScheduler()
	if err != nil {
		logger("Error creating scheduler:", err)
		os.Exit(1)
	}
}

func CreateTask(url string, method func(string) (*http.Response, error)) gocron.Task {
	return gocron.NewTask(
		func(url string, method func(string) (*http.Response, error)) {
			resp, err := method(url)
			if err != nil {
				logger("Error sending request:", err)
				return
			}
			defer resp.Body.Close()
			body, err := io.ReadAll(resp.Body)
			if err != nil {
				logger("Error reading response body:", err)
				return
			}

			logger("Response Body:", string(body), "Status:", resp.Status)
		},
		"https://catfact.ninja/fact",
		http.Get,
	)
}

func registerJob(minSec, maxSec int, url string, method string) {
	// TODO: サポートするHTTPメソッドを増やす
	methods := map[string]func(string) (*http.Response, error){
		"GET": http.Get,
		// "POST": http.Post,
	}
	httpMethod, exists := methods[strings.ToUpper(method)]
	if !exists {
		logger("Error: Unsupported HTTP method:", method)
		return
	}

	_, err := scheduler.NewJob(
		gocron.DurationRandomJob(
			time.Duration(minSec)*time.Second,
			time.Duration(maxSec)*time.Second,
		),
		CreateTask(
			url,
			httpMethod,
		),
	)
	if err != nil {
		logger("Error creating job:", err)
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
		logger("Error stoping jobs:", err)
		os.Exit(1)
	} else {
		logger(time.Now(), "Scheduler Stop Jobs!")
	}
}

// Start job handler
func StartHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	response := map[string]string{"message": "Scheduler Start! To stop, please access /stop"}
	json.NewEncoder(w).Encode(response)

	scheduler.Start()
	logger("Scheduler Restart!")
}

func currentTimeStr() string {
	tz, err := time.LoadLocation("Asia/Tokyo")
	if err != nil {
		logger("Error loading location:", err)
		return ""
	}
	return time.Now().In(tz).Format("2006-01-02 15:04:05 MST")
}

func logger(msgs ...interface{}) {
	var messages []string
	for _, msg := range msgs {
		switch v := msg.(type) {
		case string:
			messages = append(messages, v)
		case error:
			messages = append(messages, v.Error())
		default:
			messages = append(messages, fmt.Sprintf("%v", v))
		}
	}
	// TODO: colorize, log level
	fmt.Printf("%s %s\n", currentTimeStr(), strings.Join(messages, " "))
}

func main() {
	// TODO: 環境変数から取得する
	registerJob(4, 6, "https://catfact.ninja/fact", "GET")

	// start the scheduler
	scheduler.Start()
	logger("Scheduler Start!")

	// register handlers
	http.HandleFunc("/stop", StopHandler)
	http.HandleFunc("/start", StartHandler)

	// start the HTTP server in a goroutine
	go func() {
		logger("HTTP server is running on port 8080")
		if err := http.ListenAndServe(":8080", nil); err != nil {
			logger("Error starting HTTP server:", err)
		}
	}()

	// block main goroutine forever
	select {}
}
