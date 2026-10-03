package domain

import "time"

// LineItemInput source row for one issue.
type LineItemInput struct {
	ID       int64
	Title    string
	Category string
}

// NewDocument builds the root object. generatedAt is supplied by the caller (typically time.Now()).
func NewDocument(title, author string, generatedAt time.Time, items []LineItemInput) Document {
	if generatedAt.IsZero() {
		generatedAt = time.Now()
	}
	lines := make([]LineItem, 0, len(items))
	for _, item := range items {
		lines = append(lines, LineItem{
			Title:    item.Title,
			Category: item.Category,
		})
	}
	return Document{
		Report: Report{
			Title:       title,
			Date:        generatedAt.Format("2006-01-02"),
			GeneratedAt: generatedAt.Format("2006-01-02 15:04"),
		},
		Summary: Summary{
			Author:   author,
			Overview: "Executive summary text. Issue titles are listed below; full detail starts at section 2.",
		},
		Items: lines,
	}
}
