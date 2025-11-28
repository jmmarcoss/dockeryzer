package security

import "strings"

// CIS-4.1: Container should not run as root
type NoRootUserRule struct{}

func (r NoRootUserRule) Check(df string) CISResult {
	if !strings.Contains(strings.ToLower(df), "user") {
		return CISResult{
			RuleID:      "CIS-4.1",
			Description: "Container should not run as root",
			Passed:      false,
			Severity:    "HIGH",
			Message:     "Missing USER instruction",
		}
	}
	return CISResult{RuleID: "CIS-4.1", Passed: true}
}

// CIS-4.5: Avoid latest tag
type NoLatestTagRule struct{}

func (r NoLatestTagRule) Check(df string) CISResult {
	if strings.Contains(df, ":latest") {
		return CISResult{
			RuleID:      "CIS-4.5",
			Description: "Avoid latest tag",
			Passed:      false,
			Severity:    "MEDIUM",
			Message:     "Image uses latest tag",
		}
	}
	return CISResult{RuleID: "CIS-4.5", Passed: true}
}

// CIS-4.6: HEALTHCHECK
type HealthcheckRule struct{}

func (r HealthcheckRule) Check(df string) CISResult {
	if !strings.Contains(strings.ToLower(df), "healthcheck") {
		return CISResult{
			RuleID:      "CIS-4.6",
			Description: "Container must define HEALTHCHECK",
			Passed:      false,
			Severity:    "LOW",
			Message:     "Missing HEALTHCHECK instruction",
		}
	}
	return CISResult{RuleID: "CIS-4.6", Passed: true}
}
