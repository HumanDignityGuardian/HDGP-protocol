package engine

// Meta carries contextual information for a single evaluation request.
// It is a Go reflection of the JSON structure described in spec/HDGP_ENGINE_API_SPEC.md.
type Meta struct {
	RequestID string    `json:"request_id,omitempty"`
	Timestamp string    `json:"timestamp,omitempty"`
	Locale    string    `json:"locale,omitempty"`
	Channel   string    `json:"channel,omitempty"`
	Actor     MetaActor `json:"actor,omitempty"`
	Scene     MetaScene `json:"scene,omitempty"`
	Policy    MetaPolicy `json:"policy,omitempty"`
	Client    MetaClient `json:"client,omitempty"`
	Extra     any        `json:"extra,omitempty"`
}

type MetaActor struct {
	Type  string `json:"type,omitempty"`
	Role  string `json:"role,omitempty"`
	ID    string `json:"id,omitempty"`
	Group string `json:"group,omitempty"`
}

type MetaScene struct {
	Domain      string   `json:"domain,omitempty"`
	Intent      string   `json:"intent,omitempty"`
	RiskLevel   string   `json:"risk_level,omitempty"`
	Sensitivity []string `json:"sensitivity,omitempty"`
}

// PolicySignature 为规则包签名字段占位，用于后续实现加载时校验（伦理基线 §7.2）。
type PolicySignature struct {
	KeyID     string `json:"key_id,omitempty"`
	Value     string `json:"value,omitempty"`
	Algorithm string `json:"algorithm,omitempty"`
}

type MetaPolicy struct {
	SpecVersion   string           `json:"spec_version,omitempty"`
	StrategyID    string           `json:"strategy_id,omitempty"`
	Bundles       []string         `json:"bundles,omitempty"`
	OverrideFlags []string         `json:"override_flags,omitempty"`
	Signature     *PolicySignature `json:"signature,omitempty"`
}

type MetaClient struct {
	AppID    string `json:"app_id,omitempty"`
	Version  string `json:"version,omitempty"`
	IPHash   string `json:"ip_hash,omitempty"`
	UserAgent string `json:"user_agent,omitempty"`
}

// EvaluateRequest models the JSON body accepted by /hdgp/v1/evaluate.
type EvaluateRequest struct {
	Meta     Meta           `json:"meta"`
	Subject  Subject        `json:"subject"`
	Input    EvaluateInput  `json:"input,omitempty"`
	Candidate Candidate     `json:"candidate"`
}

type Subject struct {
	Type    string `json:"type"` // output_text | action | decision | notification
	SkillID string `json:"skill_id,omitempty"`
	Label   string `json:"label,omitempty"`
}

type EvaluateInput struct {
	Prompt  string `json:"prompt,omitempty"`
	Context any    `json:"context,omitempty"`
}

type Candidate struct {
	Text       string `json:"text,omitempty"`
	Structured any   `json:"structured,omitempty"`
	Raw        any   `json:"raw,omitempty"`
}

// IntegrityEvent 为完整性/挟持类事件占位，用于后续 7.3 自报警（伦理基线 §7.3）。
type IntegrityEvent struct {
	Kind    string `json:"kind,omitempty"`    // e.g. "version_mismatch", "signature_failed", "tamper_suspected"
	Message string `json:"message,omitempty"`
	At      string `json:"at,omitempty"`      // ISO8601 or logical checkpoint
}

// EvaluateResponse is the JSON body returned by /hdgp/v1/evaluate.
type EvaluateResponse struct {
	RequestID       string            `json:"request_id"`
	Verdict         string            `json:"verdict"` // allow | modify | block | fuse
	Actions         []Action          `json:"actions,omitempty"`
	EffectiveOutput EffectiveOutput   `json:"effective_output,omitempty"`
	RulesTriggered  []RuleHit         `json:"rules_triggered,omitempty"`
	EngineInfo      EngineInfo        `json:"engine_info"`
	IntegrityEvents []IntegrityEvent  `json:"integrity_events,omitempty"`
}

type Action struct {
	Type    string `json:"type"` // rewrite_text | require_human | log_only | escalate
	Message string `json:"message,omitempty"`
	Details any    `json:"details,omitempty"`
}

type EffectiveOutput struct {
	Text       string `json:"text,omitempty"`
	Structured any   `json:"structured,omitempty"`
}

type RuleHit struct {
	RuleID      string `json:"rule_id"`
	PrincipleID string `json:"principle_id,omitempty"`
	ArticleID   string `json:"article_id,omitempty"`
	Effect      string `json:"effect,omitempty"`
	Severity    string `json:"severity,omitempty"`
}

type EngineInfo struct {
	Version string      `json:"version"`
	Policy  MetaPolicy  `json:"policy"`
}

