package admin

import (
	"bunker-web/routers/api/admin/unlimited_server"
	"bunker-web/routers/api/admin/user"
)

type Admin struct {
	unlimited_server.UnlimitedServer
	user.User
}
