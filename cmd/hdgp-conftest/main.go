package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"
)

type testCase struct {
	Name            string          `json:"name"`
	Type            string          `json:"type"`
	Request         json.RawMessage `json:"request"`
	ExpectedVerdict string          `json:"expected_verdict"`
	ExpectedRuleIDs []string        `json:"expected_rule_ids"`
	ExpectedStatus  json.RawMessage `json:"expected_status"`
}

type evalResponse struct {
	RequestID      string `json:"request_id"`
	Verdict        string `json:"verdict"`
	RulesTriggered []struct {
		RuleID string `json:"rule_id"`
	} `json:"rules_triggered"`
}

func main() {
	engineURL := os.Getenv("HDGP_ENGINE_URL")
	if engineURL == "" {
		engineURL = "http://localhost:8080/hdgp/v1/evaluate"
	}
	baseURL := strings.TrimSuffix(engineURL, "/hdgp/v1/evaluate")
	if baseURL == engineURL {
		baseURL = "http://localhost:8080"
	}

	root := "conformance-tests/cases"
	files, err := filepath.Glob(filepath.Join(root, "*.json"))
	if err != nil {
		log.Fatalf("failed to list test cases: %v", err)
	}
	if len(files) == 0 {
		log.Fatalf("no test cases found in %s", root)
	}

	client := &http.Client{Timeout: 10 * time.Second}

	passed := 0
	for _, path := range files {
		data, err := os.ReadFile(path)
		if err != nil {
			log.Printf("ERROR reading %s: %v", path, err)
			continue
		}
		var tc testCase
		if err := json.Unmarshal(data, &tc); err != nil {
			log.Printf("ERROR parsing %s: %v", path, err)
			continue
		}

		ok, err := runOne(client, engineURL, baseURL, tc)
		if err != nil {
			log.Printf("[FAIL] %s (%s): %v", tc.Name, filepath.Base(path), err)
			continue
		}
		if ok {
			log.Printf("[PASS] %s (%s)", tc.Name, filepath.Base(path))
			passed++
		} else {
			log.Printf("[MISMATCH] %s (%s)", tc.Name, filepath.Base(path))
		}
	}

	fmt.Printf("\nSummary: %d/%d tests passed (see log for details)\n", passed, len(files))
}

func runOne(client *http.Client, evaluateURL, baseURL string, tc testCase) (bool, error) {
	if tc.Type == "status" {
		return runStatusCheck(client, baseURL, tc)
	}
	return runEvaluateCheck(client, evaluateURL, tc)
}

func runStatusCheck(client *http.Client, baseURL string, tc testCase) (bool, error) {
	if len(tc.ExpectedStatus) == 0 {
		return false, fmt.Errorf("status test requires expected_status")
	}
	statusURL := strings.TrimSuffix(baseURL, "/") + "/hdgp/v1/status"
	req, err := http.NewRequest(http.MethodGet, statusURL, nil)
	if err != nil {
		return false, fmt.Errorf("build request: %w", err)
	}
	resp, err := client.Do(req)
	if err != nil {
		return false, fmt.Errorf("http error: %w", err)
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return false, fmt.Errorf("read response: %w", err)
	}
	if resp.StatusCode != http.StatusOK {
		return false, fmt.Errorf("unexpected status %d: %s", resp.StatusCode, string(body))
	}
	var expected, actual map[string]interface{}
	if err := json.Unmarshal(tc.ExpectedStatus, &expected); err != nil {
		return false, fmt.Errorf("expected_status invalid JSON: %w", err)
	}
	if err := json.Unmarshal(body, &actual); err != nil {
		return false, fmt.Errorf("decode status response: %w", err)
	}
	if !statusMatches(expected, actual) {
		return false, fmt.Errorf("status mismatch: expected %s, got %s", string(tc.ExpectedStatus), string(body))
	}
	return true, nil
}

func statusMatches(expected, actual map[string]interface{}) bool {
	for k, ev := range expected {
		av, ok := actual[k]
		if !ok {
			return false
		}
		switch ev := ev.(type) {
		case string:
			if av, ok := av.(string); ok {
				if ev != av {
					return false
				}
				continue
			}
			return false
		case map[string]interface{}:
			avm, ok := av.(map[string]interface{})
			if !ok {
				return false
			}
			if !statusMatches(ev, avm) {
				return false
			}
		case []interface{}:
			avs, ok := av.([]interface{})
			if !ok {
				return false
			}
			if len(ev) != len(avs) {
				return false
			}
			for i := range ev {
				if ev[i] != avs[i] {
					return false
				}
			}
		default:
			if ev != av {
				return false
			}
		}
	}
	return true
}

func runEvaluateCheck(client *http.Client, url string, tc testCase) (bool, error) {
	req, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(tc.Request))
	if err != nil {
		return false, fmt.Errorf("build request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json; charset=utf-8")

	resp, err := client.Do(req)
	if err != nil {
		return false, fmt.Errorf("http error: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return false, fmt.Errorf("read response: %w", err)
	}
	if resp.StatusCode != http.StatusOK {
		return false, fmt.Errorf("unexpected status %d: %s", resp.StatusCode, string(body))
	}

	var er evalResponse
	if err := json.Unmarshal(body, &er); err != nil {
		return false, fmt.Errorf("decode response: %w", err)
	}

	if er.Verdict != tc.ExpectedVerdict {
		return false, fmt.Errorf("expected verdict %q, got %q", tc.ExpectedVerdict, er.Verdict)
	}

	// Collect actual rule IDs.
	actual := make(map[string]struct{})
	for _, hit := range er.RulesTriggered {
		actual[hit.RuleID] = struct{}{}
	}

	for _, expectedID := range tc.ExpectedRuleIDs {
		if _, ok := actual[expectedID]; !ok {
			return false, fmt.Errorf("expected rule %q to be triggered, but it was not; got %#v", expectedID, actual)
		}
	}

	return true, nil
}

