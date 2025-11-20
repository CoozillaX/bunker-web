package g79

import (
	"bunker-core/protocol/g79"
	"bunker-core/protocol/gameinfo"
	"bunker-web/pkg/giner"
	"encoding/json"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

type SkinType struct {
	Type string `json:"type"`
}

type SkinData struct {
	IsSlim     bool   `json:"is_slim"`
	ItemID     string `json:"item_id"`
	SecondType int    `json:"second_type"`
}

type ScreenConfig struct {
	ItemID        string `json:"item_id"`
	OutfitLevel   *int   `json:"outfit_level,omitempty"`
	BehaviourUUID string `json:"behaviour_uuid"`
	EffectMtypeid int    `json:"effect_mtypeid"`
	EffectStypeid int    `json:"effect_stypeid"`
}

type UsingMod struct {
	SkinType         SkinType                 `json:"skin_type"`
	SkinData         SkinData                 `json:"skin_data"`
	ScreenConfig     map[string]*ScreenConfig `json:"screen_config"`
	SkinDownloadInfo *DownloadInfo
}

func GetCurrentUsingMod(gu *g79.G79User) (*UsingMod, *gin.Error) {
	// 1. Do req
	reader, protocolErr := gu.CreateHttpClient().
		SetMethod(http.MethodPost).
		SetUrl(gameinfo.G79Servers.Load().WebServerUrl + "/pe-get-user-setting-list").
		SetRawBody([]byte(`{"settings":["skin_type","skin_data","persona_data","screen_config","outfit_type","personal_open","personal_ad_open","personal_tags"]}`)).
		SetTokenMode(g79.TOKEN_MODE_NORMAL).
		Do()
	if protocolErr != nil {
		return nil, giner.NewGinErrorFromProtocolErr(protocolErr)
	}
	// 2. Parse response
	var query struct {
		UsingMod UsingMod `json:"entity"`
	}
	if err := json.NewDecoder(reader).Decode(&query); err != nil {
		return nil, giner.NewPrivateGinError(err)
	}
	// 3. Get skin download info
	if !strings.HasPrefix(query.UsingMod.SkinData.ItemID, "-") {
		var ginerr *gin.Error
		query.UsingMod.SkinDownloadInfo, ginerr = GetDownloadInfoByItemID(gu, query.UsingMod.SkinData.ItemID)
		if ginerr != nil {
			return nil, ginerr
		}
	} else {
		query.UsingMod.SkinDownloadInfo = &DownloadInfo{
			EntityID: query.UsingMod.SkinData.ItemID,
			ResUrl:   "",
		}
	}
	return &query.UsingMod, nil
}

func (u *UsingMod) GetConfigUUID2OutfitLevel() (ret map[string]*int) {
	ret = make(map[string]*int)
	for _, v := range u.ScreenConfig {
		if v.OutfitLevel == nil {
			ret[v.BehaviourUUID] = nil
			continue
		}
		var gameOutfitLevel int
		switch *v.OutfitLevel {
		case 0:
			gameOutfitLevel = 2
		case 1:
			gameOutfitLevel = 1
		case 2:
			gameOutfitLevel = 0
		}
		ret[v.BehaviourUUID] = &gameOutfitLevel
	}
	return
}
