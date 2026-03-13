package engine

import (
	"encoding/json"
	"net/http"
	"sync"
	"time"
)

// MVP 版上诉记录，仅在内存中保存，便于后续扩展为完整工作流。

type AppealRequest struct {
	RequestID   string `json:"request_id,omitempty"`
	Reason      string `json:"reason,omitempty"`
	Contact     string `json:"contact,omitempty"`
	SubmittedBy string `json:"submitted_by,omitempty"`
}

type AppealRecord struct {
	Timestamp time.Time     `json:"timestamp"`
	Appeal    AppealRequest `json:"appeal"`
}

var (
	appealMu   sync.Mutex
	appealLog  []AppealRecord
	maxAppeals = 200
)

func NewAppealHandler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodPost:
			handleAppealPost(w, r)
		case http.MethodGet:
			handleAppealList(w, r)
		default:
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		}
	})
}

func handleAppealPost(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()

	var ar AppealRequest
	if err := decoder.Decode(&ar); err != nil {
		http.Error(w, "invalid JSON: "+err.Error(), http.StatusBadRequest)
		return
	}

	appealMu.Lock()
	defer appealMu.Unlock()

	rec := AppealRecord{
		Timestamp: time.Now().UTC(),
		Appeal:    ar,
	}
	appealLog = append(appealLog, rec)
	if len(appealLog) > maxAppeals {
		appealLog = appealLog[len(appealLog)-maxAppeals:]
	}

	w.WriteHeader(http.StatusAccepted)
}

func handleAppealList(w http.ResponseWriter, r *http.Request) {
	limit := clampLimit(r.URL.Query().Get("limit"), maxAppeals, 20)

	appealMu.Lock()
	defer appealMu.Unlock()

	entries := lastEntries(appealLog, limit)

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	enc := json.NewEncoder(w)
	enc.SetEscapeHTML(false)
	if err := enc.Encode(entries); err != nil {
		http.Error(w, "failed to encode appeals", http.StatusInternalServerError)
		return
	}
}

