package api

import (
	"github.com/sazonovItas/gochat-tcp/cmd/gochat/app/core"
)

type Api struct {
	app *core.Core
}

func NewApi(core *core.Core) *Api {
	return &Api{app: core}
}
