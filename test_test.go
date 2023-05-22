package mlogger

import (
	"context"
	"fmt"
	"github.com/Alp4ka/mlogger/contactpoints/matrix"
	"github.com/Alp4ka/mlogger/field"
	"github.com/Alp4ka/mlogger/misc"
	"github.com/Alp4ka/mlogger/templates"
	"testing"
	"time"
)

const t = `
My test template! 

**Time:**
*{{ .LogTime.Format "Jan 02, 2006 15:04:05 UTC" }}*

**Level:**
*{{ .LogLevel }}*

**Origin:**
*{{ .LogSource }}*
{{ if len .LogContextFields }}
**Context Fields:**
{{- end }}
{{ range .LogContextFields }}
{{ .Key }}: {{ .Value }}
{{ end }}
{{ if len .LogFields }}
**Fields:**
{{- end }}
{{ range .LogFields }}
{{ .Key }}: {{ .Value }}
{{ end }}
**Message:**
*{{ .LogMessage }}*`

func Test_Main(test *testing.T) {
	m := matrix.NewContactPoint(matrix.Config{})
	cfg := Config{Level: misc.LevelInfo,
		Template: templates.Config{Pattern: t, Use: true}}

	f1 := field.Bool("test_bool", true)
	ctx := field.WithContextFields(context.Background(), f1)
	logger, err := NewProduction(ctx, cfg, m)
	if err != nil {
		panic(err.Error())
	}

	ReplaceGlobals(logger)
	L().Info(
		"test message",
		field.Int("test_int", 123),
		field.String("test_string", "hello world!"),
		field.Error(fmt.Errorf("test_error")))

	time.Sleep(time.Second)
}
