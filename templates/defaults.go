package templates

import (
	"github.com/Alp4ka/mlogger/misc"
	"text/template"
)

var (
	_tmplMap = map[misc.ParseMode]*template.Template{
		misc.MarkdownMode: _markdownTmpl,
	}

	_markdownTmpl, _ = template.New(_tmplName).Parse(markdownTmplPattern)
)

func DefaultTemplate(mode misc.ParseMode) *Template {
	if val, ok := _tmplMap[mode]; ok {
		return &Template{Tmpl: val, Mode: mode}
	}
	return &Template{Tmpl: _tmplMap[misc.DefaultMode], Mode: misc.DefaultMode}
}
