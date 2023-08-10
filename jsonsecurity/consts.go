package jsonsecurity

import "regexp"

const MaskSymbol = "*"

var _noWhitespaceRegexCompiled = regexp.MustCompile(`\S`)
