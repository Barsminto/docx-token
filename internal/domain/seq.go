package domain

import "fmt"

// AssignIssueSeq sets per-concern section numbers from the template (continues after "1.", "2.", …).
// SectionNumber equals the main investigation chapter (3, 4, 5, …) used in "{Record.SectionNumber}.2" lines.
func (d *Document) AssignIssueSeq(firstSeq int) {
	for i := range d.Items {
		seq := firstSeq + i
		d.Items[i].Seq = fmt.Sprintf("%d", seq)
		d.Items[i].ListIndex = fmt.Sprintf("%d", i+1)
		d.Items[i].ListLetter = string(rune('a' + i))
		d.Items[i].StmtLabelA = string(rune('a' + 2*i))
		d.Items[i].StmtLabelB = string(rune('b' + 2*i))
		d.Items[i].SectionNumber = d.Items[i].Seq
		title := d.Items[i].Title
		if title == "" {
			title = fmt.Sprintf("CONCERN %d", i+1)
		}
		d.Items[i].NumberedTitle = fmt.Sprintf("%s. INVESTIGATION: %s", d.Items[i].Seq, title)
	}
}
