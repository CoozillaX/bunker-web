package g79

import (
	"bunker-core/protocol/g79"
	"bunker-core/protocol/gameinfo"
	"bunker-web/pkg/giner"
	"context"
	"encoding/json"
	"net/http"

	"github.com/gin-gonic/gin"
)

func QueryRealmsServer(ctx context.Context, gu *g79.G79User, serverCode string) (map[string]any, *gin.Error) {
	details, protocolErr := gu.QueryRealmsServer(ctx, serverCode)
	if protocolErr != nil {
		return nil, giner.NewGinErrorFromProtocolErr(protocolErr)
	}
	return details, nil
}

func QueryRealmsServerOnlinePlayers(ctx context.Context, gu *g79.G79User, serverID string, offset, length, orderType int) ([]map[string]any, *gin.Error) {
	reqBody, _ := json.Marshal(map[string]any{
		"sid":        serverID,
		"offset":     offset,
		"length":     length,
		"order_type": orderType,
	})
	reader, protocolErr := gu.CreateHttpClient(ctx).
		SetMethod(http.MethodPost).
		SetUrl(gameinfo.G79Servers.Load().ApiGatewayUrl + "/domain-server/get-online-players").
		SetRawBody(reqBody).
		SetTokenMode(g79.TOKEN_MODE_NORMAL).
		Do()
	if protocolErr != nil {
		return nil, giner.NewGinErrorFromProtocolErr(protocolErr)
	}
	var query struct {
		Entity []map[string]any `json:"entity"`
	}
	decoder := json.NewDecoder(reader)
	decoder.UseNumber()
	if err := decoder.Decode(&query); err != nil {
		return nil, giner.NewPrivateGinError(err)
	}
	return query.Entity, nil
}

func TurnOnRealmsServer(ctx context.Context, gu *g79.G79User, serverID string) *gin.Error {
	return updateRealmsServerStatus(ctx, gu, serverID, "/domain-server/start-server")
}

func TurnOffRealmsServer(ctx context.Context, gu *g79.G79User, serverID string) *gin.Error {
	return updateRealmsServerStatus(ctx, gu, serverID, "/domain-server/stop-server")
}

func updateRealmsServerStatus(ctx context.Context, gu *g79.G79User, serverID, endpoint string) *gin.Error {
	reqBody, _ := json.Marshal(map[string]any{
		"sid": serverID,
	})
	_, protocolErr := gu.CreateHttpClient(ctx).
		SetMethod(http.MethodPost).
		SetUrl(gameinfo.G79Servers.Load().ApiGatewayUrl + endpoint).
		SetRawBody(reqBody).
		SetTokenMode(g79.TOKEN_MODE_NORMAL).
		Do()
	if protocolErr != nil {
		return giner.NewGinErrorFromProtocolErr(protocolErr)
	}
	return nil
}

type RealmsServerBlocklistPlayer struct {
	UserID     int64  `json:"uid"`         // 玩家 UID
	Nickname   string `json:"nickname"`    // 玩家昵称
	AvatarURL  string `json:"head"`        // 玩家头像地址
	CreateTime int64  `json:"create_time"` // 加入黑名单时间
}

func QueryRealmsServerBlocklist(ctx context.Context, gu *g79.G79User, serverID string) ([]RealmsServerBlocklistPlayer, *gin.Error) {
	reqBody, _ := json.Marshal(map[string]any{
		"sid": serverID,
	})
	reader, protocolErr := gu.CreateHttpClient(ctx).
		SetMethod(http.MethodPost).
		SetUrl(gameinfo.G79Servers.Load().ApiGatewayUrl + "/domain-server/get-blacklist").
		SetRawBody(reqBody).
		SetTokenMode(g79.TOKEN_MODE_NORMAL).
		Do()
	if protocolErr != nil {
		return nil, giner.NewGinErrorFromProtocolErr(protocolErr)
	}
	var query struct {
		Entities []RealmsServerBlocklistPlayer `json:"entities"`
	}
	if err := json.NewDecoder(reader).Decode(&query); err != nil {
		return nil, giner.NewPrivateGinError(err)
	}
	return query.Entities, nil
}

func AddRealmsServerBlocklistPlayer(ctx context.Context, gu *g79.G79User, serverID, uid string) *gin.Error {
	return updateRealmsServerBlocklist(ctx, gu, serverID, uid, "/domain-server/add-to-blacklist")
}

func RemoveRealmsServerBlocklistPlayer(ctx context.Context, gu *g79.G79User, serverID, uid string) *gin.Error {
	return updateRealmsServerBlocklist(ctx, gu, serverID, uid, "/domain-server/remove-from-blacklist")
}

func updateRealmsServerBlocklist(ctx context.Context, gu *g79.G79User, serverID, uid, endpoint string) *gin.Error {
	reqBody, _ := json.Marshal(map[string]any{
		"sid": serverID,
		"fid": uid,
	})
	_, protocolErr := gu.CreateHttpClient(ctx).
		SetMethod(http.MethodPost).
		SetUrl(gameinfo.G79Servers.Load().ApiGatewayUrl + endpoint).
		SetRawBody(reqBody).
		SetTokenMode(g79.TOKEN_MODE_NORMAL).
		Do()
	if protocolErr != nil {
		return giner.NewGinErrorFromProtocolErr(protocolErr)
	}
	return nil
}

