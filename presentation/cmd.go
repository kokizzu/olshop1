package presentation

import (
	"olshop1/domain"

	"github.com/rs/zerolog"
)

type CLI struct {
	*domain.Domain
}

func (c *CLI) Run(args []string, log *zerolog.Logger) {
	if len(args) < 1 {
		c.Log.Print(`must start with one of: `, allCommands)
		return
	}
	if len(args) < 2 {
		c.Log.Print(`must provide json payload`)
		return
	}

	cmdRun(c.Domain, args[0], []byte(args[1]))
}
