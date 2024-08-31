package handlers

import (
	"encoding/json"
	"net/http"
	"os"
	"time"

	"github.com/ishikawa096/auto-requester/scheduler"
	"github.com/ishikawa096/auto-requester/utils"
)

// Stop job handler
func StopHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	response := map[string]string{"message": "job stop! To restart, please access /start"}
	json.NewEncoder(w).Encode(response)

	err := scheduler.Scheduler.StopJobs()
	if err != nil {
		utils.Logger("Error stoping jobs:", err)
		os.Exit(1)
	} else {
		utils.Logger(time.Now(), "Scheduler Stop Jobs!")
	}
}
