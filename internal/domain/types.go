package domain

// Report header / metadata tokens.
type Report struct {
	Title       string
	Date        string
	GeneratedAt string
}

// Summary optional summary block (legacy keys).
type Summary struct {
	Author   string
	Overview string
}

// LineItem one concern row: list + full investigation section in {BEGIN_RECORD}.
type LineItem struct {
	Seq                  string
	ListIndex            string
	ListLetter           string // a, b, c … for executive summary concern lines
	StmtLabelA           string // interview sub-bullet (a, c, …)
	StmtLabelB           string // interview sub-bullet (b, d, …)
	Title                string
	NumberedTitle        string
	SectionNumber        string
	ConcernBrief         string
	ConcernSummary       string
	Category             string
	FactsAndEvidence     string
	TeamsChatFindings    string
	ColleagueAStatements string
	ColleagueBStatements string
	FindingOutcome       string
}

// Document root object passed to the renderer.
type Document struct {
	Report  Report
	Summary Summary
	Items   []LineItem
}
