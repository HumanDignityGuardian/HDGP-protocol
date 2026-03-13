package engine

import (
	"encoding/json"
	"net/http"
)

// StatusResponse is returned by GET /hdgp/v1/status for health and version exposure.
// Used to support "运行版本与规则包版本的对外暴露" (HDGP_ETHICS_BASELINE §7, ENGINE_API_SPEC).
type StatusResponse struct {
	EngineVersion string       `json:"engine_version"`
	SpecVersion   string       `json:"spec_version"`
	Policy        PolicyStatus `json:"policy"`
}

// PolicyStatus describes the loaded policy/strategy for audit and integrity checks.
type PolicyStatus struct {
	SpecVersion  string   `json:"spec_version"`
	StrategyID   string   `json:"strategy_id"`
	Bundles      []string `json:"bundles"`
}

const (
	statusSpecVersion  = "HDGP-1.0"
	statusStrategyID   = "S-global-default"
	statusBundleCore  = "B-CORE-1.0.0"
)

// NewStatusHandler returns an HTTP handler for GET /hdgp/v1/status.
// Exposes engine version and policy version for monitoring and "self-alert" capability
// when deployed version or signature does not match expected (see ethics baseline §7.3).
func NewStatusHandler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}

		resp := StatusResponse{
			EngineVersion: engineVersion,
			SpecVersion:   statusSpecVersion,
			Policy: PolicyStatus{
				SpecVersion: statusSpecVersion,
				StrategyID:  statusStrategyID,
				Bundles:     []string{statusBundleCore},
			},
		}

		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		enc := json.NewEncoder(w)
		enc.SetEscapeHTML(false)
		_ = enc.Encode(resp)
	})
}
