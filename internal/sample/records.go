package sample

import (
	"time"

	"github.com/liboyang/docx-token/internal/domain"
)

func Document(generatedAt time.Time) domain.Document {
	return domain.NewDocument("Business Data Report", "System Generated", generatedAt, Items())
}

func Items() []domain.LineItemInput {
	return []domain.LineItemInput{
		{ID: 1, Title: "INC-1001 Hardware refresh", Category: "Hardware"},
		{ID: 2, Title: "INC-1002 License renewal", Category: "Software"},
		{ID: 3, Title: "INC-1003 Cloud storage", Category: "Cloud"},
	}
}
