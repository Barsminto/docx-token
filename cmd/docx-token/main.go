package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/liboyang/docx-token/docxtoken"
	"github.com/liboyang/docx-token/internal/project"
	"github.com/liboyang/docx-token/internal/sample"
	"github.com/liboyang/docx-token/internal/template"
)

func main() {
	defaultTpl, err := project.TemplatePath()
	exitOnErr(err)
	defaultOut, err := project.OutputPath()
	exitOnErr(err)

	templatePath := flag.String("template", defaultTpl, "docx template path")
	outputPath := flag.String("out", defaultOut, "output docx path")
	skipTemplate := flag.Bool("skip-template", false, "do not auto-create template if missing")
	writeTemplate := flag.Bool("write-template", false, "overwrite template with built-in layout")
	flag.Parse()

	exitOnErr(os.MkdirAll(filepath.Dir(*templatePath), 0o755))
	if *writeTemplate {
		exitOnErr(template.WriteReportTemplate(*templatePath))
		fmt.Printf("Template overwritten: %s\n", *templatePath)
	} else if !*skipTemplate {
		created, err := template.EnsureReportTemplate(*templatePath)
		exitOnErr(err)
		if created {
			fmt.Printf("Template created: %s\n", *templatePath)
		}
	}

	data := sampleDataMap(time.Now())
	exitOnErr(os.MkdirAll(filepath.Dir(*outputPath), 0o755))
	exitOnErr(docxtoken.Fill(*templatePath, *outputPath, data))

	fmt.Printf("Generated: %s\n", *outputPath)
}

func sampleDataMap(now time.Time) map[string]interface{} {
	records := make([]map[string]interface{}, 0, len(sample.Items()))
	for _, item := range sample.Items() {
		records = append(records, map[string]interface{}{
			"Title":    item.Title,
			"Category": item.Category,
		})
	}
	return map[string]interface{}{
		"Report.Title":       "Business Data Report",
		"Summary.Author":     "System Generated",
		"Report.GeneratedAt": now.Format("2006-01-02 15:04"),
		docxtoken.RecordsKey: records,
	}
}

func exitOnErr(err error) {
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}
}
