package webauthn

import (
	"bunker-web/routers/api/webauthn/login"
	"bunker-web/routers/api/webauthn/register"
)

type Webauthn struct {
	login.Login
	register.Register
}
