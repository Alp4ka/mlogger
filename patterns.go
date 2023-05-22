package templates

import (
	"github.com/Alp4ka/mlogger/field"
	"github.com/Alp4ka/mlogger/misc"
	"time"
)

const (
	markdownTmplPattern = `**Time:**
*{{ .LogTime.Format "Jan 02, 2006 15:04:05 UTC" }}*

**Level:**
*{{ .LogLevel }}*

**Origin:**
*{{ .LogSource }}*
{{ if len .LogContextFields }}
**Context Fields:**
{{- end }}
{{ range .LogContextFields }}
*{{ .Key }}*: {{ .Value }}
{{ end }}
{{ if len .LogFields }}
**Fields:**
{{- end }}
{{ range .LogFields }}
*{{ .Key }}*: {{ .Value }}
{{ end }}
**Message:**
*{{ .LogMessage }}*`
)

type Placeholder struct {
	LogTime          time.Time
	LogLevel         misc.Level
	LogSource        string
	LogContextFields field.Fields
	LogFields        []field.Field
	LogMessage       string
}
