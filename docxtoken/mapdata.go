package docxtoken

import (
	"fmt"
	"strings"
	"time"

	"github.com/liboyang/docx-token/internal/domain"
)

// RecordsKey is the preferred map key for concern rows used in template loops.
const RecordsKey = "records"

var recordListKeys = []string{RecordsKey, "items", "Records", "Items"}

func BuildFromMap(data map[string]interface{}) (domain.Document, error) {
	if data == nil {
		data = map[string]interface{}{}
	}
	flat := flattenMap("", data)

	items, err := parseRecordInputs(data)
	if err != nil {
		return domain.Document{}, err
	}

	generatedAt := time.Now()
	if s, ok := flat["Report.GeneratedAt"]; ok && s != "" {
		if t, err := parseTime(s); err == nil {
			generatedAt = t
		}
	}

	title := flat["projectName"]
	if title == "" {
		title = flat["Project.Name"]
	}
	if title == "" {
		title = flat["Report.Title"]
	}
	if title == "" {
		title = "PROJECT"
	}
	author := flat["accountableExecutive"]
	if author == "" {
		author = flat["Summary.Author"]
	}
	if author == "" {
		author = flat["Executive.Accountable"]
	}

	doc := domain.NewDocument(title, author, generatedAt, items)
	return doc, nil
}

func ScalarPlaceholders(data map[string]interface{}) map[string]string {
	out := flattenMap("", data)
	for _, key := range recordListKeys {
		delete(out, key)
	}
	for k := range out {
		if strings.HasPrefix(k, "Record.") || strings.HasPrefix(k, "Record_") ||
			strings.HasPrefix(k, "issue.") || strings.HasPrefix(k, "issue_") {
			delete(out, k)
		}
	}
	return out
}

func flattenMap(prefix string, data map[string]interface{}) map[string]string {
	out := make(map[string]string)
	for k, v := range data {
		key := k
		if prefix != "" {
			key = prefix + "." + k
		}
		switch val := v.(type) {
		case map[string]interface{}:
			for fk, fv := range flattenMap(key, val) {
				out[fk] = fv
			}
		case []interface{}, []map[string]interface{}:
			continue
		default:
			out[key] = stringify(val)
		}
	}
	return out
}

func stringify(v interface{}) string {
	if v == nil {
		return ""
	}
	switch x := v.(type) {
	case string:
		return x
	case fmt.Stringer:
		return x.String()
	default:
		return fmt.Sprint(x)
	}
}

func parseRecordInputs(data map[string]interface{}) ([]domain.LineItemInput, error) {
	raw := issuesFromCase(data)
	if raw == nil {
		for _, key := range recordListKeys {
			if v, ok := data[key]; ok {
				raw = v
				break
			}
		}
	}
	if raw == nil {
		return nil, nil
	}

	switch list := raw.(type) {
	case []map[string]interface{}:
		return mapsToInputs(list)
	case []interface{}:
		maps := make([]map[string]interface{}, 0, len(list))
		for _, elem := range list {
			m, ok := elem.(map[string]interface{})
			if !ok {
				return nil, fmt.Errorf("records entries must be map[string]interface{}")
			}
			maps = append(maps, m)
		}
		return mapsToInputs(maps)
	default:
		return nil, fmt.Errorf("records must be a slice")
	}
}

func mapsToInputs(rows []map[string]interface{}) ([]domain.LineItemInput, error) {
	out := make([]domain.LineItemInput, 0, len(rows))
	for i, row := range rows {
		title := field(row, "Title", "title")
		if title == "" {
			title = fmt.Sprintf("CONCERN %d", i+1)
		}
		brief := field(row, "brief", "ConcernBrief", "concernBrief", "Brief")
		summary := field(row, "ConcernSummary", "concernSummary", "Summary")
		if summary == "" {
			summary = brief
		}
		out = append(out, domain.LineItemInput{
			Title:                title,
			ConcernBrief:         brief,
			ConcernSummary:       summary,
			Category:             field(row, "Category", "category"),
			FactsAndEvidence:     field(row, "FactsAndEvidence", "factsAndEvidence"),
			TeamsChatFindings:    field(row, "TeamsChatFindings", "teamsChatFindings"),
			ColleagueAStatements: field(row, "ColleagueAStatements", "colleagueAStatements"),
			ColleagueBStatements: field(row, "ColleagueBStatements", "colleagueBStatements"),
			FindingOutcome:       field(row, "FindingOutcome", "findingOutcome"),
		})
	}
	return out, nil
}

func issuesFromCase(data map[string]interface{}) interface{} {
	caseObj, ok := data["case"].(map[string]interface{})
	if !ok {
		return nil
	}
	return caseObj["issues"]
}

func field(row map[string]interface{}, keys ...string) string {
	for _, k := range keys {
		if v, ok := row[k]; ok {
			return stringify(v)
		}
	}
	return ""
}

func parseTime(s string) (time.Time, error) {
	if t, err := time.Parse(time.RFC3339, s); err == nil {
		return t, nil
	}
	if t, err := time.ParseInLocation("2006-01-02 15:04", s, time.Local); err == nil {
		return t, nil
	}
	return time.ParseInLocation("2006-01-02", s, time.Local)
}
