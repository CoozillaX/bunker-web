package routers

import (
	"bunker-web/routers/api"
	"bunker-web/routers/openapi"
)

type Routers struct {
	api.API
	openapi.OpenAPI
}

var routers Routers
