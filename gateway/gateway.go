package gateway

import (
	"context"
	"github.com/Alp4ka/mlogger/contactpoints"
	"github.com/Alp4ka/mlogger/field"
	"github.com/Alp4ka/mlogger/misc"
	"github.com/Alp4ka/mlogger/templates"
	"golang.org/x/sync/errgroup"
	"time"
)

type Gateway struct {
	template *templates.Template
	contacts []contactpoints.ContactPoint
}

func (g Gateway) Msg(ctx context.Context, level misc.Level, msg string, fields ...field.Field) error {
	rendered, err := g.template.Render(
		&templates.Placeholder{
			LogTime:          time.Now(),
			LogLevel:         level,
			LogSource:        "None",
			LogContextFields: field.FieldsFromCtx(ctx),
			LogFields:        fields,
			LogMessage:       msg,
		},
	)
	if err != nil {
		return err
	}

	eg := &errgroup.Group{}
	for _, cp := range g.contacts {
		if cp != nil {
			var c contactpoints.ContactPoint
			c = cp
			eg.Go(func() error {
				return c.Msg(ctx, level, rendered)
			})
		}
	}
	return eg.Wait()
}

func (g Gateway) WithTemplate(tmpl *templates.Template) Gateway {
	g.template = tmpl
	return g
}

func (g Gateway) WithContactPoints(replace bool, contacts ...contactpoints.ContactPoint) Gateway {
	if replace {
		g.contacts = contacts
	} else {
		g.contacts = append(g.contacts, contacts...)
	}
	return g
}

func CreateGateway() Gateway {
	return Gateway{
		template: templates.DefaultTemplate(misc.DefaultMode),
	}
}
