package conf

import (
	"os"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/kokizzu/gotro/D/Pg"
	"github.com/kokizzu/gotro/S"
)

type PostgresConf struct {
	User   string
	Pass   string
	Host   string
	Port   int
	DbName string
}

func EnvPostgres() PostgresConf {
	// postgresql://root@127.0.0.1:26257/defaultdb?sslmode=disable
	pgUser := os.Getenv("POSTGRES_USER")
	if pgUser == `` {
		pgUser = `root`
	}
	pgHost := os.Getenv("POSTGRES_HOST")
	if pgHost == `` {
		pgHost = `127.0.0.1`
	}
	pgPort := S.ToInt(os.Getenv("POSTGRES_PORT"))
	if pgPort <= 0 || pgPort > 65535 {
		pgPort = 26257
	}
	pgDbName := os.Getenv("POSTGRES_DBNAME")
	if pgDbName == `` {
		pgDbName = `defaultdb`
	}

	return PostgresConf{
		User:   pgUser,
		Pass:   os.Getenv("POSTGRES_PASS"),
		Host:   pgHost,
		Port:   pgPort,
		DbName: pgDbName,
	}
}
func (p PostgresConf) Connect() *Pg.Adapter {
	adapter := &Pg.Adapter{
		Reconnect: func() *pgxpool.Pool {
			return Pg.Connect1(p.User, p.Pass, p.Host, p.DbName, p.Port, 32)
		},
	}
	adapter.Pool = adapter.Reconnect()
	return adapter
}
