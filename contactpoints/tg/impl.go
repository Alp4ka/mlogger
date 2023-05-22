package tg

import (
	"context"
	cp "github.com/Alp4ka/mlogger/contactpoints"
	"github.com/Alp4ka/mlogger/misc"
	"sync"
	"time"

	tb "gopkg.in/tucnak/telebot.v2"
)

const (
	_defaultPollTimeout = 10 * time.Second
)

type ContactPoint struct {
	cfg Config
	bot *tb.Bot

	once sync.Once
}

func (g *ContactPoint) init() error {
	var err error
	g.bot, err = tb.NewBot(tb.Settings{
		Token:  g.cfg.Token,
		Poller: &tb.LongPoller{Timeout: _defaultPollTimeout},
	})
	return err
}

func (g *ContactPoint) Msg(ctx context.Context, level misc.Level, msg string) error {
	if level < g.cfg.Level {
		return nil
	}

	var err error
	g.once.Do(
		func() {
			err = g.init()
		},
	)

	if err != nil {
		return err
	}

	_, err = g.bot.Send(
		g.cfg.ChatID,
		msg,
		&tb.SendOptions{DisableWebPagePreview: true, ParseMode: tb.ModeMarkdownV2},
	)

	return err
}

func NewContactPoint(cfg Config) *ContactPoint {
	return &ContactPoint{cfg: cfg}
}

var _ cp.ContactPoint = (*ContactPoint)(nil)
