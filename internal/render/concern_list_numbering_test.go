package render

import (
	"strings"
	"testing"

	"github.com/liboyang/docx-token/internal/domain"
	"github.com/liboyang/docx-token/internal/template"
)

func TestConcernListNotTreatedAsInvestigationChapter(t *testing.T) {
	xml := `<w:body>
<w:p><w:r><w:t>` + template.DotMarker + ` GI reviewed {subjectList}:</w:t></w:r></w:p>
<w:p><w:r><w:t>Concern {issue_0.listIndex}: {issue_0.brief}</w:t></w:r></w:p>
<w:p><w:r><w:t>` + template.DotMarker + ` The review identified that ...</w:t></w:r></w:p>
<w:p><w:r><w:rPr><w:b/></w:rPr><w:t>INVESTIGATION: CONCERN 1</w:t></w:r></w:p>
</w:body>`
	out := patchRecordSubsections(xml, []domain.LineItem{{Title: "CONCERN 1", SectionNumber: "3"}})
	out = reapplyConcernListNumIDs(out)
	var concernPara string
	for _, p := range paragraphRE.FindAllString(out, -1) {
		if strings.Contains(p, "Concern {issue_0") {
			concernPara = p
			break
		}
	}
	if concernPara == "" {
		t.Fatal("concern paragraph not found")
	}
	concernChunk := concernPara
	if strings.Contains(concernChunk, `numId w:val="100"`) {
		t.Fatalf("concern list should not use investigation numId 100: %s", concernChunk)
	}
	if !strings.Contains(concernChunk, `numId w:val="1"`) {
		t.Fatalf("executive concern list should use letter list numId 1: %s", concernChunk)
	}
}
