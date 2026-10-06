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
// Loops in the template use {BEGIN:dot.path:prefix} … {END:dot.path:prefix}.
// Inside the loop use {prefix.field}; listIndex (1,2,3…) is added automatically.
// Legacy GI templates use case.issues + prefix issue ({issue.brief}).
//
// Example:
//
//	err := docxtoken.Generate(docxtoken.Options{
//	    TemplatePath: "/path/to/template.docx",
//	    OutputPath:   "/path/to/report.docx",
//	    Data: map[string]interface{}{
//	        "projectName": "PROJECT XXX",
//	        "case": map[string]interface{}{
//	            "case_issues": []map[string]interface{}{
//	                {"title": "CONCERN 1", "brief": "...", "code": "C-1"},
//	            },
//	        },
//	    },
//	})
//	// Template: {BEGIN:case.case_issues:record} … {record.code} … {END:case.case_issues:record}
