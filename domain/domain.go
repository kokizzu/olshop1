package domain

import (
	"olshop1/conf"

	"github.com/kokizzu/gotro/D/Pg"
	"github.com/kokizzu/gotro/M"
	"github.com/rs/zerolog"
)

type Domain struct {
	IsBgSvc     bool
	Log         *zerolog.Logger
	UploadDir   string
	CacheDir    string
	Superadmins M.SB

	WebCfg conf.WebConf

	Postgres *Pg.Adapter
}
