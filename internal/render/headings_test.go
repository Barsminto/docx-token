package render

import (
	"strings"
	"testing"

	"github.com/liboyang/docx-token/internal/domain"
)

func TestPatchIssueHeadingsCombinedToken(t *testing.T) {
	xml := `<w:p><w:r><w:t>{issue_0.numberedTitle}</w:t></w:r></w:p>`
	items := []domain.LineItem{{Seq: "3", Title: "CONCERN 1", NumberedTitle: "3. INVESTIGATION: CONCERN 1"}}
	out := patchIssueHeadings(xml, items)
	if strings.Contains(out, "{issue_0") {
		t.Fatalf("tokens should be replaced: %s", out)
	}
	if !strings.Contains(out, "3. INVESTIGATION: CONCERN 1") {
		t.Fatal("expected numbered title")
	}
}
