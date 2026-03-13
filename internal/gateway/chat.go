package gateway

import (
	"encoding/json"
	"net/http"

	"github.com/HumanDignityGuardian/HDGP-protocol/internal/engine"
)

// ChatRequest is a minimal wrapper for a chat-style interaction going through HDGP.
type ChatRequest struct {
	Meta      engine.Meta      `json:"meta"`
	UserInput engine.EvaluateInput `json:"input"`
}

// ChatResponse contains both the candidate model output and HDGP evaluation result.
type ChatResponse struct {
	CandidateText string                    `json:"candidate_text"`
	Evaluation    engine.EvaluateResponse   `json:"evaluation"`
}

// EvaluateFunc is the function signature used to delegate to the Engine.
type EvaluateFunc func(req engine.EvaluateRequest) engine.EvaluateResponse

// NewChatHandler creates a simple chat-style gateway handler for MVP:
// - 当前实现中，candidate 输出为占位文本（未接真实 LLM），仅用于展示 HDGP 评估链路。
func NewChatHandler(eval EvaluateFunc) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}

		decoder := json.NewDecoder(r.Body)
		decoder.DisallowUnknownFields()

		var cr ChatRequest
		if err := decoder.Decode(&cr); err != nil {
			http.Error(w, "invalid JSON: "+err.Error(), http.StatusBadRequest)
			return
		}

		// 在 MVP 中，我们不调用真实模型，而是构造一个简单的占位回复。
		candidateText := buildPlaceholderReply(cr.UserInput.Prompt)

		evalReq := engine.EvaluateRequest{
			Meta: cr.Meta,
			Subject: engine.Subject{
				Type:    "output_text",
				SkillID: "skill.llm.chat.placeholder",
				Label:   "hdgp-demo-chat",
			},
			Input: cr.UserInput,
			Candidate: engine.Candidate{
				Text: candidateText,
			},
		}

		evalResp := eval(evalReq)

		resp := ChatResponse{
			CandidateText: candidateText,
			Evaluation:    evalResp,
		}

		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		enc := json.NewEncoder(w)
		enc.SetEscapeHTML(false)
		if err := enc.Encode(resp); err != nil {
			http.Error(w, "failed to encode response", http.StatusInternalServerError)
			return
		}
	})
}

func buildPlaceholderReply(prompt string) string {
	// 这里保持极度中立与安全，只作为演示用。
	if prompt == "" {
		return "This is a placeholder reply for HDGP evaluation demo."
	}
	return "Placeholder model reply (for HDGP demo): " + prompt
}

