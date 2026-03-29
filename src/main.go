package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strconv"
)

func (lm *LockManager) Unlock(path string) (int, int) {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return StatusIdle, http.StatusNotFound
	}

	lm.mu.Lock()
	defer lm.mu.Unlock()

	currentStatus, exists := lm.locks[path]
	if !exists || currentStatus == StatusIdle {
		return StatusIdle, http.StatusNoContent
	}

	// Unlock (set to Idle)
	lm.locks[path] = StatusIdle
	return StatusIdle, http.StatusOK
}

func main() {
	lm := NewLockManager()

	http.HandleFunc("/lock", func(w http.ResponseWriter, r *http.Request) {
		path := r.URL.Query().Get("path")
		statusStr := r.URL.Query().Get("status")

		if path == "" || statusStr == "" {
			http.Error(w, "Missing path or status", http.StatusBadRequest)
			return
		}

		status, err := strconv.Atoi(statusStr)
		if err != nil || status < 1 || status > 5 {
			http.Error(w, "Invalid status", http.StatusBadRequest)
			return
		}

		newStatus, httpStatus := lm.Lock(path, status)

		if httpStatus == http.StatusNotFound {
			w.WriteHeader(http.StatusNotFound)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(httpStatus)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"status": newStatus,
		})
	})

	http.HandleFunc("/unlock", func(w http.ResponseWriter, r *http.Request) {
		path := r.URL.Query().Get("path")
		if path == "" {
			http.Error(w, "Missing path", http.StatusBadRequest)
			return
		}

		newStatus, httpStatus := lm.Unlock(path)

		if httpStatus == http.StatusNotFound {
			w.WriteHeader(http.StatusNotFound)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(httpStatus)
		if httpStatus == http.StatusOK {
			json.NewEncoder(w).Encode(newStatus)
		}
	})

	fmt.Println("Server starting on :8080...")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		fmt.Printf("Server failed: %s\n", err)
	}
}
