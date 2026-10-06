package template

// Structural markers (removed during template build or fill). Data placeholders use {path.to.field}.
const (
	DotMarker   = "{.}" // decimal subsection (1.1, 8.1, 3.1) — Word list ilvl 1
	TildeMarker = "{~}" // letter sub-item (a., b.)

	BeginCaseIssuesMarker           = "{BEGIN:case.issues}"
	EndCaseIssuesMarker             = "{END:case.issues}"
	BeginCaseIssuesBackgroundMarker = "{BEGIN:case.issues:background}"
	EndCaseIssuesBackgroundMarker   = "{END:case.issues:background}"
	BeginCaseIssuesRecordMarker     = "{BEGIN:case.issues:record}"
	EndCaseIssuesRecordMarker       = "{END:case.issues:record}"

	// Legacy aliases (same strings as above where applicable).
	BeginIssueListMarker      = BeginCaseIssuesMarker
	EndIssueListMarker        = EndCaseIssuesMarker
	BeginBackgroundListMarker = BeginCaseIssuesBackgroundMarker
	EndBackgroundListMarker   = EndCaseIssuesBackgroundMarker
	BeginRecordMarker         = BeginCaseIssuesRecordMarker
	EndRecordMarker           = EndCaseIssuesRecordMarker

	SubsecCh1Marker        = DotMarker
	SubsecCh1RestartMarker = DotMarker
	SubsecCh2Marker        = DotMarker
	SubsecRecordMarker     = DotMarker
	SubsecConclusionMarker = DotMarker
	SubsecRecommendMarker  = DotMarker
	SubsecLetterMarker     = TildeMarker
)

// LoopRegion duplicates paragraphs between Begin and End; ItemPrefix is the loop variable (issue → issue_0).
type LoopRegion struct {
	Begin      string
	End        string
	ItemPrefix string
}

// LoopRegions order matters.
var LoopRegions = []LoopRegion{
	{Begin: BeginCaseIssuesMarker, End: EndCaseIssuesMarker, ItemPrefix: "issue"},
	{Begin: BeginCaseIssuesBackgroundMarker, End: EndCaseIssuesBackgroundMarker, ItemPrefix: "issue"},
	{Begin: BeginCaseIssuesRecordMarker, End: EndCaseIssuesRecordMarker, ItemPrefix: "issue"},
}
