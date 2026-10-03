package template

// Loop region markers — copy whole paragraphs between begin/end; only {Record.*} keys are renamed.
const (
	BeginIssueListMarker = "{BEGIN_ISSUE_LIST}"
	EndIssueListMarker   = "{END_ISSUE_LIST}"
	BeginRecordMarker    = "{BEGIN_RECORD}"
	EndRecordMarker      = "{END_RECORD}"
)

// LoopRegion defines one repeatable block in the template.
type LoopRegion struct {
	Begin string
	End   string
}

// LoopRegions order matters: expand the summary list first, then issue detail sections.
var LoopRegions = []LoopRegion{
	{Begin: BeginIssueListMarker, End: EndIssueListMarker},
	{Begin: BeginRecordMarker, End: EndRecordMarker},
}
