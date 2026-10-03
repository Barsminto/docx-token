package render

import (
	"strings"
	"testing"

	"github.com/liboyang/docx-token/internal/domain"
)

func TestPatchIssueHeadingsCombinedToken(t *testing.T) {
	xml := `<w:p><w:r><w:t>{Record_0.Seq}. {Record_0.Title}</w:t></w:r></w:p>`
	items := []domain.LineItem{{Seq: "2", Title: "INC-1", NumberedTitle: "2. INC-1"}}
	out := patchIssueHeadings(xml, items)
	if strings.Contains(out, "{Record_0") {
		t.Fatalf("tokens should be replaced: %s", out)
	}
	if !strings.Contains(out, "2. INC-1") {
		t.Fatal("expected numbered title")
	}
}
