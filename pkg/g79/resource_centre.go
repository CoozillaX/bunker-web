package g79

import (
	"bunker-core/protocol/g79"
	"bunker-core/protocol/gameinfo"
	"bunker-web/pkg/giner"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

type DownloadInfo struct {
	EntityID string `json:"entity_id"`
	ResUrl   string `json:"res_url"`
}

func GetDownloadInfoByItemID(gu *g79.G79User, id string) (*DownloadInfo, *gin.Error) {
	// 1. Do req
	reader, protocolErr := gu.CreateHttpClient().
		SetMethod(http.MethodPost).
		SetUrl(gameinfo.G79Servers.Load().ApiGatewayUrl + "/pe-download-item/get-download-info").
		SetRawBody(fmt.Appendf(nil, `{"item_id":"%s"}`, id)).
		SetTokenMode(g79.TOKEN_MODE_NORMAL).
		Do()
	if protocolErr != nil {
		return nil, giner.NewGinErrorFromProtocolErr(protocolErr)
	}
	// 2. Parse response
	var query struct {
		Entity DownloadInfo `json:"entity"`
	}
	if err := json.NewDecoder(reader).Decode(&query); err != nil {
		return nil, giner.NewPrivateGinError(err)
	}
	return &query.Entity, nil
}
