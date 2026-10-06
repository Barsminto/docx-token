package sample

import "time"

// FillMap demo data for docxtoken.Fill (generic paths + GI sample issues).
func FillMap(now time.Time) map[string]interface{} {
	issues := make([]map[string]interface{}, 0, len(Items()))
	for _, item := range Items() {
		issues = append(issues, map[string]interface{}{
			"title":   item.Title,
			"brief":   item.ConcernBrief,
			"summary": item.ConcernSummary,
		})
	}
	return map[string]interface{}{
		"projectName":          "XXX",
		"draftDate":            now.Format("01 February 2006"),
		"accountableExecutive": "[to be left blank for the Investigator to complete]",
		"subjectList":          "list all subject names, job title and PSID",
		"background.period":    "February 2026",
		"case": map[string]interface{}{
			"name":   "Sample investigation case",
			"issues": issues,
		},
		"executiveSummary": map[string]interface{}{
			"identified":       "...",
			"concernRef":       "...",
			"substantiation":   "not substantiated",
			"evidence":         "did not establish",
			"recommendations":  "...",
		},
		"conclusion": map[string]interface{}{
			"support":         "Colleague A and Colleague B may have breached Bank policies",
			"policyBreaches":  "...",
		},
		"recommendation": map[string]interface{}{
			"1": "...",
			"2": "...",
			"3": "...",
		},
		"Report.GeneratedAt": now.Format("2006-01-02 15:04"),
	}
}
