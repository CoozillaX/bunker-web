package phoenix

import (
	"bunker-core/protocol/g79"
	"bunker-web/models"
	"bunker-web/pkg/fbtoken"
	"bunker-web/pkg/giner"
	"bunker-web/pkg/sessions"
	"bunker-web/services/user"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/patrickmn/go-cache"

	g79_utils "bunker-web/pkg/g79"
)

type LoginRequest struct {
	FBToken         string `json:"login_token"`
	UserName        string `json:"username"`
	Password        string `json:"password"`
	ServerCode      string `json:"server_code" binding:"min=1,max=40"`
	ServerPasscode  string `json:"server_passcode"`
	ClientPublicKey string `json:"client_public_key"`
}

type LoginResponse struct {
	giner.BasicResponse
	ChainInfo   string          `json:"chainInfo,omitempty"`
	IPAddress   string          `json:"ip_address,omitempty"`
	GrowthLevel int             `json:"growth_level,omitempty"`
	Dry         bool            `json:"dry,omitempty"`
	Token       string          `json:"token,omitempty"`
	ResponseTo  string          `json:"respond_to,omitempty"`
	SkinInfo    *SkinInfo       `json:"skin_info,omitempty"`
	OutfitInfo  map[string]*int `json:"outfit_info,omitempty"`
	ServerMsg   string          `json:"server_msg,omitempty"`
}

var versionCache = cache.New(24*time.Hour, time.Hour) // cache[serverCode]bedrockVersion

func requestServerInfo(
	mu models.MpayUser,
	req *LoginRequest,
) (*g79.G79User, *g79.RentalServerInfo, *gin.Error) {
	// change engine version by cache
	if value, ok := versionCache.Get(req.ServerCode); ok {
		if err := mu.UpdateGameInfoByBedrockVersion(value.(string)); err != nil {
			return nil, nil, giner.NewPublicGinError(fmt.Sprintf("无法获取游戏信息: %s", err.Error()))
		}
	}
	// g79 login
	gu, ginerr := g79_utils.HandleG79Login(mu)
	if ginerr != nil {
		return nil, nil, ginerr
	}
	// chain info
	rentalInfo, protocolErr := gu.ImpactRentalServer(req.ServerCode, req.ServerPasscode, req.ClientPublicKey)
	if protocolErr != nil {
		return nil, nil, giner.NewGinErrorFromProtocolErr(protocolErr)
	}
	// cache version
	rentalBedrockVersion := strings.TrimSuffix(rentalInfo.MCVersion, "-release")
	versionCache.SetDefault(req.ServerCode, rentalBedrockVersion)
	// check version
	if gu.GetBedrockVersion() != rentalBedrockVersion {
		// re-login and get chain with updated engine version
		return requestServerInfo(mu, req)
	}
	return gu, rentalInfo, nil
}

func (*Phoenix) Login(c *gin.Context) {
	// Get session
	bearer, _ := c.Get("bearer")
	session, _ := sessions.GetSessionByBearer(bearer.(string))
	// Parse request
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.Error(giner.NewPublicGinError("无效参数"))
		return
	}
	// Login and get usr
	usr, ginerr := user.PhoenixLogin(c.ClientIP(), req.FBToken, req.UserName, req.Password)
	if usr != nil {
		session.Store("usr", usr)
	}
	if ginerr != nil {
		c.Error(ginerr)
		return
	}
	// Setup session info
	session.Store("isPhoenix", true)
	// Handle dry Login
	if req.ServerCode == "::DRY::" && req.ServerPasscode == "::DRY::" {
		newToken, _ := fbtoken.Encrypt(usr.Username, usr.Password)
		c.JSON(http.StatusOK, &LoginResponse{
			BasicResponse: giner.BasicResponse{
				Success: true,
				Message: "ok",
			},
			Dry:   true,
			Token: newToken,
		})
		// Create log
		c.Set("log", "Phoenix登录成功(DRY)")
		return
	}
	// Auto create helper if not exist
	if usr.HelperMpayUser == nil || usr.HelperMpayUser.GetToken() == "" {
		c.Error(giner.NewPublicGinError("辅助用户不存在, 请先创建辅助用户"))
		return
	}
	// Check client public key
	if req.ClientPublicKey == "" {
		c.Error(giner.NewPublicGinError("无效参数"))
		return
	}
	// Query server info
	gu, serverInfo, ginerr := requestServerInfo(usr.HelperMpayUser, &req)
	if ginerr != nil {
		c.Error(ginerr)
		return
	}
	// Owner ID check for normal limited users
	if ginerr := user.GameLicenseCheck(usr, req.ServerCode, serverInfo.OwnerID); ginerr != nil {
		c.Error(ginerr)
		return
	}
	// Fetch growth level
	lv, _, _, ginerr := g79_utils.GetG79LauncherLevel(gu)
	if ginerr != nil {
		c.Error(ginerr)
		return
	}
	// Get helper using mod
	usingMod, ginerr := g79_utils.GetCurrentUsingMod(gu)
	if ginerr != nil {
		c.Error(ginerr)
		return
	}
	// Pack info and return
	c.JSON(http.StatusOK, &LoginResponse{
		BasicResponse: giner.BasicResponse{
			Success: true,
			Message: "ok",
		},
		ChainInfo:   serverInfo.ChainInfo,
		IPAddress:   serverInfo.IPAddress,
		GrowthLevel: lv,
		SkinInfo: &SkinInfo{
			EntityID: usingMod.SkinDownloadInfo.EntityID,
			ResUrl:   usingMod.SkinDownloadInfo.ResUrl,
			IsSlim:   usingMod.SkinData.IsSlim,
		},
		OutfitInfo: usingMod.GetConfigUUID2OutfitLevel(),
	})
	session.Store("entityID", gu.EntityID)
	session.Store("engineVersion", gu.GetEngineVersion())
	session.Store("patchVersion", gu.GetPatchVersion())
	session.Store("platform", gu.GetSystemName())
	// Create log
	c.Set("log", fmt.Sprintf(
		"Phoenix登录成功(Normal), Helper: %s, Code: %s, Version: %s",
		gu.Username,
		req.ServerCode,
		serverInfo.MCVersion,
	))
}
