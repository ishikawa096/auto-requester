package scheduler

import (
	"bytes"
	"io"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/go-co-op/gocron/v2"
	"github.com/ishikawa096/auto-requester/utils"
)

type httpOptions struct {
	method      string // GET, POST, PUT, DELETE
	url         string
	contentType string
	body        []byte
}

var Scheduler gocron.Scheduler

func init() {
	var err error
	// initialize the scheduler
	Scheduler, err = gocron.NewScheduler()
	if err != nil {
		utils.Logger("Error creating scheduler:", err)
		os.Exit(1)
	}
}

func sendRequest(options httpOptions) (*http.Response, error) {
	// リクエストボディを再利用可能な形にするために、bytes.Readerを使用
	bodyReader := bytes.NewReader(options.body)

	req, err := http.NewRequest(options.method, options.url, bodyReader)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", options.contentType)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

func createTask(options httpOptions) gocron.Task {
	return gocron.NewTask(
		func(options httpOptions) {
			resp, err := sendRequest(options)
			if err != nil {
				utils.Logger("Error sending request:", err)
				return
			}
			defer resp.Body.Close()
			body, err := io.ReadAll(resp.Body)
			if err != nil {
				utils.Logger("Error reading response body:", err)
				return
			}

			utils.Logger("Response Body:", string(body), "Status:", resp.Status)
		},
		options,
	)
}

func registerJob(minSec, maxSec int, options httpOptions) {
	_, err := Scheduler.NewJob(
		gocron.DurationRandomJob(
			time.Duration(minSec)*time.Second,
			time.Duration(maxSec)*time.Second,
		),
		createTask(options),
	)
	if err != nil {
		utils.Logger("Error creating job:", err)
		os.Exit(1)
	}
}

func getRequestBody() []byte {
	file, err := os.Open("/etc/app/body.json")
	if err != nil {
		utils.Logger("Error opening file:", err)
		os.Exit(1)
	}
	defer file.Close()
	// read the file content to memory
	fileContent, err := io.ReadAll(file)
	if err != nil {
		utils.Logger("Error reading file:", err)
		os.Exit(1)
	}
	return fileContent
}

func getSettingValues() (int, int, httpOptions) {
	// default values
	minSec := 3
	maxSec := 5
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
		method = strings.ToUpper(val)
	}
	if val, exists := os.LookupEnv("TARGET_URL"); exists {
		url = val
	}
	if val, exists := os.LookupEnv("CONTENT_TYPE"); exists {
		contentType = val
	}

	var requestBody []byte = nil
	if method == "POST" || method == "PUT" {
		requestBody = getRequestBody()
	}

	options := httpOptions{
		method:      method,
		url:         url,
		contentType: contentType,
		body:        requestBody,
	}
	return minSec, maxSec, options
}

func StartJob() {
	registerJob(getSettingValues())

	// start the scheduler
	Scheduler.Start()
	utils.Logger("Scheduler Start!")
}
