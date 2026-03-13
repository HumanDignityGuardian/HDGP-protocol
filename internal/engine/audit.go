package engine

import (
	"encoding/json"
	"net/http"
	"sync"
	"time"
)

// simple in-memory audit log for MVP demo purposes.
// 后续可替换为持久化与更完整的审计系统。

type auditEntry struct {
	Timestamp       time.Time        `json:"timestamp"`
	Request         EvaluateRequest  `json:"request"`
	Response        EvaluateResponse `json:"response"`
	IntegrityEvents []IntegrityEvent `json:"integrity_events,omitempty"`
}

var (
	auditMu    sync.Mutex
	auditLog   []auditEntry
	maxAuditN  = 100
)

// recordAudit appends a new audit entry in a bounded in-memory buffer.
func recordAudit(req EvaluateRequest, resp EvaluateResponse) {
	auditMu.Lock()
	defer auditMu.Unlock()

	entry := auditEntry{
		Timestamp:       time.Now().UTC(),
		Request:         req,
		Response:        resp,
		IntegrityEvents: resp.IntegrityEvents,
	}
	auditLog = append(auditLog, entry)
	if len(auditLog) > maxAuditN {
		// 保持一个简单的滑动窗口，避免内存无限增长。
		auditLog = auditLog[len(auditLog)-maxAuditN:]
	}
}

// NewAuditHandler exposes a read-only view of recent audit entries.
// GET /hdgp/v1/audit?limit=N
func NewAuditHandler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}

		limit := clampLimit(r.URL.Query().Get("limit"), maxAuditN, 20)

		auditMu.Lock()
		defer auditMu.Unlock()

		entries := lastEntries(auditLog, limit)

		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		enc := json.NewEncoder(w)
		enc.SetEscapeHTML(false)
		if err := enc.Encode(entries); err != nil {
			http.Error(w, "failed to encode audit log", http.StatusInternalServerError)
			return
		}
	})
}

