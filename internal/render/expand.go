package render

import (
	"archive/zip"
	"bytes"
	"fmt"
	"io"
	"strings"

	"github.com/liboyang/docx-token/internal/domain"
	"github.com/liboyang/docx-token/internal/template"
)

func expandRecordBlocks(docxBytes []byte, loopSpecs []template.LoopSpec, loopCounts map[string]int) ([]byte, error) {
	zr, err := zip.NewReader(bytes.NewReader(docxBytes), int64(len(docxBytes)))
	if err != nil {
		return nil, err
	}

	fileOrder := make([]string, 0, len(zr.File))
	files := make(map[string][]byte, len(zr.File))
	for _, f := range zr.File {
		rc, err := f.Open()
		if err != nil {
			return nil, err
		}
		data, err := io.ReadAll(rc)
		rc.Close()
		if err != nil {
			return nil, err
		}
		fileOrder = append(fileOrder, f.Name)
		files[f.Name] = data
	}

	raw, ok := files["word/document.xml"]
	if !ok {
		return nil, fmt.Errorf("word/document.xml not found")
	}
	expanded, err := expandAllLoops(string(raw), loopSpecs, loopCounts)
	if err != nil {
		return nil, err
	}
	files["word/document.xml"] = []byte(expanded)
	if num, ok := files["word/numbering.xml"]; ok {
		files["word/numbering.xml"] = []byte(mergeLoopLetterNumberingDefs(string(num), loopSpecs, loopCounts))
	}

	buf := new(bytes.Buffer)
	zw := zip.NewWriter(buf)
	for _, name := range fileOrder {
		w, err := zw.Create(name)
		if err != nil {
			return nil, err
		}
		if _, err := w.Write(files[name]); err != nil {
			return nil, err
		}
	}
	if err := zw.Close(); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func expandAllLoops(xml string, specs []template.LoopSpec, counts map[string]int) (string, error) {
	out := xml
	for si, spec := range specs {
		if !strings.Contains(out, spec.Begin) {
			continue
		}
		n := counts[spec.Begin]
		next, err := expandLoopRegion(out, spec, n, si)
		if err != nil {
			return "", fmt.Errorf("%s: %w", spec.Begin, err)
		}
		out = next
	}
	return out, nil
}

// expandLoopRegion duplicates paragraphs between markers; preserves all XML (styles, static subheadings).
func expandLoopRegion(xml string, spec template.LoopSpec, recordCount int, specIndex int) (string, error) {
	begin, end := spec.Begin, spec.End
	prefix := spec.ItemPrefix
	if prefix == "" {
		prefix = domain.LoopItemPrefix
	}
	indices := paragraphRE.FindAllStringIndex(xml, -1)
	if len(indices) == 0 {
		return "", fmt.Errorf("document has no paragraphs")
	}
	paragraphs := paragraphRE.FindAllString(xml, -1)

	beginIdx, endIdx := -1, -1
	for i, p := range paragraphs {
		if strings.Contains(p, begin) {
			beginIdx = i
		}
		if strings.Contains(p, end) {
			endIdx = i
		}
	}
	if beginIdx < 0 || endIdx < 0 || endIdx <= beginIdx {
		return "", fmt.Errorf("markers %s / %s not found or invalid order", begin, end)
	}

	block := strings.Join(paragraphs[beginIdx+1:endIdx], "")
	// Repair only inside the loop block so body footnote/endnote runs are not stripped.
	block = repairFragmentedPlaceholders(block)
	loopToken := "{" + prefix + "."
	if strings.Contains(block, loopToken) {
		var repeated strings.Builder
		for i := 0; i < recordCount; i++ {
			chunk := strings.ReplaceAll(block, loopToken, fmt.Sprintf("{%s_%d.", prefix, i))
			chunk = applyLoopCopyLetterNumID(chunk, loopLetterNumID(specIndex, i))
			repeated.WriteString(chunk)
		}
		block = repeated.String()
	}

	head := xml[:indices[beginIdx][0]]
	suffix := xml[indices[endIdx][1]:]
	return head + block + suffix, nil
}
