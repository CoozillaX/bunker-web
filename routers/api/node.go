package api

import (
	"bunker-web/routers/api/admin"
	"bunker-web/routers/api/helper"
	"bunker-web/routers/api/notice"
	"bunker-web/routers/api/owner"
	"bunker-web/routers/api/phoenix"
	"bunker-web/routers/api/user"
	"bunker-web/routers/api/webauthn"
)

type API struct {
	admin.Admin
	notice.Notice
	helper.Helper
	owner.Owner
	phoenix.Phoenix
	user.User
	webauthn.Webauthn
}
