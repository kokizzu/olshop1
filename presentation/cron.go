package presentation

import (
	"github.com/rs/zerolog"

	"olshop1/domain"
)

type Cron struct {
	*domain.Domain
	Log *zerolog.Logger
}

func (c *Cron) Start(*zerolog.Logger) {

}
