package matrix

import (
	"github.com/Alp4ka/mlogger/misc"
)

type Config struct {
	Level         misc.Level
	HomeserverURL string
	UserID        string
	RoomID        string
	AccessToken   string
}
