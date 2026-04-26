package handler

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/google/uuid"
)

func PayHandler(w http.ResponseWriter, r *http.Request) {

	time.Sleep(2 * time.Second)

	resp := map[string]interface{}{
		"status":         "paid",
		"amount":         1000,
		"transaction_id": uuid.New().String(),
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(resp)
}
