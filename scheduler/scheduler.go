package scheduler

import (
	"bytes"
	"io"
	"net/http"
	"os"
	"time"

	"github.com/go-co-op/gocron/v2"
	"github.com/ishikawa096/auto-requester/utils"
)

type requestOptions struct {
	method      string // GET, POST, PUT, DELETE
	url         string
	contentType string
	body        []byte
	randomize   bool
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

func sendRequest(options requestOptions) (*http.Response, error) {
	var requestBody []byte = nil
	if options.randomize {
		requestBody = selectRandomElement(options.body)
	} else {
		requestBody = options.body
	}

	// To use the request body multiple times, we need to create a reader
	bodyReader := bytes.NewReader(requestBody)

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

func createTask(options requestOptions) gocron.Task {
	return gocron.NewTask(
		func(options requestOptions) {
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

func registerJob(minSec, maxSec int, options requestOptions) {
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

func StartJob() {
	registerJob(getConfigs())

	// start the scheduler
	Scheduler.Start()
	utils.Logger("Scheduler Start!")
}
