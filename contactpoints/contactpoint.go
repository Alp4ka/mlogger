package contactpoints

import (
	"context"
	"github.com/Alp4ka/mlogger/misc"
)

type ContactPoint interface {
	Msg(ctx context.Context, level misc.Level, msg string) error
}
