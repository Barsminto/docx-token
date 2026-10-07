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

	specs := []template.LoopSpec{
		{Begin: template.BeginIssueListMarker, End: template.EndIssueListMarker, ItemPrefix: "issue"},
		{Begin: template.BeginRecordMarker, End: template.EndRecordMarker, ItemPrefix: "issue"},
	}
	counts := map[string]int{
		template.BeginIssueListMarker: 2,
		template.BeginRecordMarker:      2,
	}
	out, err := expandAllLoops(xml, specs, counts)
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

func TestExpandLoopFragmentedListIndexInBody(t *testing.T) {
	begin := "{BEGIN:case.items:cn}"
	end := "{END:case.items:cn}"
	// Word often splits placeholders in body paragraphs across runs; title may stay in one run.
	xml := `<w:body>
<w:p><w:r><w:t>` + begin + `</w:t></w:r></w:p>
<w:p><w:r><w:t>{cn.listIndex}. Title</w:t></w:r></w:p>
<w:p><w:r><w:t>{</w:t></w:r><w:r><w:t>cn.listIndex</w:t></w:r><w:r><w:t>}</w:t></w:r></w:p>
<w:p><w:r><w:t>` + end + `</w:t></w:r></w:p>
</w:body>`
	specs := []template.LoopSpec{
		{Begin: begin, End: end, DataPath: "case.items", ItemPrefix: "cn"},
	}
	counts := map[string]int{begin: 2}
	docXML := repairFragmentedPlaceholders(xml)
	out, err := expandAllLoops(docXML, specs, counts)
	if err != nil {
		t.Fatal(err)
	}
	if strings.Contains(out, "{cn.listIndex}") || strings.Contains(out, "{cn.") {
		t.Fatalf("loop prefix should be expanded: %s", out)
	}
	if !strings.Contains(out, "{cn_0.listIndex}") || !strings.Contains(out, "{cn_1.listIndex}") {
		t.Fatalf("expected indexed placeholders: %s", out)
	}
}
