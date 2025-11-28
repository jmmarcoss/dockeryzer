package security

import "fmt"

func PrintCISResults(results []CISResult) {
	fmt.Println("\nSecurity Analysis based on CIS Docker Benchmark:\n")

	score := 0

	for _, r := range results {
		status := "PASS"
		if !r.Passed {
			status = "FAIL"
		} else {
			score++
		}

		fmt.Printf("[%s] %s - %s\n", status, r.RuleID, r.Description)
		if !r.Passed {
			fmt.Printf("  Severity: %s\n  Issue: %s\n\n", r.Severity, r.Message)
		}
	}

	percent := (score * 100) / len(results)
	fmt.Printf("Security Score: %d%%\n", percent)
}
