package scheduler

import (
	"encoding/json"
	"time"

	"github.com/ishikawa096/auto-requester/utils"
	"golang.org/x/exp/rand"
)

// IF the body is a JSON array, select a random element from it.
// Otherwise, return the body as is.
func selectRandomElement(body []byte) []byte {
	var jsonArray []interface{}
	if err := json.Unmarshal(body, &jsonArray); err == nil {
		if len(jsonArray) > 0 {
			rand.Seed(uint64(time.Now().UnixNano()))
			randomIndex := rand.Intn(len(jsonArray))
			randomElement, err := json.Marshal(jsonArray[randomIndex])
			if err != nil {
				utils.Logger("failed to marshal random element: %v", err)
				return body
			}
			return randomElement
		}
	}
	return body
}
