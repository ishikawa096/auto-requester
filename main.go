package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/go-co-op/gocron/v2"
)

type httpMethodOptions struct {
	contentType string
	body        io.Reader
}

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

func sendRequest(method string, url string, options httpMethodOptions) (*http.Response, error) {
	switch strings.ToUpper(method) {
	case "GET":
		return http.Get(url)
	case "POST":
		return http.Post(url, options.contentType, options.body)
	default:
		return nil, fmt.Errorf("unsupported HTTP method: %s", method)
	}
}

func CreateTask(method string, url string, options httpMethodOptions) gocron.Task {
	return gocron.NewTask(
		func(method string, url string, options httpMethodOptions) {
			resp, err := sendRequest(method, url, options)
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
		method,
		url,
		options,
	)
}

func registerJob(minSec, maxSec int, method string, url string, options httpMethodOptions) {
	_, err := scheduler.NewJob(
		gocron.DurationRandomJob(
			time.Duration(minSec)*time.Second,
			time.Duration(maxSec)*time.Second,
		),
		CreateTask(
			method,
			url,
			options,
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

func getSettingValues() (int, int, string, string, httpMethodOptions) {
	// default values
	minSec := 4
	maxSec := 6
	url := "http://localhost:3000"
	method := "GET"
	contentType := "application/json"

	if val, exists := os.LookupEnv("INTERVAL_MIN_SEC"); exists {
		if v, err := strconv.Atoi(val); err == nil {
			minSec = v
		}
	}
	if val, exists := os.LookupEnv("INTERVAL_MAX_SEC"); exists {
		if v, err := strconv.Atoi(val); err == nil {
			maxSec = v
		}
	}
	if val, exists := os.LookupEnv("HTTP_METHOD"); exists {
		method = val
	}
	if val, exists := os.LookupEnv("TARGET_URL"); exists {
		url = val
	}
	if val, exists := os.LookupEnv("HTTP_HEADERS"); exists {
		contentType = val
	}
	file, err := os.Open("/etc/app/request.json")
	if err != nil {
		logger("Error opening file:", err)
		os.Exit(1)
	}
	defer file.Close()
	// read the file content to memory
	fileContent, err := io.ReadAll(file)
	if err != nil {
		logger("Error reading file:", err)
		os.Exit(1)
	}
	options := httpMethodOptions{
		contentType: contentType,
		body:        strings.NewReader(string(fileContent)),
	}
	return minSec, maxSec, method, url, options
}

func main() {
	registerJob(getSettingValues())

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
