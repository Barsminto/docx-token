package template

import (
	"fmt"
	"os"

	"github.com/gomutex/godocx"
	"github.com/liboyang/docx-token/internal/domain"
)

// WriteReportTemplate writes an English docx template matching the issue summary + detail layout.
func WriteReportTemplate(path string) error {
	doc, err := godocx.NewDocument()
	if err != nil {
		return err
	}

	if _, err := doc.AddHeading(domain.Wrap(domain.Key("Report", "Title")), 1); err != nil {
		return err
	}
	doc.AddParagraph(
		"Report date: " + domain.Wrap(domain.Key("Report", "Date")) +
			" | Prepared by: " + domain.Wrap(domain.Key("Summary", "Author")),
	)

	if _, err := doc.AddHeading("1. Summary", 1); err != nil {
		return err
	}
	doc.AddParagraph(domain.Wrap(domain.Key("Summary", "Overview")))
	if _, err := doc.AddHeading("Issue titles", 2); err != nil {
		return err
	}

	doc.AddParagraph(BeginIssueListMarker)
	doc.AddParagraph("• " + domain.Wrap(domain.LoopField("Title")))
	doc.AddParagraph(EndIssueListMarker)

	doc.AddParagraph(BeginRecordMarker)

	if _, err := doc.AddHeading(domain.Wrap(domain.LoopField("NumberedTitle")), 1); err != nil {
		return err
	}
	if _, err := doc.AddHeading("Category", 2); err != nil {
		return err
	}
	doc.AddParagraph(domain.Wrap(domain.LoopField("Category")))
	if _, err := doc.AddHeading("Impact", 2); err != nil {
		return err
	}
	doc.AddParagraph("Fixed narrative text for this issue. Replace in Word as needed; no token required.")
	if _, err := doc.AddHeading("Next steps", 3); err != nil {
		return err
	}
	doc.AddParagraph("Fixed checklist content preserved on each loop iteration.")

	doc.AddParagraph(EndRecordMarker)

	if err := doc.SaveTo(path); err != nil {
		return fmt.Errorf("save template: %w", err)
	}
	if err := doc.Close(); err != nil {
		return err
	}
	if err := attachDefaultHeader(path); err != nil {
		return fmt.Errorf("attach header: %w", err)
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
