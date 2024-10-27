package main

import (
	"context"
	"fmt"
	"os"
	"time"

	"olshop1/conf"
	"olshop1/domain"
	"olshop1/presentation"

	"github.com/kokizzu/gotro/D/Pg"
	"github.com/kokizzu/gotro/L"
	"github.com/kokizzu/gotro/M"
	"github.com/kokizzu/gotro/S"
	"github.com/rs/zerolog"
	"golang.org/x/sync/errgroup"
)

var VERSION = ``
var log *zerolog.Logger

func main() {
	conf.VERSION = VERSION

	// note: set instance id when there's multiple instance
	// lexid.Config.Separator = `~THE_INSTANCE_ID`

	fmt.Println(conf.PROJECT_NAME + ` ` + S.IfEmpty(VERSION, `local-dev`))

	log = conf.InitLogger()
	conf.LoadEnv()

	args := os.Args
	if len(args) < 2 {
		L.Print(`must start with: run, web, cron, migrate, or config as first argument`)
		L.Print(args)
		return
	}

	if args[1] == `config` {
		L.Describe(M.SX{
			`web`:      conf.EnvWebConf(),
			`postgres`: conf.EnvPostgres(),
		})
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	eg, ctx := errgroup.WithContext(ctx)
	var closers []func() error

	// connect to postgres
	var pConn *Pg.Adapter
	eg.Go(func() error {
		pConf := conf.EnvPostgres()
		pConn = pConf.Connect()
	})

	L.PanicIf(eg.Wait(), `eg.Wait`) // if error, make sure no error on: docker compose up
	for _, closer := range closers {
		closer := closer
		defer closer()
	}

	// create domain object
	d := &domain.Domain{
		IsBgSvc: false,
		Log:     log,

		UploadDir: conf.UploadDir(),
		CacheDir:  conf.CacheDir(),

		Superadmins: conf.EnvSuperAdmins(),

		Postgres: pConn,

		WebCfg: conf.EnvWebConf(),
	}

	mode := S.ToLower(os.Args[1])

	// check table existence
	if mode != `migrate` {
		L.Print(`verifying table schema, if failed, run: go run main.go migrate`)
		// TODO: migrate tables
	}

	// start
	switch mode {
	case `web`:
		ws := &presentation.WebServer{
			Domain: d,
		}
		conf.LoadCountries("./static/country_data/countries.tsv")
		ws.Start(log)
	case `cli`:
		cli := &presentation.CLI{
			Domain: d,
		}
		cli.Run(os.Args[2:], log)
	case `cron`:
		cron := &presentation.Cron{
			Domain: d,
		}
		cron.Start(log)
	case `migrate`:
		model.RunMigration(log, pConn, cConn, pConn, cConn, pConn)

	default:
		log.Error().Str(`mode`, mode).Msg(`unknown mode`)
	}

}
