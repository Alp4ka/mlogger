package mlogger

import (
	"github.com/Alp4ka/mlogger/jsonsecurity"
	"github.com/Alp4ka/mlogger/misc"
	"github.com/Alp4ka/mlogger/templates"
	"io"
)

type Config struct {
	Source       string
	Level        misc.Level
	ShowLogUUID  bool
	Template     templates.Config
	JSONSecurity jsonsecurity.Config
	Writer       io.Writer
}
