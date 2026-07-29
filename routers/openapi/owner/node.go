package owner

import (
	"bunker-web/routers/openapi/owner/realms_server"
	"bunker-web/routers/openapi/owner/rental_server"
)

type Owner struct {
	rental_server.RentalServer
	realms_server.RealmsServer
}
