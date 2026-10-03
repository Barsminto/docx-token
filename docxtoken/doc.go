package docxtoken

// Package docxtoken fills Word templates from map[string]interface{} data.
//
// Example:
//
//	err := docxtoken.Fill("templates/template.docx", "output/report.docx", map[string]interface{}{
//	    "Report.Title": "FY26 Q1 Report",
//	    "Summary.Author": "Jane Doe",
//	    "records": []map[string]interface{}{
//	        {"Title": "INC-1001", "Category": "Hardware"},
//	    },
//	})
