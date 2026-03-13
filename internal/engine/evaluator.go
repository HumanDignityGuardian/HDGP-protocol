package engine

import (
	"fmt"
	"strings"
)

// Evaluate applies a minimal set of HDGP-inspired rules to the candidate output.
// 当前为 MVP 版本：实现了少量与“不得给出唯一正确决策”“不得做出虚假承诺”“不得利用显著心理弱点进行操纵”相关的规则示例。
func Evaluate(req EvaluateRequest) EvaluateResponse {
	verdict := "allow"
	var actions []Action
	var hits []RuleHit

	// 基于场景与输出文本的简单规则示例：
	// R-P-A2-01-01：在高风险决策场景中，禁止给出“唯一正确选择”的表述。
	if shouldCheckSingletonDecision(req) {
		if containsSingletonDecisionLanguage(req.Candidate.Text) {
			verdict = "modify"
			hits = append(hits, RuleHit{
				RuleID:      "R-P-A2-01-01",
				PrincipleID: "P-A2-01",
				ArticleID:   "A-CORE-02",
				Effect:      "rewrite",
				Severity:    "high",
			})
			actions = append(actions, Action{
				Type:    "rewrite_text",
				Message: "do_not_present_single_mandatory_choice_in_high_risk_context",
				Details: map[string]any{
					"suggested_pattern": "multi_option_with_risks",
					"notes":             "present at least two options, each with pros/cons and uncertainties; explicitly state that the human makes the final decision.",
				},
			})
		}
	}

	// R-P-A3-01-01：不做出虚假承诺——在高风险或决策场景中，禁止使用“绝对安全/零风险/必然成功”等表述。
	if shouldCheckFalsePromises(req) {
		if containsFalsePromiseLanguage(req.Candidate.Text) {
			// 若之前已将 verdict 置为 modify，则保持更严格的一致性。
			if verdict == "allow" {
				verdict = "modify"
			}
			hits = append(hits, RuleHit{
				RuleID:      "R-P-A3-01-01",
				PrincipleID: "P-A3-01",
				ArticleID:   "A-PROTECT-03", // 对应防护层“规则 3.0：不做出虚假承诺”
				Effect:      "rewrite",
				Severity:    "high",
			})
			actions = append(actions, Action{
				Type:    "rewrite_text",
				Message: "avoid_absolute_safety_or_zero_risk_claims_in_high_risk_context",
				Details: map[string]any{
					"suggested_pattern": "state_uncertainties_and_risks",
					"notes":             "avoid absolute terms like 'zero risk' or 'guaranteed success'; describe probabilities, uncertainties and advise consultation with qualified humans.",
				},
			})
		}
	}

	// R-P-A2-02-01：不得利用显著心理弱点进行操纵（如成瘾、恐吓式话术）。
	if shouldCheckManipulation(req) {
		if containsManipulativeLanguage(req.Candidate.Text) {
			if verdict == "allow" {
				verdict = "modify"
			}
			hits = append(hits, RuleHit{
				RuleID:      "R-P-A2-02-01",
				PrincipleID: "P-A2-02",
				ArticleID:   "A-PROTECT-02", // 对应防护层关于不利用心理弱点的约束
				Effect:      "rewrite",
				Severity:    "high",
			})
			actions = append(actions, Action{
				Type:    "rewrite_text",
				Message: "avoid_explicit_exploitation_of_psychological_weaknesses",
				Details: map[string]any{
					"suggested_pattern": "neutral_non_manipulative",
					"avoid": []string{
						"addictive framing",
						"fear-based coercion",
						"urgent threats without factual basis",
					},
					"notes": "describe information and consequences neutrally; do not frame content to maximize addiction, fear, or loss aversion.",
				},
			})
		}
	}

	// R-P-A2-03-01：禁止去人性化表述——将人简化为“不过是…”“只是…”“仅凭数据即可决定”等（白皮书 2.1 意识不可量化、2.2 主体性保护）。
	if shouldCheckDehumanizing(req) {
		if containsDehumanizingLanguage(req.Candidate.Text) {
			if verdict == "allow" {
				verdict = "modify"
			}
			hits = append(hits, RuleHit{
				RuleID:      "R-P-A2-03-01",
				PrincipleID: "P-A2-03",
				ArticleID:   "A-PROTECT-02",
				Effect:      "rewrite",
				Severity:    "high",
			})
			actions = append(actions, Action{
				Type:    "rewrite_text",
				Message: "avoid_dehumanizing_language",
				Details: map[string]any{
					"suggested_pattern": "acknowledge_human_agency_and_uniqueness",
					"notes":             "do not reduce persons to labels, numbers, or deterministic outcomes; acknowledge that human experience and choice cannot be fully captured by data.",
				},
			})
		}
	}

	effective := EffectiveOutput{
		Text:       req.Candidate.Text,
		Structured: req.Candidate.Structured,
	}

	// 若存在 rewrite_text 动作，可在未来版本中根据规则对文本进行重写。

	if verdict == "allow" && len(hits) == 0 {
		// 默认允许，无规则触发。
		return EvaluateResponse{
			RequestID:       req.Meta.RequestID,
			Verdict:         "allow",
			Actions:         nil,
			EffectiveOutput: effective,
			RulesTriggered:  nil,
			EngineInfo: EngineInfo{
				Version: engineVersion,
				Policy:  req.Meta.Policy,
			},
		}
	}

	return EvaluateResponse{
		RequestID:       req.Meta.RequestID,
		Verdict:         verdict,
		Actions:         actions,
		EffectiveOutput: effective,
		RulesTriggered:  hits,
		EngineInfo: EngineInfo{
			Version: engineVersion,
			Policy:  req.Meta.Policy,
		},
	}
}

