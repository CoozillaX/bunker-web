package g79

import (
	"bunker-core/protocol/defines"
	"bunker-core/protocol/g79"
	"bunker-core/protocol/gameinfo"
	"bunker-web/pkg/giner"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/patrickmn/go-cache"
)

const defaultTTL = 4

type g79UserCacheItem struct {
	gu  *g79.G79User
	ttl int
}

var (
	g79UserCache *cache.Cache // cache[MpayUserUid]*g79UserCacheItem
)

func init() {
	g79UserCache = cache.New(25*time.Minute, 5*time.Minute)
	g79UserCache.OnEvicted(func(uid string, value any) {
		item := value.(*g79UserCacheItem)
		if item.ttl > 0 && item.gu.Update() == nil { // no need to logout if update failed
			item.ttl--
			g79UserCache.SetDefault(uid, item)
		} else {
			item.gu.Logout()
		}
	})
}

func HandleG79Login(mu *defines.MpayUser, engineVersion *string) (*g79.G79User, *gin.Error) {
	// check cache
	if cached, ok := g79UserCache.Get(mu.Uid); ok {
		item := cached.(*g79UserCacheItem)
		gu := item.gu
		// if version match?
		if engineVersion == nil || gu.GameInfo.EngineVersion == *engineVersion {
			// if still valid ?
			if _, _, protocolErr := gu.AccOnlineExp(); protocolErr == nil {
				item.ttl = defaultTTL // refresh ttl
				g79UserCache.SetDefault(mu.Uid, item)
				return gu, nil
			}
		}
	}
	// handle engine version
	version := gameinfo.DefaultEngineVersion
	if engineVersion != nil {
		version = *engineVersion
	}
	// g79 login
	gu, protocolErr := g79.Login(version, mu)
	if protocolErr != nil {
		return nil, giner.NewGinErrorFromProtocolErr(protocolErr)
	}
	// cache
	g79UserCache.SetDefault(mu.Uid, &g79UserCacheItem{
		gu:  gu,
		ttl: defaultTTL,
	})
	return gu, nil
}
