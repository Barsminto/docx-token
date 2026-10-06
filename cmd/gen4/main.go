package main

import (
	"fmt"
	"time"

	"github.com/liboyang/docx-token/docxtoken"
	"github.com/liboyang/docx-token/internal/project"
	"github.com/liboyang/docx-token/internal/sample"
)

func main() {
	now := time.Now()
	data := sample.FillMap(now)
	caseMap, _ := data["case"].(map[string]interface{})
	issues, _ := caseMap["issues"].([]map[string]interface{})
	caseMap["issues"] = issues[:4]
	tpl, err := project.TemplatePath()
	if err != nil {
		panic(err)
	}
	out := "output/report-4-concerns.docx"
	if err := docxtoken.Fill(tpl, out, data); err != nil {
		panic(err)
	}
	fmt.Println("Generated:", out)
}
