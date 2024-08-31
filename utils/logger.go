package utils

import (
	"fmt"
	"strings"
	"time"
)

func currentTimeStr() string {
	tz, err := time.LoadLocation("Asia/Tokyo")
	if err != nil {
		return ""
	}
	return time.Now().In(tz).Format("2006-01-02 15:04:05 MST")
}

func Logger(msgs ...interface{}) {
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
