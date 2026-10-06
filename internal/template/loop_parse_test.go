package template

import "testing"

func TestParseLoopMarkerBody(t *testing.T) {
	tests := []struct {
		body           string
		wantPath       string
		wantPrefix     string
	}{
		{"case.issues", "case.issues", "issue"},
		{"case.issues:background", "case.issues", "issue"},
		{"case.issues:record", "case.issues", "issue"},
		{"case.case_issues:record", "case.case_issues", "record"},
		{"data.items:row", "data.items", "row"},
	}
	for _, tc := range tests {
		path, prefix := ParseLoopMarkerBody(tc.body)
		if path != tc.wantPath || prefix != tc.wantPrefix {
			t.Fatalf("%q → path=%q prefix=%q, want %q %q", tc.body, path, prefix, tc.wantPath, tc.wantPrefix)
		}
	}
}

func TestDiscoverLoopSpecsCustomPrefix(t *testing.T) {
	xml := `<w:p>{BEGIN:case.case_issues:record}</w:p><w:p>{record.title}</w:p><w:p>{END:case.case_issues:record}</w:p>`
	specs, err := DiscoverLoopSpecs(xml)
	if err != nil {
		t.Fatal(err)
	}
	if len(specs) != 1 || specs[0].DataPath != "case.case_issues" || specs[0].ItemPrefix != "record" {
		t.Fatalf("specs: %+v", specs)
	}
}
