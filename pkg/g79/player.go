package g79

import (
	"bunker-core/protocol/g79"
	"bunker-core/protocol/gameinfo"
	"bunker-web/pkg/giner"
	"encoding/json"
	"net/http"

	"github.com/gin-gonic/gin"
)

type X19Player struct {
	Id                    int    `json:"id"`
	Nickname              string `json:"nickname"`
	Signature             string `json:"signature"`
	HeadImage             string `json:"headImage"`
	FrameId               string `json:"frame_id"`
	MomentId              string `json:"moment_id"`
	PublicFlag            bool   `json:"public_flag"`
	FriendRecommend       int    `json:"friend_recommend"`
	FriendApply           int    `json:"friend_apply"`
	Mark                  string `json:"mark"`
	IsFriend              bool   `json:"is_friend"`
	LogoutTimestamp       int64  `json:"tLogout"`
	PersonalPageViewCount int    `json:"personal_page_view_count"`
	PersonalPageLikeCount int    `json:"personal_page_like_count"`
	FriendCnt             int    `json:"friend_cnt"`
	MyPublicFollowCnt     int    `json:"my_public_follow_cnt"`
	PeGrowth              *struct {
		Lv       int   `json:"lv"`
		Exp      int   `json:"exp"`
		NeedExp  int   `json:"need_exp"`
		Decorate []int `json:"decorate"`
		IsVip    int   `json:"is_vip"`
	} `json:"pe_growth"`
}

func GetPlayerInfoByName(gu *g79.G79User, name string) ([]map[string]any, *gin.Error) {
	reqBody, _ := json.Marshal(map[string]any{
		"name_or_mail": name,
	})
	reader, protocolErr := gu.CreateHttpClient().
		SetMethod(http.MethodPost).
		SetUrl(gameinfo.G79Servers.Load().WebServerUrl + "/user-search-friend/").
		SetRawBody(reqBody).
		SetTokenMode(g79.TOKEN_MODE_NORMAL).
		Do()
	if protocolErr != nil {
		return nil, giner.NewGinErrorFromProtocolErr(protocolErr)
	}
	var query struct {
		Entities []map[string]any `json:"entities"`
	}
	if err := json.NewDecoder(reader).Decode(&query); err != nil {
		return nil, giner.NewPrivateGinError(err)
	}
	return query.Entities, nil
}

func GetPlayerDetailsByUid(gu *g79.G79User, uid int) (map[string]any, *gin.Error) {
	reqBody, _ := json.Marshal(map[string]any{
		"entity_id": uid,
	})
	reader, protocolErr := gu.CreateHttpClient().
		SetMethod(http.MethodPost).
		SetUrl(gameinfo.G79Servers.Load().WebServerUrl + "/user-detail/query/other").
		SetRawBody(reqBody).
		SetTokenMode(g79.TOKEN_MODE_NORMAL).
		Do()
	if protocolErr != nil {
		return nil, giner.NewGinErrorFromProtocolErr(protocolErr)
	}
	var query struct {
		Entity map[string]any `json:"entity"`
	}
	if err := json.NewDecoder(reader).Decode(&query); err != nil {
		return nil, giner.NewPrivateGinError(err)
	}
	return query.Entity, nil
}

func GetPlayerStateByUid(gu *g79.G79User, uid int) (map[string]any, *gin.Error) {
	reqBody, _ := json.Marshal(map[string]any{
		"search_id": uid,
	})
	reader, protocolErr := gu.CreateHttpClient().
		SetMethod(http.MethodPost).
		SetUrl(gameinfo.G79Servers.Load().WebServerUrl + "/user-stat/get-user-state").
		SetRawBody(reqBody).
		SetTokenMode(g79.TOKEN_MODE_NORMAL).
		Do()
	if protocolErr != nil {
		return nil, giner.NewGinErrorFromProtocolErr(protocolErr)
	}
	var query struct {
		Entity map[string]any `json:"entity"`
	}
	if err := json.NewDecoder(reader).Decode(&query); err != nil {
		return nil, giner.NewPrivateGinError(err)
	}
	return query.Entity, nil
}
