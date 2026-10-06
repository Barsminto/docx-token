package render

import (
	"strings"
	"testing"
)

func TestTrimEmptyParagraphsBeforeInvestigationReport(t *testing.T) {
	in := `<w:p><w:r><w:t>Summary end</w:t></w:r></w:p>` +
		`<w:p><w:r><w:t></w:t></w:r></w:p>` +
		`<w:p><w:r><w:t></w:t></w:r></w:p>` +
		`<w:p><w:r><w:t>INVESTIGATION REPORT</w:t></w:r></w:p>`
	out := trimEmptyParagraphsBeforeHeading(in, "INVESTIGATION REPORT")
	if strings.Count(out, "<w:p") != 2 {
		t.Fatalf("expected 2 paragraphs, got: %s", out)
	}
}

func TestEnsurePageBreakBeforeInvestigationReport(t *testing.T) {
	in := `<w:p><w:r><w:t>INVESTIGATION REPORT</w:t></w:r></w:p>`
	out := ensurePageBreakBeforeHeading(in, "INVESTIGATION REPORT")
	if !strings.Contains(out, "w:pageBreakBefore") {
		t.Fatalf("expected page break: %s", out)
	}
}

func TestReplaceParagraphTextPreservesIndent(t *testing.T) {
	in := `<w:p><w:pPr><w:ind w:left="720" w:hanging="360"/><w:numPr><w:ilvl w:val="1"/><w:numId w:val="3"/></w:numPr></w:pPr>` +
		`<w:r><w:rPr><w:rFonts w:ascii="Arial"/></w:rPr><w:t>{.} hello</w:t></w:r></w:p>`
	out := stripMarkersInParagraph(in, "{.}")
	out = trimParagraphTextRuns(out)
	if !strings.Contains(out, `w:left="720"`) || !strings.Contains(out, `w:hanging="360"`) {
		t.Fatalf("indent lost: %s", out)
	}
	if strings.Contains(out, "{.}") {
		t.Fatal("marker should be removed")
	}
}

func TestInheritParagraphRunFontsFromPPr(t *testing.T) {
	in := `<w:p><w:pPr><w:rPr><w:rFonts w:ascii="Arial" w:hAnsi="Arial"/><w:sz w:val="18"/></w:rPr></w:pPr><w:r><w:t>hello</w:t></w:r></w:p>`
	out := inheritParagraphRunFonts(in)
	if !strings.Contains(out, `<w:r><w:rPr><w:rFonts w:ascii="Arial"`) {
		t.Fatalf("run should inherit Arial: %s", out)
	}
}
