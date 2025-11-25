package sessions

import (
	"sync"
	"time"

	"github.com/patrickmn/go-cache"
)

const (
	SESSION_EXPIRE_TIME = time.Hour * 24 * 7
	SESSION_COOKIE_NAME = "BUNKER_WEB_SESSION"
)

var (
	sessions *cache.Cache
	users    *cache.Cache
)

func init() {
	sessions = cache.New(SESSION_EXPIRE_TIME, time.Minute*5)
	users = cache.New(SESSION_EXPIRE_TIME, time.Minute*5)
}

func CreateSessionByBearer(bearer string) *sync.Map {
	session := &sync.Map{}
	sessions.Add(bearer, session, SESSION_EXPIRE_TIME)
	return session
}

func GetSessionByBearer(bearer string) (*sync.Map, bool) {
	session, ok := sessions.Get(bearer)
	if !ok {
		return nil, false
	}
	sessionMap, ok := session.(*sync.Map)
	if !ok {
		return nil, false
	}
	username, ok := sessionMap.Load("session_username")
	if ok {
		users.Replace(username.(string), bearer, SESSION_EXPIRE_TIME)
	}
	sessions.Replace(bearer, session, SESSION_EXPIRE_TIME)
	return sessionMap, true
}

func BindSessionToUsername(bearer, username string) bool {
	session, ok := GetSessionByBearer(bearer)
	if !ok {
		return false
	}
	DeleteSessionByUsername(username)
	session.Store("session_username", username)
	users.Add(username, bearer, SESSION_EXPIRE_TIME)
	return true
}

func DeleteSessionByUsername(username string) bool {
	bearer, ok := users.Get(username)
	if !ok {
		return false
	}
	sessions.Delete(bearer.(string))
	users.Delete(username)
	return true
}

func DeleteSessionByBearer(bearer string) bool {
	session, ok := sessions.Get(bearer)
	if !ok {
		return false
	}
	sessionMap, ok := session.(*sync.Map)
	if !ok {
		return false
	}
	username, ok := sessionMap.Load("session_username")
	if ok {
		users.Delete(username.(string))
	}
	sessions.Delete(bearer)
	return true
}
