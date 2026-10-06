package docxtoken

import (
	"path/filepath"
	"strings"
	"testing"

	"github.com/liboyang/docx-token/internal/template"
)

func TestCustomLoopPrefixFill(t *testing.T) {
	dir := t.TempDir()
	tpl := filepath.Join(dir, "t.docx")
	xml := `<w:document xmlns:w="http://schemas.openxmlformats.org/wordprocessingml/2006/main"><w:body>
<w:p><w:r><w:t>{BEGIN:case.case_issues:record}</w:t></w:r></w:p>
<w:p><w:r><w:t>{record.code}: {record.title}</w:t></w:r></w:p>
<w:p><w:r><w:t>{END:case.case_issues:record}</w:t></w:r></w:p>
</w:body></w:document>`
	if err := writeMinimalDocx(tpl, xml); err != nil {
		t.Fatal(err)
	}
	data := map[string]interface{}{
		"case": map[string]interface{}{
			"case_issues": []map[string]interface{}{
				{"code": "C-1", "title": "Alpha"},
				{"code": "C-2", "title": "Beta"},
			},
		},
	}
	out := filepath.Join(dir, "out.docx")
	if err := Fill(tpl, out, data); err != nil {
		t.Fatal(err)
	}
	body, err := readDocBody(out)
	if err != nil {
		t.Fatal(err)
	}
	if strings.Contains(body, "{record") || strings.Contains(body, "BEGIN:") {
		t.Fatalf("unexpanded: %s", body)
	}
	if !strings.Contains(body, "C-1") || !strings.Contains(body, "C-2") {
		t.Fatalf("missing row data: %s", body)
	}
	specs, err := template.DiscoverLoopSpecsFromDocx(mustRead(tpl))
	if err != nil || len(specs) != 1 || specs[0].ItemPrefix != "record" {
		t.Fatalf("specs %+v err %v", specs, err)
	}
}
