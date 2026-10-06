package render

import (
	"strings"
	"testing"

	"github.com/liboyang/docx-token/internal/template"
)

func TestExpandIssueListAndDetailLoops(t *testing.T) {
	xml := `<w:body>
<w:p><w:r><w:t>` + template.BeginIssueListMarker + `</w:t></w:r></w:p>
<w:p><w:r><w:t>{issue.title}</w:t></w:r></w:p>
<w:p><w:r><w:t>` + template.EndIssueListMarker + `</w:t></w:r></w:p>
<w:p><w:r><w:t>` + template.BeginRecordMarker + `</w:t></w:r></w:p>
<w:p><w:r><w:t>{issue.numberedTitle}</w:t></w:r></w:p>
<w:p><w:r><w:t>{issue.category}</w:t></w:r></w:p>
<w:p><w:r><w:t>Static body</w:t></w:r></w:p>
<w:p><w:r><w:t>` + template.EndRecordMarker + `</w:t></w:r></w:p>
</w:body>`

	out, err := expandAllLoops(xml, 2)
	if err != nil {
		t.Fatal(err)
	}
	if strings.Contains(out, template.BeginRecordMarker) {
		t.Fatal("markers should be removed")
	}
	if !strings.Contains(out, "{issue_0.category}") || !strings.Contains(out, "{issue_1.category}") {
		t.Fatalf("unexpected: %s", out)
	}
	if strings.Count(out, "Static body") != 2 {
		t.Fatal("static paragraphs should be duplicated per concern")
	}
}
