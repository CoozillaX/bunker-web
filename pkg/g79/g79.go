package g79

import (
	"bunker-core/protocol/g79"
	"bunker-web/models"
	"bunker-web/pkg/giner"
	"bunker-web/pkg/sessions"
	"context"
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
		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()
		item := value.(*g79UserCacheItem)
		if item.ttl > 0 && item.gu.Update(ctx) == nil { // no need to logout if update failed
			item.ttl--
			g79UserCache.SetDefault(uid, item)
		} else {
			item.gu.Logout(ctx)
		}
	})
}

func HandleG79Login(ctx context.Context, mu models.MpayUser) (*g79.G79User, *gin.Error) {
	if ctx.Err() != nil {
		return nil, giner.NewPublicGinError(ctx.Err().Error())
	}
	// need to initialise first for cache check
	if protocolErr := mu.Initialise(); protocolErr != nil {
		return nil, giner.NewGinErrorFromProtocolErr(protocolErr)
	}
	// check cache
	if cached, ok := g79UserCache.Get(mu.GetUid() + mu.GetEngineVersion()); ok {
		item := cached.(*g79UserCacheItem)
		gu := item.gu
		// if version match?
		if mu.GetEngineVersion() == gu.GetEngineVersion() && mu.GetPatchVersion() == gu.GetPatchVersion() {
			// if still valid?
			if _, _, protocolErr := gu.AccOnlineExp(ctx); protocolErr == nil {
				item.ttl = defaultTTL // refresh ttl
				g79UserCache.SetDefault(mu.GetUid()+mu.GetEngineVersion(), item)
				return gu, nil
			}
		}
	}
	// g79 login
	gu, protocolErr := g79.Login(ctx, mu)
	if protocolErr != nil {
		return nil, giner.NewGinErrorFromProtocolErr(protocolErr)
	}
	// cache
	g79UserCache.SetDefault(mu.GetUid()+mu.GetEngineVersion(), &g79UserCacheItem{
		gu:  gu,
		ttl: defaultTTL,
	})
	return gu, nil
}

func GetHelperG79User(c *gin.Context) (*g79.G79User, *gin.Error) {
	usr, ok := getUserFromContext(c)
	if !ok || usr.HelperMpayUser == nil || usr.HelperMpayUser.GetToken() == "" {
		return nil, giner.NewPublicGinError("未创建辅助用户")
	}
	return HandleG79Login(c.Request.Context(), usr.HelperMpayUser)
}

func GetOwnerG79User(c *gin.Context) (*g79.G79User, *gin.Error) {
	usr, ok := getUserFromContext(c)
	if !ok || usr.OwnerMpayUser == nil || usr.OwnerMpayUser.GetToken() == "" {
		return nil, giner.NewPublicGinError("请先绑定游戏账号")
	}
	return HandleG79Login(c.Request.Context(), usr.OwnerMpayUser)
}

func getUserFromContext(c *gin.Context) (*models.User, bool) {
	if u, exists := c.Get("usr"); exists {
		usr, ok := u.(*models.User)
		return usr, ok
	}
	bearer, exists := c.Get("bearer")
	if !exists {
		return nil, false
	}
	bearerString, ok := bearer.(string)
	if !ok {
		return nil, false
	}
	session, ok := sessions.GetSessionByBearer(bearerString)
	if !ok {
		return nil, false
	}
	u, exists := session.Load("usr")
	if !exists {
		return nil, false
	}
	usr, ok := u.(*models.User)
	return usr, ok
}
