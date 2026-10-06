package template

import (
	"strings"
	"testing"
)

func TestInjectNumPrConcernLine(t *testing.T) {
	p := `<w:p><w:r><w:t>Concern {issue.listIndex}: {issue.brief}</w:t></w:r></w:p>`
	out := injectNumPr(p, 1, 0)
	if !strings.Contains(out, `w:val="1"`) {
		t.Fatalf("missing numPr: %s", out)
	}
}

func TestPatchSubsectionParagraphs(t *testing.T) {
	xml := `<w:p><w:r><w:t>2. INVESTIGATION METHODOLOGY</w:t></w:r></w:p>` +
		`<w:p><w:r><w:t>` + DotMarker + ` GI applied...</w:t></w:r></w:p>`
	out := patchSubsectionParagraphs(xml, false)
	if !strings.Contains(out, `numId w:val="5"`) || strings.Contains(out, DotMarker) {
		t.Fatalf("subsection patch failed: %s", out)
	}
}

func TestPatchConcernListParagraphs(t *testing.T) {
	xml := `<w:p><w:r><w:t>` + BeginIssueListMarker + `</w:t></w:r></w:p>` +
		`<w:p><w:r><w:t>Concern {issue.listIndex}: {issue.brief}</w:t></w:r></w:p>` +
		`<w:p><w:r><w:t>` + EndIssueListMarker + `</w:t></w:r></w:p>`
	out := patchConcernListParagraphs(xml, BeginIssueListMarker, EndIssueListMarker, 1)
	if !strings.Contains(out, "numId") {
		t.Fatal("expected numbering on concern line")
	}
}
