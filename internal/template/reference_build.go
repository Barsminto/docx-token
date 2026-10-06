package template

import (
	"fmt"

	"github.com/gomutex/godocx"
	"github.com/gomutex/godocx/docx"
	"github.com/liboyang/docx-token/internal/domain"
)

// WriteReferenceTemplate writes a full GI layout plus a visible marker guide (for copying when authoring new templates).
// Structural markers {.} {~} {GI_BOLD} are left in the document XML; numbering is still attached.
func WriteReferenceTemplate(path string) error {
	doc, err := godocx.NewDocument()
	if err != nil {
		return err
	}
	appendReferenceGuide(doc)
	appendGIReportBody(doc)

	if err := doc.SaveTo(path); err != nil {
		return fmt.Errorf("save template: %w", err)
	}
	if err := doc.Close(); err != nil {
		return err
	}
	if err := attachReferenceHeader(path); err != nil {
		return fmt.Errorf("attach header: %w", err)
	}
	if err := attachDocumentNumbering(path, NumberingAttachOptions{PreserveMarkers: true}); err != nil {
		return fmt.Errorf("attach document numbering: %w", err)
	}
	return nil
}

func appendReferenceGuide(doc *docx.RootDoc) {
	doc.AddParagraph("—— TEMPLATE REFERENCE (copy structure below; delete this guide in production) ——")
	doc.AddParagraph("Loops (repeat per case.issues row):")
	doc.AddParagraph(BeginCaseIssuesMarker + " … " + EndCaseIssuesMarker + "  (executive concern list)")
	doc.AddParagraph(BeginCaseIssuesBackgroundMarker + " … " + EndCaseIssuesBackgroundMarker)
	doc.AddParagraph(BeginCaseIssuesRecordMarker + " … " + EndCaseIssuesRecordMarker + "  (full INVESTIGATION block)")
	doc.AddParagraph("Subsection markers (start of paragraph):")
	doc.AddParagraph(DotMarker + "  decimal subsection → 1.1, 3.1, 8.1 (Word list ilvl 1)")
	doc.AddParagraph(TildeMarker + "  letter item → a., b. (under CONCLUSION)")
	doc.AddParagraph(GiBoldMarker + "  optional chapter title marker → bold (built-in generator only)")
	doc.AddParagraph("Scalars: {projectName} {draftDate} {accountableExecutive} {subjectList} {background.period}")
	doc.AddParagraph("Nested: {case.name} {executiveSummary.identified} {conclusion.support} {recommendation.1} …")
	doc.AddParagraph("Loop fields (inside regions above): {issue.listIndex} {issue.listLetter} {issue.brief} {issue.title} {issue.summary}")
	doc.AddParagraph("{issue.sectionNumber} {issue.numberedTitle} {issue.stmtLabelA} {issue.stmtLabelB}")
	doc.AddParagraph("Filled at generate time: {Conclusion.Number} {Recommendations.Number} {Report.Title} {Report.Date} {Report.GeneratedAt}")
	doc.AddParagraph("Legacy aliases: {BEGIN_ISSUE_LIST} = " + BeginCaseIssuesMarker + " (prefer new names)")
	doc.AddParagraph("—— END REFERENCE ——")
	doc.AddParagraph("")
}

