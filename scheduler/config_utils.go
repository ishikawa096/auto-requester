package scheduler

import (
	"io"
	"os"
	"strconv"
	"strings"

	"github.com/ishikawa096/auto-requester/utils"
)

// default env values
const (
	defaultMinSec    = 3
	defaultMaxSec    = 5
	defaultUrl       = "http://localhost:3000"
	defaultMethod    = "GET"
	defaultType      = "application/json"
	defaultRandomize = true
	defaultFilePath  = "/etc/app/body.json"
)

func getStrEnv(envKey string, defaultValue string) string {
	if value, exists := os.LookupEnv(envKey); exists {
		return value
	}
	return defaultValue
}

func getIntEnv(envKey string, defaultValue int) int {
	valueStr, exists := os.LookupEnv(envKey)
	if !exists {
		return defaultValue
	}
	value, err := strconv.Atoi(valueStr)
	if err != nil {
		utils.Logger("🚨 invalid value for", utils.Yellow(envKey), err)
		return defaultValue
	}
	return value
}

func getBoolEnv(envKey string, defaultValue bool) bool {
	valueStr, exists := os.LookupEnv(envKey)
	if !exists {
		return defaultValue
	}
	value, err := strconv.ParseBool(valueStr)
	if err != nil {
		utils.Logger("🚨 invalid value for", utils.Yellow(envKey), err)
		return defaultValue
	}
	return value
}

// Reads the JSON file and returns its content as a byte slice.
func getRequestBody() []byte {
	filePath := getStrEnv("FILE_PATH", defaultFilePath)
	file, err := os.Open(filePath)
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

func getConfigs() (int, int, requestOptions) {
	minSec := getIntEnv("INTERVAL_MIN_SEC", defaultMinSec)
	maxSec := getIntEnv("INTERVAL_MAX_SEC", defaultMaxSec)
	url := getStrEnv("TARGET_URL", defaultUrl)
	method := strings.ToUpper(getStrEnv("HTTP_METHOD", defaultMethod))
	contentType := getStrEnv("CONTENT_TYPE", defaultType)
	randomize := getBoolEnv("RANDOMIZE", defaultRandomize)

	var requestBody []byte = nil
	if method == "POST" || method == "PUT" {
		requestBody = getRequestBody()
	}

	options := requestOptions{
		method:      method,
		url:         url,
		contentType: contentType,
		body:        requestBody,
		randomize:   randomize,
	}
	return minSec, maxSec, options
}
