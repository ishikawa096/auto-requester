package handlers

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestStopHandler(t *testing.T) {
	req, err := http.NewRequest("GET", "/stop", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(StopHandler)

	handler.ServeHTTP(rr, req)

	// check the status code
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	// check response body
	expected := map[string]string{"message": "Sucessfully stopped the job! To restart, please access /start"}
	var actual map[string]string
	if err := json.Unmarshal(rr.Body.Bytes(), &actual); err != nil {
		t.Fatalf("could not unmarshal response: %v", err)
	}

	if actual["message"] != expected["message"] {
		t.Errorf("handler returned unexpected body: got %v want %v",
			actual["message"], expected["message"])
	}
}
