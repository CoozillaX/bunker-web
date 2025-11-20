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

type G79MailBoxItem struct {
	MailID  string `json:"mailid"`
	FromUID string `json:"fromuid"`
	Reward  []struct {
		Tp         any    `json:"tp"`
		Nb         any    `json:"nb"` // count
		Name       string `json:"name"`
		Tips       string `json:"tips"`
		ID         string `json:"id"`
		Url        string `json:"url"` // reward image url
		JumpTarget string `json:"jumpTarget"`
		JumpType   string `json:"jumpType"`
		JumpUrl    string `json:"jumpUrl"`
	} `json:"reward"`
	ReceivedReward string `json:"received_reward"`
	CreateTime     string `json:"create_time"`
	IsRead         string `json:"isread"`
	MailTP         string `json:"mailtp"`
	MType          int    `json:"mtype"`
	Title          string `json:"title"`
	Duration       int    `json:"duration"`
	Nickname       string `json:"nickname"`
}

func GetGeneralMailList(gu *g79.G79User) ([]*G79MailBoxItem, *gin.Error) {
	var result []*G79MailBoxItem
	currentOffset := "0"
	for {
		// 1. Do req
		reader, protocolErr := gu.CreateHttpClient().
			SetMethod(http.MethodPost).
			SetUrl(gameinfo.G79Servers.Load().ApiGatewayUrl + "/get-general-mail-list/").
			SetRawBody([]byte(fmt.Sprintf(`{"mailid":"%s"}`, currentOffset))).
			SetTokenMode(g79.TOKEN_MODE_NORMAL).
			Do()
		if protocolErr != nil {
			return nil, giner.NewGinErrorFromProtocolErr(protocolErr)
		}
		// 2. Parse response
		var query struct {
			Entities []*G79MailBoxItem `json:"entities"`
		}
		if err := json.NewDecoder(reader).Decode(&query); err != nil {
			return nil, giner.NewPrivateGinError(err)
		}
		// 3. Check length
		resultLength := len(query.Entities)
		if resultLength < 1 {
			return result, nil
		}
		// 4. Update result
		result = append(result, query.Entities...)
		// 5. Update offset
		currentOffset = query.Entities[resultLength-1].MailID
	}
}

func ReadMails(gu *g79.G79User, mailItemList []*G79MailBoxItem) *gin.Error {
	for _, mailItem := range mailItemList {
		if mailItem.IsRead == "1" {
			continue
		}
		_, protocolErr := gu.CreateHttpClient().
			SetMethod(http.MethodPost).
			SetUrl(gameinfo.G79Servers.Load().WebServerUrl + "/get-detail-mail/").
			SetRawBody([]byte(fmt.Sprintf(`{"mailid":"%s"}`, mailItem.MailID))).
			SetTokenMode(g79.TOKEN_MODE_NORMAL).
			Do()
		if protocolErr != nil {
			return giner.NewGinErrorFromProtocolErr(protocolErr)
		}
	}
	return nil
}

func ClaimMailRewards(gu *g79.G79User, mailItemList []*G79MailBoxItem) *gin.Error {
	idList := make([]string, 0)
	for _, mailItem := range mailItemList {
		if mailItem.ReceivedReward == "1" {
			continue
		}
		{
			needClaim := false
			for _, reward := range mailItem.Reward {
				if c, ok := reward.Nb.(float64); ok && c > 0 {
					needClaim = true
					break
				}
				if c, ok := reward.Nb.(string); ok && c != "0" {
					needClaim = true
					break
				}
			}
			if !needClaim {
				continue
			}
		}
		idList = append(idList, mailItem.MailID)
		// resp.state: 2: success, 3: claimed
	}
	if len(idList) == 0 {
		return nil
	}
	reqBody, _ := json.Marshal(map[string]any{
		"mailidlist": idList,
	})
	_, protocolErr := gu.CreateHttpClient().
		SetMethod(http.MethodPost).
		SetUrl(gameinfo.G79Servers.Load().WebServerUrl + "/get-all-mail-reward/").
		SetRawBody(reqBody).
		SetTokenMode(g79.TOKEN_MODE_NORMAL).
		Do()
	return giner.NewGinErrorFromProtocolErr(protocolErr)
}