type RealmsServerInviteCode struct {
	Id         string `json:"id"`         // 邀请码 ID
	ServerID   string `json:"sid"`        // 山头服 SID
	Code       string `json:"code"`       // 邀请码
	ShareText  string `json:"share_text"` // 分享文案
	Permission string `json:"permission"` // 加入后的权限
}

func QueryRealmsServerInviteCodes(ctx context.Context, gu *g79.G79User, serverID string) ([]RealmsServerInviteCode, *gin.Error) {
	reqBody, _ := json.Marshal(map[string]any{
		"sid": serverID,
	})
	reader, protocolErr := gu.CreateHttpClient(ctx).
		SetMethod(http.MethodPost).
		SetUrl(gameinfo.G79Servers.Load().ApiGatewayUrl + "/domain-server/get-invite-code-list").
		SetRawBody(reqBody).
		SetTokenMode(g79.TOKEN_MODE_NORMAL).
		Do()
	if protocolErr != nil {
		return nil, giner.NewGinErrorFromProtocolErr(protocolErr)
	}
	var query struct {
		Entities []RealmsServerInviteCode `json:"entities"`
	}
	if err := json.NewDecoder(reader).Decode(&query); err != nil {
		return nil, giner.NewPrivateGinError(err)
	}
	return query.Entities, nil
}

func GenerateRealmsServerInviteCode(ctx context.Context, gu *g79.G79User, serverID string) (string, *gin.Error) {
	reqBody, _ := json.Marshal(map[string]any{
		"sid": serverID,
	})
	reader, protocolErr := gu.CreateHttpClient(ctx).
		SetMethod(http.MethodPost).
		SetUrl(gameinfo.G79Servers.Load().ApiGatewayUrl + "/domain-server/generate-invite-code").
		SetRawBody(reqBody).
		SetTokenMode(g79.TOKEN_MODE_NORMAL).
		Do()
	if protocolErr != nil {
		return "", giner.NewGinErrorFromProtocolErr(protocolErr)
	}
	var query struct {
		Entity string `json:"entity"`
	}
	if err := json.NewDecoder(reader).Decode(&query); err != nil {
		return "", giner.NewPrivateGinError(err)
	}
	return query.Entity, nil
}

func RemoveRealmsServerInviteCode(ctx context.Context, gu *g79.G79User, serverID, code string) *gin.Error {
	return updateRealmsServerInviteCode(ctx, gu, serverID, code, nil, "/domain-server/remove-invite-code")
}

func UpdateRealmsServerInviteCode(ctx context.Context, gu *g79.G79User, serverID, code, shareText string) *gin.Error {
	return updateRealmsServerInviteCode(ctx, gu, serverID, code, &shareText, "/domain-server/update-invite-code")
}

func updateRealmsServerInviteCode(ctx context.Context, gu *g79.G79User, serverID, code string, shareText *string, endpoint string) *gin.Error {
	req := map[string]any{
		"sid":  serverID,
		"code": code,
	}
	if shareText != nil {
		req["share_text"] = *shareText
	}
	reqBody, _ := json.Marshal(req)
	_, protocolErr := gu.CreateHttpClient(ctx).
		SetMethod(http.MethodPost).
		SetUrl(gameinfo.G79Servers.Load().ApiGatewayUrl + endpoint).
		SetRawBody(reqBody).
		SetTokenMode(g79.TOKEN_MODE_NORMAL).
		Do()
	if protocolErr != nil {
		return giner.NewGinErrorFromProtocolErr(protocolErr)
	}
	return nil
}

func QueryRealmsServerMembers(ctx context.Context, gu *g79.G79User, serverID string) ([]map[string]any, *gin.Error) {
	reqBody, _ := json.Marshal(map[string]any{
		"sid": serverID,
	})
	reader, protocolErr := gu.CreateHttpClient(ctx).
		SetMethod(http.MethodPost).
		SetUrl(gameinfo.G79Servers.Load().ApiGatewayUrl + "/domain-server/get-server-members").
		SetRawBody(reqBody).
		SetTokenMode(g79.TOKEN_MODE_NORMAL).
		Do()
	if protocolErr != nil {
		return nil, giner.NewGinErrorFromProtocolErr(protocolErr)
	}
	var query struct {
		Entities []map[string]any `json:"entities"`
	}
	decoder := json.NewDecoder(reader)
	decoder.UseNumber()
	if err := decoder.Decode(&query); err != nil {
		return nil, giner.NewPrivateGinError(err)
	}
	return query.Entities, nil
}

