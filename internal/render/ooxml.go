package render

import "regexp"

// paragraphRE matches a full w:p element (styles live in pPr; any formatting is preserved when copying).
var paragraphRE = regexp.MustCompile(`(?s)<w:p\b[^>]*>.*?</w:p>`)
