package tg

import (
	"github.com/Alp4ka/mlogger/misc"
	tb "gopkg.in/tucnak/telebot.v2"
)

type Config struct {
	Level  misc.Level
	ChatID tb.ChatID
	Token  string
}
