package docxtoken

import (
	"fmt"
	"strings"
	"time"

	"github.com/liboyang/docx-token/internal/domain"
)

// RecordsKey is the preferred map key for issue rows used in template loops.
const RecordsKey = "records"

var recordListKeys = []string{RecordsKey, "items", "Records", "Items"}

// BuildFromMap builds a Document from a flat or nested map.
//
// Top-level scalar keys use placeholder names, e.g. Report.Title, Summary.Author.
// Nested objects are flattened with dots: Report -> Title becomes Report.Title.
// Issue rows: data[records] as []map[string]interface{} (or []interface{} of maps).
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

	title := flat["Report.Title"]
	if title == "" {
		title = "Report"
	}
	author := flat["Summary.Author"]
	if author == "" {
		author = "System Generated"
	}

	doc := domain.NewDocument(title, author, generatedAt, items)
	if s := flat["Summary.Overview"]; s != "" {
		doc.Summary.Overview = s
	}
	if s := flat["Report.Date"]; s != "" {
		doc.Report.Date = s
	}
	if s := flat["Report.GeneratedAt"]; s != "" {
		doc.Report.GeneratedAt = s
	}
	return doc, nil
}

// ScalarPlaceholders returns non-record keys as placeholder map (merged after record expansion).
func ScalarPlaceholders(data map[string]interface{}) map[string]string {
	out := flattenMap("", data)
	for _, key := range recordListKeys {
		delete(out, key)
	}
	for k := range out {
		if strings.HasPrefix(k, "Record.") || strings.HasPrefix(k, "Record_") {
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
	var raw interface{}
	for _, key := range recordListKeys {
		if v, ok := data[key]; ok {
			raw = v
			break
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
		title := stringify(row["Title"])
		if title == "" {
			title = stringify(row["title"])
		}
		category := stringify(row["Category"])
		if category == "" {
			category = stringify(row["category"])
		}
		if title == "" {
			return nil, fmt.Errorf("records[%d]: Title is required", i)
		}
		out = append(out, domain.LineItemInput{
			ID:       int64(i + 1),
			Title:    title,
			Category: category,
		})
	}
	return out, nil
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
