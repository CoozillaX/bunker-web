package g79

import (
	"bunker-core/protocol/g79"
	"bunker-web/models"
	"bunker-web/pkg/giner"
	"strconv"
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

func HandleG79Login(mu models.MpayUser) (*g79.G79User, *gin.Error) {
	// check cache
	cacheKey := strconv.FormatUint(uint64(mu.GetID()), 10)
	if cached, ok := g79UserCache.Get(cacheKey); ok {
		item := cached.(*g79UserCacheItem)
		gu := item.gu
		// if version match?
		if mu.GetEngineVersion() == gu.GetEngineVersion() {
			// if still valid ?
			if _, _, protocolErr := gu.AccOnlineExp(); protocolErr == nil {
				item.ttl = defaultTTL // refresh ttl
				g79UserCache.SetDefault(cacheKey, item)
				return gu, nil
			}
		}
	}
	// g79 login
	gu, protocolErr := g79.Login(mu)
	if protocolErr != nil {
		return nil, giner.NewGinErrorFromProtocolErr(protocolErr)
	}
	// cache
	g79UserCache.SetDefault(cacheKey, &g79UserCacheItem{
		gu:  gu,
		ttl: defaultTTL,
	})
	return gu, nil
}
