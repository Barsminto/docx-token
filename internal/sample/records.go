package sample

import (
	"time"

	"github.com/liboyang/docx-token/internal/domain"
)

func Document(generatedAt time.Time) domain.Document {
	return domain.NewDocument("XXX", "Investigator", generatedAt, Items())
}

func Items() []domain.LineItemInput {
	return []domain.LineItemInput{
		{
			Title:          "CONCERN 1",
			ConcernBrief:   "Alleged sharing of confidential client information via personal email",
			ConcernSummary: "Alleged sharing of confidential client information via personal email",
		},
		{
			Title:          "CONCERN 2",
			ConcernBrief:   "Unapproved access to trading systems outside business hours",
			ConcernSummary: "Unapproved access to trading systems outside business hours",
		},
		{
			Title:          "CONCERN 3",
			ConcernBrief:   "Inappropriate language in Teams chats with external counterparties",
			ConcernSummary: "Inappropriate language in Teams chats with external counterparties",
		},
		{
			Title:          "CONCERN 4",
			ConcernBrief:   "Failure to escalate a suspected AML alert within required timelines",
			ConcernSummary: "Failure to escalate a suspected AML alert within required timelines",
		},
		{
			Title:          "CONCERN 5",
			ConcernBrief:   "Conflict of interest in vendor selection for IT services",
			ConcernSummary: "Conflict of interest in vendor selection for IT services",
		},
	}
}
