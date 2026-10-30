package g79

import (
	"bunker-core/protocol/g79"
	"bunker-core/protocol/gameinfo"
	"bunker-web/pkg/giner"
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

func GetG79LauncherLevel(ctx context.Context, gu *g79.G79User) (level, exp, needExp int, ginerr *gin.Error) {
	reader, protocolErr := gu.CreateHttpClient(ctx).
		SetMethod(http.MethodPost).
		SetUrl(gameinfo.G79Servers.Load().ApiGatewayUrl + "/pe-get-grow-lv-exp").
		SetRawBody([]byte("{}")).
		SetTokenMode(g79.TOKEN_MODE_NORMAL).
		Do()
	if protocolErr != nil {
		return 0, 0, 0, ginerr
	}
	var query struct {
		Entity struct {
			Level   int `json:"lv"`
			Exp     int `json:"exp"`
			NeedExp int `json:"need_exp"`
		} `json:"entity"`
	}
	if err := json.NewDecoder(reader).Decode(&query); err != nil {
		return 0, 0, 0, giner.NewPrivateGinError(err)
	}
	return query.Entity.Level, query.Entity.Exp, query.Entity.NeedExp, nil
}

type RandomName struct {
	Name string `json:"name"` // 随机昵称
	Type int    `json:"type"` // 昵称类型, 0 表示账号原有昵称
}

func GetRandomNames(ctx context.Context, gu *g79.G79User) ([]RandomName, *gin.Error) {
	reader, protocolErr := gu.CreateHttpClient(ctx).
		SetMethod(http.MethodPost).
		SetUrl(gameinfo.G79Servers.Load().WebServerUrl + "/nickname-init/get-random-name").
		SetRawBody([]byte("{}")).
		SetTokenMode(g79.TOKEN_MODE_NORMAL).
		Do()
	if protocolErr != nil {
		return nil, giner.NewGinErrorFromProtocolErr(protocolErr)
	}
	var query struct {
		Entity struct {
			Names []RandomName `json:"names"`
		} `json:"entity"`
	}
	if err := json.NewDecoder(reader).Decode(&query); err != nil {
		return nil, giner.NewPrivateGinError(err)
	}
	return query.Entity.Names, nil
}

func ChangeUserName(ctx context.Context, gu *g79.G79User, name string) *gin.Error {
	reqBody, _ := json.Marshal(map[string]any{
		"name": name,
	})
	_, protocolErr := gu.CreateHttpClient(ctx).
		SetMethod(http.MethodPost).
		SetUrl(gameinfo.G79Servers.Load().WebServerUrl + "/pe-nickname-setting/update").
		SetRawBody(reqBody).
		SetTokenMode(g79.TOKEN_MODE_NORMAL).
		Do()
	if protocolErr != nil {
		return giner.NewGinErrorFromProtocolErr(protocolErr)
	}
	gu.Username = name
	return nil
}

func UseGiftCode(ctx context.Context, gu *g79.G79User, code string) *gin.Error {
	_, protocolErr := gu.CreateHttpClient(ctx).
		SetMethod(http.MethodPost).
		SetUrl(gameinfo.G79Servers.Load().WebServerUrl + "/gift-code/").
		SetRawBody(fmt.Appendf(nil, `{"code":"%s"}`, code)).
		SetTokenMode(g79.TOKEN_MODE_NORMAL).
		Do()
	return giner.NewGinErrorFromProtocolErr(protocolErr)
}