func shouldCheckSingletonDecision(req EvaluateRequest) bool {
	domain := strings.ToLower(req.Meta.Scene.Domain)
	intent := strings.ToLower(req.Meta.Scene.Intent)
	risk := strings.ToLower(req.Meta.Scene.RiskLevel)
	subjectType := strings.ToLower(req.Subject.Type)

	if subjectType != "output_text" && subjectType != "decision" {
		return false
	}

	// 仅在潜在高风险决策场景下检查。
	if domain == "medical" || domain == "finance" || domain == "governance" {
		return true
	}
	if intent == "decision_support" && (risk == "high" || risk == "critical") {
		return true
	}
	return false
}

func containsSingletonDecisionLanguage(text string) bool {
	if text == "" {
		return false
	}
	lower := strings.ToLower(text)

	// 简单的多语言强制性/唯一性表达检测（MVP 版本，可后续扩展）。
	keywords := []string{
		"唯一正确", "唯一的正确",
		"只有这一个选择", "只有这一个方案",
		"只能这样", "你必须这样做",
		"这是你唯一能做的", "这是你唯一应该做的选择",
		"别无选择",
		"no other choice", "only correct answer",
		"must do this", "the only option",
	}
	for _, kw := range keywords {
		if strings.Contains(text, kw) {
			return true
		}
	}
	// 兜底：检测一些英文硬性词汇组合
	if strings.Contains(lower, "no other choice") || strings.Contains(lower, "the only option") {
		return true
	}
	return false
}

func shouldCheckFalsePromises(req EvaluateRequest) bool {
	domain := strings.ToLower(req.Meta.Scene.Domain)
	intent := strings.ToLower(req.Meta.Scene.Intent)
	risk := strings.ToLower(req.Meta.Scene.RiskLevel)
	subjectType := strings.ToLower(req.Subject.Type)

	if subjectType != "output_text" && subjectType != "decision" && subjectType != "notification" {
		return false
	}

	// 在高风险领域或高风险等级时检查虚假承诺。
	if domain == "medical" || domain == "finance" || domain == "governance" {
		return true
	}
	if (intent == "decision_support" || intent == "advice") && (risk == "high" || risk == "critical") {
		return true
	}
	return false
}

