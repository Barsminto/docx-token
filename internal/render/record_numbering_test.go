package render

import (
	"strings"
	"testing"

	"github.com/liboyang/docx-token/internal/domain"
	"github.com/liboyang/docx-token/internal/template"
)

func TestPatchRecordSubsections(t *testing.T) {
	xml := `<w:body>
<w:p><w:r><w:rPr><w:b/></w:rPr><w:t>INVESTIGATION: {issue_0.title}</w:t></w:r></w:p>
<w:p><w:r><w:t>` + template.SubsecRecordMarker + ` facts</w:t></w:r></w:p>
<w:p><w:r><w:t>Review of Communication Record</w:t></w:r></w:p>
<w:p><w:r><w:t>` + template.SubsecRecordMarker + ` Teams</w:t></w:r></w:p>
<w:p><w:r><w:rPr><w:b/></w:rPr><w:t>INVESTIGATION: {issue_1.title}</w:t></w:r></w:p>
<w:p><w:r><w:t>` + template.SubsecRecordMarker + ` facts 2</w:t></w:r></w:p>
</w:body>`
	items := []domain.LineItem{
		{SectionNumber: "3"},
		{SectionNumber: "4"},
	}
	out := patchRecordSubsections(xml, items)
	if strings.Contains(out, template.SubsecRecordMarker) {
		t.Fatal("marker should be stripped")
	}
	if !strings.Contains(out, `numId w:val="100"`) || !strings.Contains(out, `numId w:val="101"`) {
		t.Fatalf("expected per-record num ids: %s", out)
	}
}

func TestPatchRecordSubsectionsWithoutDotMarker(t *testing.T) {
	xml := `<w:body>
<w:p><w:r><w:rPr><w:b/></w:rPr><w:t>INVESTIGATION: {issue_0.title}</w:t></w:r></w:p>
<w:p><w:pPr><w:numPr><w:ilvl w:val="1"/><w:numId w:val="5"/></w:numPr></w:pPr><w:r><w:t>[facts or evidence identified].</w:t></w:r></w:p>
</w:body>`
	out := patchRecordSubsections(xml, []domain.LineItem{{SectionNumber: "3"}})
	if !strings.Contains(out, `numId w:val="100"`) || !strings.Contains(out, `ilvl w:val="1"`) {
		t.Fatalf("record body should share investigation numId 100 ilvl 1: %s", out)
	}
	if strings.Contains(out, `numId w:val="5"`) {
		t.Fatalf("stale methodology numId should be replaced: %s", out)
	}
}

func TestMergeRecordNumberingUsesDecimalAbstractFromTemplate(t *testing.T) {
	// Word-style numbering: abstract 0 = decimal, abstract 1 = letters (numId 6 → abstract 0).
	numbering := `<w:numbering xmlns:w="http://schemas.openxmlformats.org/wordprocessingml/2006/main">
<w:num w:numId="6"><w:abstractNumId w:val="0"/></w:num>
</w:numbering>`
	out := mergeRecordNumberingDefs(numbering, []domain.LineItem{{SectionNumber: "3"}})
	if !strings.Contains(out, `<w:num w:numId="100"><w:abstractNumId w:val="0"/>`) {
		t.Fatalf("record num should use decimal abstract 0: %s", out)
	}
}