func appendGIReportBody(doc *docx.RootDoc) {
	doc.AddParagraph("GROUP INVESTIGATION")
	doc.AddParagraph("DRAFT: " + domain.Wrap("draftDate") + " – the date of the report")
	doc.AddParagraph("PROJECT " + domain.Wrap("projectName"))
	doc.AddParagraph(domain.Wrap("case.name"))
	addBoldChapterTitle(doc, "EXECUTIVE SUMMARY")

	doc.AddParagraph("Accountable Executive: " + domain.Wrap("accountableExecutive"))
	doc.AddParagraph(DotMarker + " Group Investigations ('GI') reviewed concerns involving " + domain.Wrap("subjectList") + ":")
	doc.AddParagraph(BeginCaseIssuesMarker)
	doc.AddParagraph(issueListLine())
	doc.AddParagraph(EndCaseIssuesMarker)
	doc.AddParagraph(DotMarker + " The review identified that " + domain.Wrap("executiveSummary.identified"))
	doc.AddParagraph(
		DotMarker + " The concern that " + domain.Wrap("executiveSummary.concernRef") + " was " +
			domain.Wrap("executiveSummary.substantiation") + ". The evidence " + domain.Wrap("executiveSummary.evidence") + " ...",
	)
	doc.AddParagraph(DotMarker + " Recommendations " + domain.Wrap("executiveSummary.recommendations"))

	addBoldChapterTitle(doc, "INVESTIGATION REPORT")
	addBoldChapterTitle(doc, "1. BACKGROUND")
	doc.AddParagraph(
		DotMarker + " In " + domain.Wrap("background.period") +
			", concerns were escalated to GI that " + domain.Wrap("subjectList") + ":",
	)
	doc.AddParagraph(BeginCaseIssuesBackgroundMarker)
	doc.AddParagraph(issueListLine())
	doc.AddParagraph(EndCaseIssuesBackgroundMarker)
	doc.AddParagraph("(collectively, the 'Concerns')")

	addBoldChapterTitle(doc, "2. INVESTIGATION METHODOLOGY")
	doc.AddParagraph(DotMarker + " GI applied the following data review parameters to investigate the Concerns:")
	doc.AddParagraph("a. Custodians: ...")
	doc.AddParagraph("b. Data Reviewed: ...")
	doc.AddParagraph("c. Review Period: ...")
	doc.AddParagraph(DotMarker + " As a limitation, GI did not ...")

	doc.AddParagraph(BeginCaseIssuesRecordMarker)

	addBoldChapterTitle(doc, "INVESTIGATION: "+domain.Wrap(domain.LoopField("title")))
	doc.AddParagraph("Optional heading token: " + domain.Wrap(domain.LoopField("numberedTitle")))
	doc.AddParagraph("Concern: " + domain.Wrap(domain.LoopField("summary")))
	doc.AddParagraph(DotMarker + " [facts or evidence identified, grouped by procedures performed or data type].")
	doc.AddParagraph("Review of Communication Record")
	doc.AddParagraph(DotMarker + " Review of the Teams chats of the Relevant Colleague identified the following:")
	doc.AddParagraph("...")
	doc.AddParagraph("Colleague Interviews")
	doc.AddParagraph(DotMarker + " Colleague A stated that:")
	doc.AddParagraph(TildeMarker + " " + domain.Wrap(domain.LoopField("stmtLabelA")) + ". [key interview statements].")
	doc.AddParagraph(DotMarker + " Colleague B stated that:")
	doc.AddParagraph(TildeMarker + " " + domain.Wrap(domain.LoopField("stmtLabelB")) + ". [key interview statements].")
	doc.AddParagraph("Findings")
	doc.AddParagraph(
		DotMarker + " Concern " + domain.Wrap(domain.LoopField("listIndex")) + " was [not / substantiated] ...",
	)

	doc.AddParagraph(EndCaseIssuesRecordMarker)

	addBoldChapterTitle(doc, "CONCLUSION")
	doc.AddParagraph("Chapter number (auto): " + domain.Wrap("Conclusion.Number"))
	doc.AddParagraph(DotMarker + " Accordingly, the evidence supports that " + domain.Wrap("conclusion.support"))
	doc.AddParagraph(TildeMarker + " " + domain.Wrap("conclusion.policyBreaches"))

	addBoldChapterTitle(doc, "RECOMMENDATIONS")
	doc.AddParagraph("Chapter number (auto): " + domain.Wrap("Recommendations.Number"))
	doc.AddParagraph(DotMarker + " " + domain.Wrap("recommendation.1"))
	doc.AddParagraph(DotMarker + " " + domain.Wrap("recommendation.2"))
	doc.AddParagraph(DotMarker + " " + domain.Wrap("recommendation.3"))

	doc.AddParagraph("GROUP INVESTIGATIONS")
}
