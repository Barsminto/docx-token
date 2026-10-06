package template

import (
	"fmt"
	"os"

	"github.com/gomutex/godocx"
	"github.com/gomutex/godocx/docx"
	"github.com/liboyang/docx-token/internal/domain"
)

func addBoldChapterTitle(doc *docx.RootDoc, text string) {
	doc.AddParagraph(GiBoldMarker + text)
}

func issueListLine() string {
	return "Concern " + domain.Wrap(domain.LoopField("listIndex")) + ": " + domain.Wrap(domain.LoopField("brief"))
}

// WriteReportTemplate GI layout using {projectName}, {case.name}, {issue.*}, {BEGIN:case.issues}.
func WriteReportTemplate(path string) error {
	doc, err := godocx.NewDocument()
	if err != nil {
		return err
	}
	appendGIReportBody(doc)

	if err := doc.SaveTo(path); err != nil {
		return fmt.Errorf("save template: %w", err)
	}
	if err := doc.Close(); err != nil {
		return err
	}
	if err := attachDefaultHeader(path); err != nil {
		return fmt.Errorf("attach header: %w", err)
	}
	if err := attachDocumentNumbering(path, NumberingAttachOptions{}); err != nil {
		return fmt.Errorf("attach document numbering: %w", err)
	}
	return nil
}

func EnsureReportTemplate(path string) (created bool, err error) {
	if fileExists(path) {
		return false, nil
	}
	if err := WriteReportTemplate(path); err != nil {
		return false, err
	}
	return true, nil
}

func fileExists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}
