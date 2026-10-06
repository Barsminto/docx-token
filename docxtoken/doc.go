package docxtoken

// Package docxtoken fills Word (.docx) templates from map[string]interface{} data.
//
// Public API:
//   - Generate(Options) — preferred; explicit template/output paths
//   - Fill(templatePath, outputPath, data) — shorthand
//
// Only import github.com/liboyang/docx-token/docxtoken from other projects.
// internal/* packages are not a stable API.
//
// Example:
//
//	err := docxtoken.Generate(docxtoken.Options{
//	    TemplatePath: "/path/to/templates/template.docx",
//	    OutputPath:   "/path/to/output/report.docx",
//	    Data: map[string]interface{}{
//	        "projectName": "PROJECT XXX",
//	        "draftDate":   "01 February 2026",
//	        "case": map[string]interface{}{
//	            "name": "Case title",
//	            "issues": []map[string]interface{}{
//	                {"title": "CONCERN 1", "brief": "...", "summary": "..."},
//	            },
//	        },
//	    },
//	})
