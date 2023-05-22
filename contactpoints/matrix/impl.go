package matrix

import (
	"context"
	cp "github.com/Alp4ka/mlogger/contactpoints"
	"github.com/Alp4ka/mlogger/misc"
	mb "gitlab.com/silkeh/matrix-bot"
	"sync"
)

type ContactPoint struct {
	cfg Config
	bot *mb.Client

	once sync.Once
}

func (g *ContactPoint) init() error {
	var err error

	config := &mb.ClientConfig{AllowedRooms: []string{g.cfg.RoomID}}
	g.bot, err = mb.NewClient(
		g.cfg.HomeserverURL,
		g.cfg.UserID,
		g.cfg.AccessToken,
		config,
	)

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

	_, err = g.bot.NewRoom(g.cfg.RoomID).SendMarkdown(msg)
	return err
}

func NewContactPoint(cfg Config) *ContactPoint {
	return &ContactPoint{cfg: cfg}
}

var _ cp.ContactPoint = (*ContactPoint)(nil)
