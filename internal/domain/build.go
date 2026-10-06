package domain

import "time"

// LineItemInput source row for one concern (maps to template {Record.*} fields).
type LineItemInput struct {
	Title                string
	ConcernBrief         string
	ConcernSummary       string
	Category             string
	FactsAndEvidence     string
	TeamsChatFindings    string
	ColleagueAStatements string
	ColleagueBStatements string
	FindingOutcome       string
}

func (i LineItemInput) ToLineItem() LineItem {
	return LineItem{
		Title:                i.Title,
		ConcernBrief:         i.ConcernBrief,
		ConcernSummary:       i.ConcernSummary,
		Category:             i.Category,
		FactsAndEvidence:     i.FactsAndEvidence,
		TeamsChatFindings:    i.TeamsChatFindings,
		ColleagueAStatements: i.ColleagueAStatements,
		ColleagueBStatements: i.ColleagueBStatements,
		FindingOutcome:       i.FindingOutcome,
	}
}

// NewDocument builds the root object. generatedAt is supplied by the caller (typically time.Now()).
func NewDocument(title, author string, generatedAt time.Time, items []LineItemInput) Document {
	if generatedAt.IsZero() {
		generatedAt = time.Now()
	}
	lines := make([]LineItem, 0, len(items))
	for _, item := range items {
		lines = append(lines, item.ToLineItem())
	}
	return Document{
		Report: Report{
			Title:       title,
			Date:        generatedAt.Format("2006-01-02"),
			GeneratedAt: generatedAt.Format("2006-01-02 15:04"),
		},
		Summary: Summary{
			Author:   author,
			Overview: "",
		},
		Items: lines,
	}
}
