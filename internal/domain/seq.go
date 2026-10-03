package domain

import "fmt"

// AssignIssueSeq sets Record.Seq from the first number read from the template (continues existing numbering).
func (d *Document) AssignIssueSeq(firstSeq int) {
	for i := range d.Items {
		d.Items[i].Seq = fmt.Sprintf("%d", firstSeq+i)
		d.Items[i].NumberedTitle = fmt.Sprintf("%s. %s", d.Items[i].Seq, d.Items[i].Title)
	}
}
