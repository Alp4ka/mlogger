package templates

import (
	"bytes"
	"github.com/Alp4ka/mlogger/misc"
	"text/template"
)

const (
	_tmplName = ""
)

// Template is a structure that represents template of a definite Mode.
type Template struct {
	Tmpl *template.Template
	Mode misc.ParseMode
}

func (tmpl *Template) Render(placeholder *Placeholder) (string, error) {
	buf := new(bytes.Buffer)
	err := tmpl.Tmpl.Execute(buf, *placeholder)
	if err != nil {
		return "", err
	}
	return buf.String(), nil
}

func FromPattern(pattern string) (*Template, error) {
	tmpl, err := template.New(_tmplName).Parse(pattern)
	if err != nil {
		return nil, err
	}

	return &Template{Tmpl: tmpl, Mode: misc.DefaultMode}, nil
}
