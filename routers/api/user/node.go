package user

import (
	"bunker-web/routers/api/user/api_key"
	"bunker-web/routers/api/user/email"
)

type User struct {
	api_key.APIKey
	email.Email
}
