package admin

import (
	"bunker-web/routers/api/admin/redeem_code"
	"bunker-web/routers/api/admin/unlimited_server"
	"bunker-web/routers/api/admin/user"
)

type Admin struct {
	redeem_code.RedeemCode
	unlimited_server.UnlimitedServer
	user.User
}
