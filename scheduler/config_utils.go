package scheduler

import (
	"io"
	"os"
	"strconv"
	"strings"

	"github.com/ishikawa096/auto-requester/utils"
)

func getRequestBody() []byte {
	file, err := os.Open("/etc/app/body.json")
	if err != nil {
		// if the file does not exist, return nil
		utils.Logger("Error opening file:", err)
		return nil
	}
	defer file.Close()
	// read the file content to memory
	fileContent, err := io.ReadAll(file)
	if err != nil {
		utils.Logger("Error reading file:", err)
		return nil
	}
	return fileContent
}

// Check the RANDOMIZE env and convert it to a boolean value
func getRandomizeEnv() bool {
	randomize := false
	if randomizeStr, exists := os.LookupEnv("RANDOMIZE"); exists {
		var err error
		randomize, err = strconv.ParseBool(randomizeStr)
		if err != nil {
			utils.Logger("invalid value for RANDOMIZE:", err)
		}
	}
	return randomize
}

func getConfigs() (int, int, requestOptions) {
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

	var randomize = getRandomizeEnv()

	options := requestOptions{
		method:      method,
		url:         url,
		contentType: contentType,
		body:        requestBody,
		randomize:   randomize,
	}
	return minSec, maxSec, options
}
