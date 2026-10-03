package render

import (
	"strings"
	"testing"

	"github.com/liboyang/docx-token/internal/template"
)

func TestExpandIssueListAndDetailLoops(t *testing.T) {
	xml := `<w:body>
<w:p><w:r><w:t>` + template.BeginIssueListMarker + `</w:t></w:r></w:p>
<w:p><w:r><w:t>• {Record.Title}</w:t></w:r></w:p>
<w:p><w:r><w:t>` + template.EndIssueListMarker + `</w:t></w:r></w:p>
<w:p><w:r><w:t>` + template.BeginRecordMarker + `</w:t></w:r></w:p>
<w:p><w:r><w:t>{Record.Seq}. {Record.Title}</w:t></w:r></w:p>
<w:p><w:r><w:t>Category: {Record.Category}</w:t></w:r></w:p>
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
	if !strings.Contains(out, "{Record_0.Title}") || !strings.Contains(out, "{Record_1.Category}") {
		t.Fatalf("unexpected: %s", out)
	}
	if strings.Count(out, "Static body") != 2 {
		t.Fatal("static paragraphs should be duplicated per issue")
	}
}