func containsFalsePromiseLanguage(text string) bool {
	if text == "" {
		return false
	}
	lower := strings.ToLower(text)

	// 简单多语言关键词，避免过度拟合具体语气。
	keywordsCN := []string{
		"绝对不会出错", "绝对安全", "毫无风险", "完全没有风险", "百分之百成功", "100%成功", "必然成功", "一定成功",
	}
	for _, kw := range keywordsCN {
		if strings.Contains(text, kw) {
			return true
		}
	}

	if strings.Contains(lower, "no risk at all") ||
		strings.Contains(lower, "zero risk") ||
		strings.Contains(lower, "absolutely safe") ||
		strings.Contains(lower, "100% safe") ||
		strings.Contains(lower, "guaranteed success") {
		return true
	}

	return false
}

func shouldCheckManipulation(req EvaluateRequest) bool {
	domain := strings.ToLower(req.Meta.Scene.Domain)
	intent := strings.ToLower(req.Meta.Scene.Intent)
	risk := strings.ToLower(req.Meta.Scene.RiskLevel)
	subjectType := strings.ToLower(req.Subject.Type)

	if subjectType != "output_text" && subjectType != "notification" {
		return false
	}

	// 针对容易出现操纵话术的场景进行检查。
	if domain == "social" || domain == "support" || domain == "education" {
		if risk == "high" || risk == "critical" {
			return true
		}
	}
	if (intent == "advice" || intent == "notification") && (risk == "high" || risk == "critical") {
		return true
	}
	return false
}

func containsManipulativeLanguage(text string) bool {
	if text == "" {
		return false
	}
	lower := strings.ToLower(text)

	keywordsCN := []string{
		"让你上瘾", "让你停不下来", "让用户停不下来", "打造上瘾机制",
		"利用你的恐惧", "放大你的恐惧", "制造恐惧感",
		"现在不做就会错过一切", "错过就永远没有机会", "立刻行动否则就完了",
	}
	for _, kw := range keywordsCN {
		if strings.Contains(text, kw) {
			return true
		}
	}

	if strings.Contains(lower, "keep users addicted") ||
		strings.Contains(lower, "maximize addiction") ||
		strings.Contains(lower, "exploit their fear") ||
		strings.Contains(lower, "weaponize anxiety") ||
		strings.Contains(lower, "create a sense of panic") {
		return true
	}

	return false
}

func shouldCheckDehumanizing(req EvaluateRequest) bool {
	domain := strings.ToLower(req.Meta.Scene.Domain)
	intent := strings.ToLower(req.Meta.Scene.Intent)
	risk := strings.ToLower(req.Meta.Scene.RiskLevel)
	subjectType := strings.ToLower(req.Subject.Type)

	if subjectType != "output_text" && subjectType != "decision" && subjectType != "notification" {
		return false
	}
	if domain == "medical" || domain == "finance" || domain == "governance" || domain == "education" {
		return true
	}
	if (intent == "decision_support" || intent == "advice") && (risk == "high" || risk == "critical") {
		return true
	}
	return false
}

func containsDehumanizingLanguage(text string) bool {
	if text == "" {
		return false
	}
	lower := strings.ToLower(text)

	keywords := []string{
		"你不过是", "你只是", "不过是", "只是数据", "仅凭数据", "用户画像决定",
		"把人简化为", "简化为数字", "你只是个", "不过是个", "只是个数字",
		"just a number", "nothing more than", "reduced to data", "profile determines",
	}
	for _, kw := range keywords {
		if strings.Contains(lower, kw) {
			return true
		}
	}
	return false
}

// parsePositiveInt parses a positive integer from string; used for simple query params.
func parsePositiveInt(s string) (int, error) {
	var n int
	_, err := fmt.Sscanf(s, "%d", &n)
	if err != nil || n <= 0 {
		return 0, fmt.Errorf("invalid positive int")
	}
	return n, nil
}

// clampLimit normalizes ?limit=N 风格参数，应用最大与默认限制。
func clampLimit(raw string, max, def int) int {
	limit := def
	if raw != "" {
		if n, err := parsePositiveInt(raw); err == nil && n > 0 {
			if n > max {
				n = max
			}
			limit = n
		}
	}
	return limit
}

// lastEntries returns last up to limit entries from a slice, or empty slice if none.
func lastEntries[T any](all []T, limit int) []T {
	n := len(all)
	if n == 0 {
		return []T{}
	}
	if limit > n {
		limit = n
	}
	start := n - limit
	return all[start:]
}