func RemoveRealmsServerMember(ctx context.Context, gu *g79.G79User, serverID, uid string, removeAll bool) *gin.Error {
	reqBody, _ := json.Marshal(map[string]any{
		"sid":        serverID,
		"fid":        uid,
		"remove_all": removeAll,
	})
	return requestRealmsServerMemberUpdate(ctx, gu, reqBody, "/domain-server/remove-server-member")
}

func SetRealmsServerMemberPermission(ctx context.Context, gu *g79.G79User, serverID, uid, permission string) *gin.Error {
	reqBody, _ := json.Marshal(map[string]any{
		"sid":        serverID,
		"player_uid": uid,
		"permission": permission,
	})
	return requestRealmsServerMemberUpdate(ctx, gu, reqBody, "/domain-server/set-player-permission")
}

func InviteRealmsServerMember(ctx context.Context, gu *g79.G79User, serverID, uid, permission string) *gin.Error {
	reqBody, _ := json.Marshal(map[string]any{
		"sid":         serverID,
		"fid":         uid,
		"permission":  permission,
		"invite_type": "1",
	})
	return requestRealmsServerMemberUpdate(ctx, gu, reqBody, "/domain-server/invite-friend")
}

func QueryRealmsServerMemberInvites(ctx context.Context, gu *g79.G79User, serverID string) ([]map[string]any, *gin.Error) {
	reqBody, _ := json.Marshal(map[string]any{
		"sid": serverID,
	})
	reader, protocolErr := gu.CreateHttpClient(ctx).
		SetMethod(http.MethodPost).
		SetUrl(gameinfo.G79Servers.Load().ApiGatewayUrl + "/domain-server/get-owner-invitations").
		SetRawBody(reqBody).
		SetTokenMode(g79.TOKEN_MODE_NORMAL).
		Do()
	if protocolErr != nil {
		return nil, giner.NewGinErrorFromProtocolErr(protocolErr)
	}
	var query struct {
		Entities []map[string]any `json:"entities"`
	}
	decoder := json.NewDecoder(reader)
	decoder.UseNumber()
	if err := decoder.Decode(&query); err != nil {
		return nil, giner.NewPrivateGinError(err)
	}
	return query.Entities, nil
}

func UpdateRealmsServerMemberInvite(ctx context.Context, gu *g79.G79User, serverID, uid, permission string) *gin.Error {
	reqBody, _ := json.Marshal(map[string]any{
		"sid":        serverID,
		"fid":        uid,
		"permission": permission,
	})
	return requestRealmsServerMemberUpdate(ctx, gu, reqBody, "/domain-server/update-invitation")
}

func RemoveRealmsServerMemberInvite(ctx context.Context, gu *g79.G79User, serverID, uid string) *gin.Error {
	reqBody, _ := json.Marshal(map[string]any{
		"sid": serverID,
		"fid": uid,
	})
	return requestRealmsServerMemberUpdate(ctx, gu, reqBody, "/domain-server/withdraw-invitation")
}

func requestRealmsServerMemberUpdate(ctx context.Context, gu *g79.G79User, reqBody []byte, endpoint string) *gin.Error {
	_, protocolErr := gu.CreateHttpClient(ctx).
		SetMethod(http.MethodPost).
		SetUrl(gameinfo.G79Servers.Load().ApiGatewayUrl + endpoint).
		SetRawBody(reqBody).
		SetTokenMode(g79.TOKEN_MODE_NORMAL).
		Do()
	if protocolErr != nil {
		return giner.NewGinErrorFromProtocolErr(protocolErr)
	}
	return nil
}

func QueryRealmsServerApplications(ctx context.Context, gu *g79.G79User, serverID string) (map[string]any, *gin.Error) {
	reqBody, _ := json.Marshal(map[string]any{
		"sid": serverID,
	})
	return requestRealmsServerApplication(ctx, gu, reqBody, "/domain-server/owner-get-application-queue")
}

func ApproveRealmsServerApplications(ctx context.Context, gu *g79.G79User, serverID string, action int, items []map[string]any) (map[string]any, *gin.Error) {
	reqBody, _ := json.Marshal(map[string]any{
		"sid":    serverID,
		"action": action,
		"items":  items,
	})
	return requestRealmsServerApplication(ctx, gu, reqBody, "/domain-server/owner-approve-application-queue")
}

func requestRealmsServerApplication(ctx context.Context, gu *g79.G79User, reqBody []byte, endpoint string) (map[string]any, *gin.Error) {
	reader, protocolErr := gu.CreateHttpClient(ctx).
		SetMethod(http.MethodPost).
		SetUrl(gameinfo.G79Servers.Load().ApiGatewayUrl + endpoint).
		SetRawBody(reqBody).
		SetTokenMode(g79.TOKEN_MODE_NORMAL).
		Do()
	if protocolErr != nil {
		return nil, giner.NewGinErrorFromProtocolErr(protocolErr)
	}
	var query struct {
		Entity map[string]any `json:"entity"`
	}
	decoder := json.NewDecoder(reader)
	decoder.UseNumber()
	if err := decoder.Decode(&query); err != nil {
		return nil, giner.NewPrivateGinError(err)
	}
	return query.Entity, nil
}
