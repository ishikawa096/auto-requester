package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/ishikawa096/auto-requester/scheduler"
	"github.com/ishikawa096/auto-requester/utils"
)

// Start job handler
func StartHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	response := map[string]string{"message": "Scheduler Start! To stop, please access /stop"}
	json.NewEncoder(w).Encode(response)

	scheduler.Scheduler.Start()
	utils.Logger(utils.Green("Scheduler Restart!✨"))
}
