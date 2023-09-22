package mlogger

import (
	"context"
	"github.com/Alp4ka/mlogger/contactpoints"
	"github.com/Alp4ka/mlogger/field"
	"github.com/Alp4ka/mlogger/gateway"
	"github.com/Alp4ka/mlogger/jsonsecurity"
	"github.com/Alp4ka/mlogger/misc"
	"github.com/Alp4ka/mlogger/templates"
	"io"
	"log/slog"
	"os"
	"sync"
)

const (
	callerFuncLevelShift = 7
)

var (
	_globalMu sync.RWMutex
	_globalL  *MainLogger

	_defaultWriter = os.Stdout
)

type MainLogger struct {
	cfg *Config

	gw     gateway.Gateway
	ctx    context.Context
	masker *jsonsecurity.Masker

	logger *slog.Logger
}

func (l *MainLogger) log(level misc.Level, msg string, fields ...field.Field) {
	if level.LessThan(l.cfg.Level) {
		return
	}

	// log is located
	flds := field.FieldsFromContext(field.WithContextFields(l.ctx, fields...)).
		WithOptions(field.OptionCallerFunc(callerFuncLevelShift)).
		Prepare(l.masker)

	attrs := field.UnpackFieldsToSlogAttrs(flds)

	go func() {
		if err := l.gw.Msg(l.ctx, l.cfg.Source, misc.LevelDebug, msg, flds...); err != nil {
			l.logger.LogAttrs(l.ctx, misc.LevelWarn, err.Error(), attrs...)
		}
	}()

	l.logger.LogAttrs(l.ctx, misc.SlogLevel(level), msg, attrs...)
}

func (l *MainLogger) Log(level misc.Level, msg string, fields ...field.Field) {
	l.log(level, msg, fields...)
}

func (l *MainLogger) Debug(msg string, fields ...field.Field) {
	l.log(misc.LevelDebug, msg, fields...)
}

func (l *MainLogger) Info(msg string, fields ...field.Field) {
	l.log(misc.LevelInfo, msg, fields...)
}

func (l *MainLogger) Warn(msg string, fields ...field.Field) {
	l.log(misc.LevelWarn, msg, fields...)
}

func (l *MainLogger) Error(msg string, fields ...field.Field) {
	l.log(misc.LevelError, msg, fields...)
}

func (l *MainLogger) Fatal(msg string, fields ...field.Field) {
	l.log(misc.LevelFatal, msg, fields...)
}

func (l *MainLogger) Panic(msg string, fields ...field.Field) {
	l.log(misc.LevelPanic, msg, fields...)
}

func L(optionalCtx ...context.Context) *MainLogger {
	var l *MainLogger
	_globalMu.RLock()

	if len(optionalCtx) != 0 {
		l = &MainLogger{
			_globalL.cfg,
			_globalL.gw,
			field.WithContextFields(_globalL.ctx, field.FieldsFromContext(optionalCtx[0])...),
			_globalL.masker,
			_globalL.logger,
		}
	} else {
		l = _globalL
	}

	_globalMu.RUnlock()
	return l
}

func NewProduction(ctx context.Context, cfg Config, contacts ...contactpoints.ContactPoint) (*MainLogger, error) {
	var (
		gw   gateway.Gateway
		tmpl *templates.Template
		err  error
	)

	masker, err := jsonsecurity.NewMasker(cfg.JSONSecurity)
	if err != nil {
		return nil, err
	}

	if !cfg.Template.Use {
		tmpl = templates.DefaultTemplate(misc.DefaultMode)
	} else if tmpl, err = templates.FromPattern(cfg.Template.Pattern); err != nil || tmpl == nil {
		return nil, err
	}
	gw = gateway.CreateGateway().WithTemplate(tmpl).WithContactPoints(true, contacts...)

	// TODO: (???)
	cfg.Writer = misc.Coalesce[io.Writer](cfg.Writer, _defaultWriter)

	fields := field.FieldsFromContext(ctx)
	if cfg.Source != "" {
		fields = fields.WithOptions(field.OptionSource(cfg.Source))
	}
	fields = fields.Prepare(masker)

	return &MainLogger{
		&cfg,
		gw,
		field.WithContextFields(ctx, fields...),
		masker,
		slog.New(slog.NewJSONHandler(cfg.Writer, &slog.HandlerOptions{
			Level:       misc.SlogLevel(cfg.Level),
			ReplaceAttr: misc.SlogReplaceAttr(),
		})),
	}, nil

}

func ReplaceGlobals(logger *MainLogger) {
	_globalMu.Lock()
	_globalL = logger
	_globalMu.Unlock()
}
