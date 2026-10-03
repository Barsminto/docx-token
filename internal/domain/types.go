package domain

// Report report-level metadata and header tokens.
type Report struct {
	Title       string
	Date        string
	GeneratedAt string
}

// Summary section under "1. Summary".
type Summary struct {
	Author   string
	Overview string
}

// LineItem one issue: Title for the summary list; NumberedTitle for section heading (Seq + Title).
type LineItem struct {
	Seq           string
	Title         string
	NumberedTitle string
	Category      string
}

// Document root object; Placeholders() flattens to {Report.*}, {Summary.*}, {Record_0.Title}, …
type Document struct {
	Report  Report
	Summary Summary
	Items   []LineItem
}
