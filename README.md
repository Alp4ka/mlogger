# MLogger
Logger with additional contact points and message templates based on zerolog for golang 1.20.*

## Contact points implementations:
1) Matrix
2) Telegram

## Parse Modes
1) Markdown

# Usage 

## Example
```golang
package main

import (
  "context"
  "fmt"
  "github.com/Alp4ka/mlogger/contactpoints/matrix"
  "github.com/Alp4ka/mlogger/field"
  "github.com/Alp4ka/mlogger/misc"
  "github.com/Alp4ka/mlogger/templates"
  "github.com/Alp4ka/mlogger"
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

func main() {
	cp := matrix.NewContactPoint(matrix.Config{
		Level:         misc.LevelInfo,
		HomeserverURL: "HOMESERVER_URL",
		UserID:        "USER_ID",
		RoomID:        "ROOM_ID",
		AccessToken:   "ACCESS_TOKEN",
	})

	cfg := Config{Level: misc.LevelInfo,
		Template: templates.Config{Pattern: t, Use: true}}

	f1 := field.Bool("test_bool", true)
	ctx := field.WithContextFields(context.Background(), f1)
	logger, err := mlogger.NewProduction(ctx, cfg, cp)
	if err != nil {
		panic(err.Error())
	}

	mlogger.ReplaceGlobals(logger)
	mlogger.L(ctx).Info(
		"test message",
		field.Int("test_int", 123),
		field.String("test_string", "hello world!"),
		field.Error(fmt.Errorf("test_error")))

	time.Sleep(time.Second)
}
```

## Output
![alt console output](https://github.com/Alp4ka/mlogger/blob/main/resources/log_template.png)

![alt matrix output](https://github.com/Alp4ka/mlogger/blob/main/resources/matrix_template.png)
