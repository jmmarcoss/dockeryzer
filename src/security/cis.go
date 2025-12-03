package security

type CISResult struct {
	RuleID      string
	Description string
	Passed      bool
	Severity    string
	Message     string
}

type CISRule interface {
	Check(dockerfile string) CISResult
}

type CISAnalyzer struct {
	rules []CISRule
}

func NewCISAnalyzer() *CISAnalyzer {
	return &CISAnalyzer{
		rules: []CISRule{
			OfficialBaseImageRule{},
			ExplicitTagRule{},
			NoRootUserRule{},
			CleanCacheRule{},
			HealthcheckRule{},
			DockerIgnoreRule{},
			MinimalPortExposureRule{},
			MultiStageBuildRule{},
			CombinedRunCommandRule{},
			OptimizedOrderRule{},
		},
	}
}

func (a *CISAnalyzer) Analyze(content string) []CISResult {
	results := []CISResult{}
	for _, rule := range a.rules {
		results = append(results, rule.Check(content))
	}
	return results
}
