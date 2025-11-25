package webauthner

import (
	"bunker-web/configs"
	"bunker-web/models"
	"bunker-web/pkg/giner"
	"bunker-web/services/user"
	"bunker-web/services/webauthn_credential"
	"encoding/binary"
	"errors"
	"net/http"
	"net/url"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-webauthn/webauthn/protocol"
	"github.com/go-webauthn/webauthn/webauthn"
	"github.com/google/uuid"
	"github.com/patrickmn/go-cache"
)

const (
	WEBAUTHN_REQID_COOKIE_NAME = "BUNKER_WEB_WEBAUTHN_REQID"
	WEBAUTHN_REQ_EXPIRE_TIME   = time.Minute
)

var (
	webAuthn             *webauthn.WebAuthn
	webAuthnSessionCache *cache.Cache
)

func init() {
	u, _ := url.Parse(configs.CURRENT_WEB_DOMAIN)
	webAuthn, _ = webauthn.New(&webauthn.Config{
		RPDisplayName: "BunkerWeb",
		RPID:          u.Hostname(),
		RPOrigins:     []string{configs.CURRENT_WEB_DOMAIN},
		Timeouts: webauthn.TimeoutsConfig{
			Login: webauthn.TimeoutConfig{
				Enforce:    true,
				Timeout:    WEBAUTHN_REQ_EXPIRE_TIME,
				TimeoutUVD: WEBAUTHN_REQ_EXPIRE_TIME,
			},
			Registration: webauthn.TimeoutConfig{
				Enforce:    true,
				Timeout:    WEBAUTHN_REQ_EXPIRE_TIME,
				TimeoutUVD: WEBAUTHN_REQ_EXPIRE_TIME,
			},
		},
	})
	webAuthnSessionCache = cache.New(WEBAUTHN_REQ_EXPIRE_TIME, time.Second*10)
}

func BeginRegistration(bearer string, user *models.User) (*protocol.CredentialCreation, error) {
	// 1. Begin the registration process
	options, session, err := webAuthn.BeginRegistration(user, func(pkcco *protocol.PublicKeyCredentialCreationOptions) {
		if credentials, err := webauthn_credential.QueryRawsByUserID(user.ID); err == nil {
			for _, credential := range credentials {
				pkcco.CredentialExcludeList = append(pkcco.CredentialExcludeList, credential.Descriptor())
			}
		}
		pkcco.AuthenticatorSelection.ResidentKey = protocol.ResidentKeyRequirementPreferred
		pkcco.Attestation = protocol.PreferNoAttestation
		pkcco.Extensions = protocol.AuthenticationExtensions{
			"credProps": true,
		}
	})
	if err != nil {
		return nil, err
	}
	// 2. Store the session data
	webAuthnSessionCache.SetDefault(bearer, session)
	// 3. Return options string
	return options, nil
}

func FinishRegistration(bearer string, user *models.User, response *http.Request) *gin.Error {
	// 1. Get the session data
	sess, ok := webAuthnSessionCache.Get(bearer)
	if !ok {
		return giner.NewPublicGinError("无效请求")
	}
	session := sess.(*webauthn.SessionData)
	// 2. Finish the registration process
	credential, err := webAuthn.FinishRegistration(user, *session, response)
	if err != nil {
		return giner.NewPublicGinError(err.Error())
	}
	// 3. Check if the credential already exists
	if _, err := webauthn_credential.QueryModelByRawID(credential.ID); err == nil {
		return giner.NewPublicGinError("请重试")
	}
	// 4. Store the credential
	return giner.NewPrivateGinError(webauthn_credential.StoreToDB(credential, user.ID))
}

func BeginDiscoverableLogin() (string, *protocol.CredentialAssertion, error) {
	// 1. Begin the login process
	options, session, err := webAuthn.BeginDiscoverableLogin()
	if err != nil {
		return "", nil, err
	}
	// 2. Store the session data
	reqId := uuid.NewString()
	webAuthnSessionCache.SetDefault(reqId, session)
	// 3. Return options
	return reqId, options, nil
}

func FinishDiscoverableLogin(reqId string, response *http.Request) (*models.User, *gin.Error) {
	// 1. Get the session data
	sess, ok := webAuthnSessionCache.Get(reqId)
	if !ok {
		return nil, giner.NewPublicGinError("无效请求")
	}
	session := sess.(*webauthn.SessionData)
	// 2. Finish the login process
	credential, err := webAuthn.FinishDiscoverableLogin(
		func(rawID, userHandle []byte) (_ webauthn.User, err error) {
			// Try to use user id to find user
			// Read user id
			userID := uint(binary.BigEndian.Uint32(userHandle))
			// Get and return user
			usr, ginerr := user.QueryUserByID(userID)
			if ginerr == nil && usr != nil {
				return usr, nil
			}
			// If not found, use credential id to find user
			credential, err := webauthn_credential.QueryModelByRawID(rawID)
			if err != nil {
				return nil, err
			}
			usr, ginerr = user.QueryUserByID(credential.UserID)
			if ginerr != nil {
				return nil, errors.New("请重试")
			}
			return usr, nil
		},
		*session,
		response,
	)
	if err != nil {
		return nil, giner.NewPublicGinError(err.Error())
	}
	// 3. Get credential model by auth response
	credentialm, err := webauthn_credential.QueryModelByRawID(credential.ID)
	if err != nil {
		return nil, giner.NewPrivateGinError(err)
	}
	// 4. Get user
	return user.QueryUserByID(credentialm.UserID)
}
