package mlogger

import (
	"context"
	"github.com/Alp4ka/mlogger/contactpoints"
	"github.com/Alp4ka/mlogger/field"
	"github.com/Alp4ka/mlogger/gateway"
	"github.com/Alp4ka/mlogger/misc"
	"github.com/Alp4ka/mlogger/templates"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/pkgerrors"
	"os"
	"sync"
)

var (
	_globalMu sync.RWMutex
	_globalL  *MainLogger
)

type MainLogger struct {
	cfg Config
	gw  gateway.Gateway
	ctx context.Context

	logger zerolog.Logger
}

func (l *MainLogger) Info(msg string, fields ...field.Field) {
	if misc.LevelInfo < l.cfg.Level {
		return
	}

	logger := l.logger.With().Fields(field.Fields(fields).Map()).Logger()
	go func() {
		if err := l.gw.Msg(l.ctx, misc.LevelInfo, msg, fields...); err != nil {
			logger.Warn().Msg(err.Error())
		}
	}()
	logger.Info().Msg(msg)
}

func (l *MainLogger) Debug(msg string, fields ...field.Field) {
	if misc.LevelDebug < l.cfg.Level {
		return
	}

	logger := l.logger.With().Fields(field.Fields(fields).Map()).Logger()
	go func() {
		if err := l.gw.Msg(l.ctx, misc.LevelDebug, msg, fields...); err != nil {
			logger.Warn().Msg(err.Error())
		}
	}()
	logger.Debug().Msg(msg)
}

func (l *MainLogger) Error(msg string, fields ...field.Field) {
	if misc.LevelError < l.cfg.Level {
		return
	}

	logger := l.logger.With().Fields(field.Fields(fields).Map()).Logger()
	go func() {
		if err := l.gw.Msg(l.ctx, misc.LevelError, msg, fields...); err != nil {
			logger.Warn().Msg(err.Error())
		}
	}()
	logger.Error().Msg(msg)
}

func (l *MainLogger) Fatal(msg string, fields ...field.Field) {
	if misc.LevelFatal < l.cfg.Level {
		return
	}

	logger := l.logger.With().Fields(field.Fields(fields).Map()).Logger()
	go func() {
		if err := l.gw.Msg(l.ctx, misc.LevelFatal, msg, fields...); err != nil {
			logger.Warn().Msg(err.Error())
		}
	}()
	logger.Fatal().Msg(msg)
}

func (l *MainLogger) Warn(msg string, fields ...field.Field) {
	if misc.LevelWarn < l.cfg.Level {
		return
	}

	logger := l.logger.With().Fields(field.Fields(fields).Map()).Logger()
	logger.Warn().Msg(msg)
	go func() {
		if err := l.gw.Msg(l.ctx, misc.LevelWarn, msg, fields...); err != nil {
			logger.Warn().Msg(err.Error())
		}
	}()
}

func (l *MainLogger) Panic(msg string, fields ...field.Field) {
	if misc.LevelPanic < l.cfg.Level {
		return
	}

	logger := l.logger.With().Fields(field.Fields(fields).Map()).Logger()
	go func() {
		if err := l.gw.Msg(l.ctx, misc.LevelPanic, msg, fields...); err != nil {
			logger.Warn().Msg(err.Error())
		}
	}()
	logger.Panic().Msg(msg)
}

func L(ctx context.Context) *MainLogger {
	var l *MainLogger
	_globalMu.RLock()

	l = &MainLogger{
		_globalL.cfg,
		_globalL.gw,
		ctx,
		_globalL.logger.
			With().
			Fields(field.FieldsFromCtx(ctx).Map()).
			Logger(),
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

	zerolog.ErrorStackMarshaler = pkgerrors.MarshalStack

	if !cfg.Template.Use {
		tmpl = templates.DefaultTemplate(misc.DefaultMode)
	} else if tmpl, err = templates.FromPattern(cfg.Template.Pattern); err != nil || tmpl == nil {
		return nil, err
	}

	gw = gateway.
		CreateGateway().
		WithTemplate(tmpl).
		WithContactPoints(true, contacts...)

	return &MainLogger{
		cfg,
		gw,
		ctx,
		zerolog.
			New(os.Stdout).
			Level(zerolog.InfoLevel).
			With().
			Fields(field.FieldsFromCtx(ctx).Map()).
			Logger(),
	}, nil
}

func ReplaceGlobals(logger *MainLogger) {
	_globalMu.Lock()
	_globalL = logger
	_globalMu.Unlock()
}
