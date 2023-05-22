package mlogger

import (
	"github.com/Alp4ka/mlogger/misc"
	"github.com/Alp4ka/mlogger/templates"
)

type Config struct {
	Level    misc.Level
	Template templates.Config
}
