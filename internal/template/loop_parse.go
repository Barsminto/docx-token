package template

import (
	"fmt"
	"regexp"
	"strings"
)

var loopMarkerRE = regexp.MustCompile(`\{(BEGIN|END):([^}]+)\}`)

// LoopSpec is parsed from {BEGIN:case.case_issues:record} … {END:case.case_issues:record}.
type LoopSpec struct {
	Begin      string // full marker, e.g. {BEGIN:case.issues:record}
	End        string
	DataPath   string // dot path in map, e.g. case.case_issues
	ItemPrefix string // loop variable, e.g. record → {record.title} → {record_0.title}
}

// ParseLoopMarkerBody parses the inside of BEGIN/END, e.g. "case.case_issues:record".
func ParseLoopMarkerBody(body string) (dataPath, itemPrefix string) {
	body = strings.TrimSpace(body)
	if body == "" {
		return "", "item"
	}

	// Legacy GI markers: case.issues[:background|record] — data always case.issues, prefix issue.
	if strings.HasPrefix(body, "case.issues") {
		suffix := strings.TrimPrefix(body, "case.issues")
		if suffix == "" {
			return "case.issues", "issue"
		}
		if suffix == ":background" || suffix == ":record" {
			return "case.issues", "issue"
		}
	}

	idx := strings.LastIndex(body, ":")
	if idx > 0 {
		suffix := body[idx+1:]
		if isLoopPrefixIdent(suffix) {
			return body[:idx], suffix
		}
	}

	// Path only: default prefix from last segment (issues → issue).
	parts := strings.Split(body, ".")
	last := parts[len(parts)-1]
	prefix := last
	if last == "issues" {
		prefix = "issue"
	} else if strings.HasSuffix(last, "s") && len(last) > 1 {
		prefix = strings.TrimSuffix(last, "s")
	}
	return body, prefix
}

func isLoopPrefixIdent(s string) bool {
	if s == "" {
		return false
	}
	for i, r := range s {
		if i == 0 {
			if r != '_' && !isASCIILetter(r) {
				return false
			}
			continue
		}
		if r != '_' && !isASCIILetter(r) && (r < '0' || r > '9') {
			return false
		}
	}
	return true
}

func isASCIILetter(r rune) bool {
	return (r >= 'a' && r <= 'z') || (r >= 'A' && r <= 'Z')
}

// LoopSpecFromMarker parses {BEGIN:…} or {END:…}.
func LoopSpecFromMarker(marker string) (LoopSpec, error) {
	m := loopMarkerRE.FindStringSubmatch(marker)
	if len(m) != 3 {
		return LoopSpec{}, fmt.Errorf("not a loop marker: %s", marker)
	}
	body := m[2]
	path, prefix := ParseLoopMarkerBody(body)
	full := "{" + m[1] + ":" + body + "}"
	var begin, end string
	if m[1] == "BEGIN" {
		begin = full
		end = "{END:" + body + "}"
	} else {
		end = full
		begin = "{BEGIN:" + body + "}"
	}
	return LoopSpec{
		Begin:      begin,
		End:        end,
		DataPath:   path,
		ItemPrefix: prefix,
	}, nil
}

// DiscoverLoopSpecs finds unique loop regions in document XML (paragraph order preserved).
func DiscoverLoopSpecs(documentXML string) ([]LoopSpec, error) {
	seen := map[string]bool{}
	var specs []LoopSpec
	for _, m := range loopMarkerRE.FindAllStringSubmatch(documentXML, -1) {
		if m[1] != "BEGIN" {
			continue
		}
		spec, err := LoopSpecFromMarker(m[0])
		if err != nil {
			return nil, err
		}
		if seen[spec.Begin] {
			continue
		}
		seen[spec.Begin] = true
		if !strings.Contains(documentXML, spec.End) {
			return nil, fmt.Errorf("missing end marker %s for %s", spec.End, spec.Begin)
		}
		specs = append(specs, spec)
	}
	if len(specs) == 0 {
		for _, r := range LoopRegions {
			body := strings.TrimSuffix(strings.TrimPrefix(r.Begin, "{BEGIN:"), "}")
			path, prefix := ParseLoopMarkerBody(body)
			specs = append(specs, LoopSpec{
				Begin:      r.Begin,
				End:        r.End,
				DataPath:   path,
				ItemPrefix: prefix,
			})
		}
	}
	return specs, nil
}
