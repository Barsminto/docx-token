package domain

import (
	"fmt"
	"reflect"
	"strings"
)

// LoopItemPrefix is the loop variable in {BEGIN:…} blocks before expansion ({issue.title} → {issue_0.title}).
const LoopItemPrefix = "issue"

// Key builds a placeholder key such as Report.Title.
func Key(parts ...string) string {
	return strings.Join(parts, ".")
}

// Wrap formats a key as a docx token literal, e.g. {Report.Title}.
func Wrap(key string) string {
	return "{" + key + "}"
}

// LoopField builds a prototype loop field key, e.g. issue.brief.
func LoopField(field string) string {
	return Key(LoopItemPrefix, field)
}

// Placeholders returns a flat map for go-docx after loop expansion (issue_0.brief, …).
func (d Document) Placeholders() map[string]string {
	out := make(map[string]string)
	bindObject(out, "Report", d.Report)
	bindObject(out, "Summary", d.Summary)
	for i, item := range d.Items {
		bindIssue(out, fmt.Sprintf("%s_%d", LoopItemPrefix, i), item)
	}
	return out
}

func bindIssue(dst map[string]string, prefix string, item LineItem) {
	dst[Key(prefix, "listIndex")] = item.ListIndex
	dst[Key(prefix, "listLetter")] = item.ListLetter
	dst[Key(prefix, "title")] = item.Title
	dst[Key(prefix, "brief")] = item.ConcernBrief
	dst[Key(prefix, "summary")] = item.ConcernSummary
	dst[Key(prefix, "sectionNumber")] = item.SectionNumber
	dst[Key(prefix, "numberedTitle")] = item.NumberedTitle
	dst[Key(prefix, "stmtLabelA")] = item.StmtLabelA
	dst[Key(prefix, "stmtLabelB")] = item.StmtLabelB
	dst[Key(prefix, "category")] = item.Category
	dst[Key(prefix, "factsAndEvidence")] = item.FactsAndEvidence
	dst[Key(prefix, "teamsChatFindings")] = item.TeamsChatFindings
	dst[Key(prefix, "colleagueAStatements")] = item.ColleagueAStatements
	dst[Key(prefix, "colleagueBStatements")] = item.ColleagueBStatements
	dst[Key(prefix, "findingOutcome")] = item.FindingOutcome
}

func bindObject(dst map[string]string, prefix string, value any) {
	rv := reflect.ValueOf(value)
	rt := rv.Type()
	if rt.Kind() != reflect.Struct {
		return
	}
	for i := 0; i < rt.NumField(); i++ {
		field := rt.Field(i)
		if field.PkgPath != "" {
			continue
		}
		dst[Key(prefix, field.Name)] = fmt.Sprint(rv.Field(i).Interface())
	}
}
