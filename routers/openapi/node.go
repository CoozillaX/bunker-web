package openapi

import (
	"bunker-web/routers/openapi/helper"
	"bunker-web/routers/openapi/owner"
	"bunker-web/routers/openapi/user"
)

type OpenAPI struct {
	helper.Helper
	owner.Owner
	user.User
}
