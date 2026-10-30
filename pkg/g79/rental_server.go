package g79

import (
	"bunker-core/protocol/g79"
	"bunker-core/protocol/gameinfo"
	"bunker-web/pkg/giner"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

const (
	RentalServerStatusStopped = iota
	RentalServerStatusRunning
	//RentalServerStatusUnintialized
	RentalServerStatusInitializing      = 7
	RentalServerStatusDockerUnavailable = 8
	RentalServerStatusUpdating          = 10
)

type RentalServerInfo struct {
	EntityID     string `json:"entity_id"`
	Name         string `json:"name"`
	OwnerID      int    `json:"owner_id"`
	Visibility   int    `json:"visibility"`
	Status       int    `json:"status"`
	IconIndex    int    `json:"icon_index"`
	Capacity     int    `json:"capacity"`
	MCVersion    string `json:"mc_version"`
	PlayerCount  int    `json:"player_count"`
	LikeNum      int    `json:"like_num"`
	ServerType   string `json:"server_type"`
	Offset       any    `json:"offset"` // unknown type
	HasPwd       string `json:"has_pwd"`
	ImageUrl     string `json:"image_url"`
	WorldID      string `json:"world_id"`
	MinLevel     string `json:"min_level"`
	PVP          bool   `json:"pvp"`
	ServerName   string `json:"server_name"`
	BriefSummary string `json:"brief_summary"`
	IPAddress    string `json:"ip_address,omitempty"`
	ChainInfo    string `json:"chainInfo,omitempty"`
}

func QueryRentalServer(ctx context.Context, gu *g79.G79User, name string) (*RentalServerInfo, *gin.Error) {
	// 1. Init server info
	serverInfo := &RentalServerInfo{}
	// 2. Get entity id
	{
		reqBody, _ := json.Marshal(map[string]any{
			"server_name": name,
			"offset":      0,
		})
		reader, protocolErr := gu.CreateHttpClient(ctx).
			SetMethod(http.MethodPost).
			SetUrl(gameinfo.G79Servers.Load().WebServerUrl + "/rental-server/query/search-by-name").
			SetRawBody(reqBody).
			SetTokenMode(g79.TOKEN_MODE_NORMAL).
			Do()
		if protocolErr != nil {
			return serverInfo, giner.NewGinErrorFromProtocolErr(protocolErr)
		}
		var query struct {
			Entities []json.RawMessage `json:"entities"`
		}
		if err := json.NewDecoder(reader).Decode(&query); err != nil {
			return serverInfo, giner.NewPrivateGinError(err)
		}
		if len(query.Entities) == 0 {
			return serverInfo, giner.NewPublicGinError(fmt.Sprintf("未查找到服务器 (%v), 请在确认服务器状态正常后重试", name))
		}
		if err := json.Unmarshal(query.Entities[0], serverInfo); err != nil {
			return serverInfo, giner.NewPrivateGinError(err)
		}
	}
	return serverInfo, nil
}

func QueryMyRentalServers(ctx context.Context, gu *g79.G79User) (result []*RentalServerInfo, ginerr *gin.Error) {
	offset := 0
	for {
		reqBody, _ := json.Marshal(map[string]any{
			"offset": offset,
			"length": 100,
		})
		reader, protocolErr := gu.CreateHttpClient(ctx).
			SetMethod(http.MethodPost).
			SetUrl(gameinfo.G79Servers.Load().WebServerUrl + "/my-rental-server/query/search-by-user").
			SetRawBody(reqBody).
			SetTokenMode(g79.TOKEN_MODE_NORMAL).
			Do()
		if protocolErr != nil {
			return nil, giner.NewGinErrorFromProtocolErr(protocolErr)
		}
		var query struct {
			Entities []*struct {
				Name      string `json:"name"`
				EntityID  string `json:"entity_id"`
				WorldID   string `json:"world_id"`
				OwnerID   string `json:"owner_id"`
				MCVersion string `json:"mc_version"`
				Status    int    `json:"status"`
			} `json:"entities"`
		}
		if err := json.NewDecoder(reader).Decode(&query); err != nil {
			return nil, giner.NewPrivateGinError(err)
		}
		if len(query.Entities) == 0 {
			return result, nil
		}
		for _, entity := range query.Entities {
			ownerId, err := strconv.Atoi(entity.OwnerID)
			if err != nil {
				return nil, giner.NewPrivateGinError(err)
			}
			result = append(result, &RentalServerInfo{
				Name:      entity.Name,
				EntityID:  entity.EntityID,
				WorldID:   entity.WorldID,
				OwnerID:   ownerId,
				MCVersion: entity.MCVersion,
				Status:    entity.Status,
			})
		}
		offset += len(query.Entities)
	}
}

func GetRentalServerIP(ctx context.Context, gu *g79.G79User, serverInfo *RentalServerInfo, passcode string) *gin.Error {
	if serverInfo == nil {
		return giner.NewPublicGinError("获取租赁服信息失败，请稍后重试")
	}
	// Never leave a previously resolved address on an unsuccessful new attempt.
	serverInfo.IPAddress = ""
	address, err := gu.QueryRentalServerIP(ctx, serverInfo.EntityID, passcode)
	if err != nil {
		return giner.NewGinErrorFromProtocolErr(err)
	}
	serverInfo.IPAddress = address
	return nil
}

// The command and every poll share one deadline. No detached polling goroutine.
func updateRentalServer(parent context.Context, gu *g79.G79User, sid string, newStatus int, timeout time.Duration) *gin.Error {
	ctx, cancel := context.WithTimeout(parent, timeout)
	defer cancel()
	timeoutError := func() *gin.Error {
		if parent.Err() != nil {
			return giner.NewPublicGinError("请求已取消")
		}
		return giner.NewPublicGinError("服务器操作尚未确认完成，请先检查当前状态，避免重复操作")
	}
	if ctx.Err() != nil {
		return timeoutError()
	}
	body, _ := json.Marshal(map[string]any{"status": newStatus, "server_id": sid})
	_, err := gu.CreateHttpClient(ctx).SetMethod(http.MethodPost).
		SetUrl(gameinfo.G79Servers.Load().ApiGatewayUrl + "/rental-server-control/update").
		SetRawBody(body).SetTokenMode(g79.TOKEN_MODE_NORMAL).Do()
	if err != nil {
		if ctx.Err() != nil {
			return timeoutError()
		}
		return giner.NewGinErrorFromProtocolErr(err)
	}
	for {
		if ctx.Err() != nil {
			return timeoutError()
		}
		body, _ := json.Marshal(map[string]any{"server_id": sid})
		reader, err := gu.CreateHttpClient(ctx).SetMethod(http.MethodPost).
			SetUrl(gameinfo.G79Servers.Load().ApiGatewayUrl + "/rental-server-control/get-status").
			SetRawBody(body).SetTokenMode(g79.TOKEN_MODE_NORMAL).Do()
		if err != nil {
			if ctx.Err() != nil {
				return timeoutError()
			}
			return giner.NewGinErrorFromProtocolErr(err)
		}
		var result struct {
			Entity *struct {
				Status *int `json:"status"`
			} `json:"entity"`
		}
		if err := json.NewDecoder(reader).Decode(&result); err != nil {
			return giner.NewPrivateGinError(err)
		}
		if result.Entity != nil && result.Entity.Status != nil {
			status := *result.Entity.Status
			if status == newStatus {
				return nil
			}
			switch status {
			case RentalServerStatusStopped, RentalServerStatusRunning, 3, 4, RentalServerStatusInitializing, RentalServerStatusDockerUnavailable, RentalServerStatusUpdating:
				// Opposite stable state may still be visible immediately after the command.
			default:
				return giner.NewPublicGinError("服务器当前状态不支持此操作，请先确认服务器状态")
			}
		}
		timer := time.NewTimer(2 * time.Second)
		select {
		case <-ctx.Done():
			timer.Stop()
			return timeoutError()
		case <-timer.C:
		}
	}
}

func TurnOnRentalServer(ctx context.Context, gu *g79.G79User, name string, timeout time.Duration) *gin.Error {
	return turnRentalServer(ctx, gu, name, RentalServerStatusRunning, timeout)
}

func TurnOffRentalServer(ctx context.Context, gu *g79.G79User, name string, timeout time.Duration) *gin.Error {
	return turnRentalServer(ctx, gu, name, RentalServerStatusStopped, timeout)
}

func turnRentalServer(parent context.Context, gu *g79.G79User, name string, status int, timeout time.Duration) *gin.Error {
	ctx, cancel := context.WithTimeout(parent, timeout)
	defer cancel()
	infos, err := QueryMyRentalServers(ctx, gu)
	if err != nil {
		return err
	}
	for _, info := range infos {
		if info.Name == name {
			return updateRentalServer(ctx, gu, info.EntityID, status, timeout)
		}
	}
	return giner.NewPublicGinError("未查找到服务器")
}

func SetRentalServerLevelLimitation(ctx context.Context, gu *g79.G79User, name string, level int) *gin.Error {
	// 0. Check level
	if level < 0 || level > 50 {
		return giner.NewPublicGinError("不在可设置范围内")
	}
	// 1. Query all rental servers
	serverInfos, ginerr := QueryMyRentalServers(ctx, gu)
	if ginerr != nil {
		return ginerr
	}
	// 2. Get entity id by name
	for _, serverInfo := range serverInfos {
		if serverInfo.Name == name {
			reqBody, _ := json.Marshal(map[string]any{
				"server_id": serverInfo.EntityID,
				"min_level": level,
			})
			_, protocolErr := gu.CreateHttpClient(ctx).
				SetMethod(http.MethodPost).
				SetUrl(gameinfo.G79Servers.Load().WebServerUrl + "/my-rental-server/update").
				SetRawBody(reqBody).
				SetTokenMode(g79.TOKEN_MODE_NORMAL).
				Do()
			return giner.NewGinErrorFromProtocolErr(protocolErr)
		}
	}
	return giner.NewPublicGinError("未查找到服务器")
}

func SetRentalServerPasswordCode(ctx context.Context, gu *g79.G79User, name string, password string) *gin.Error {
	// 0. Check password
	if len(password) != 6 {
		return giner.NewPublicGinError("密码必须为6位数字")
	}
	if _, err := strconv.Atoi(password); err != nil {
		return giner.NewPublicGinError("密码必须为6位数字")
	}
	// 1. Query all rental servers
	serverInfos, ginerr := QueryMyRentalServers(ctx, gu)
	if ginerr != nil {
		return ginerr
	}
	// 2. Get entity id by name
	for _, serverInfo := range serverInfos {
		if serverInfo.Name == name {
			reqBody, _ := json.Marshal(map[string]any{
				"server_id": serverInfo.EntityID,
				"pwd":       password,
			})
			_, protocolErr := gu.CreateHttpClient(ctx).
				SetMethod(http.MethodPost).
				SetUrl(gameinfo.G79Servers.Load().WebServerUrl + "/my-rental-server/update").
				SetRawBody(reqBody).
				SetTokenMode(g79.TOKEN_MODE_NORMAL).
				Do()
			return giner.NewGinErrorFromProtocolErr(protocolErr)
		}
	}
	return giner.NewPublicGinError("未查找到服务器")
}

func SetRentalServerVisibility(ctx context.Context, gu *g79.G79User, name string, visibility int) *gin.Error {
	// 0. Check visibility
	if visibility < 0 || visibility > 2 {
		return giner.NewPublicGinError("不在可设置范围内")
	}
	// 1. Query all rental servers
	serverInfos, ginerr := QueryMyRentalServers(ctx, gu)
	if ginerr != nil {
		return ginerr
	}
	// 2. Get entity id by name
	for _, serverInfo := range serverInfos {
		if serverInfo.Name == name {
			reqBody, _ := json.Marshal(map[string]any{
				"server_id":  serverInfo.EntityID,
				"visibility": visibility,
			})
			_, protocolErr := gu.CreateHttpClient(ctx).
				SetMethod(http.MethodPost).
				SetUrl(gameinfo.G79Servers.Load().WebServerUrl + "/my-rental-server/update").
				SetRawBody(reqBody).
				SetTokenMode(g79.TOKEN_MODE_NORMAL).
				Do()
			return giner.NewGinErrorFromProtocolErr(protocolErr)
		}
	}
	return giner.NewPublicGinError("未查找到服务器")
}

type RentalServerBackupInfo struct {
	BackupID     int    `json:"backup_id"`     // 备份槽位 ID
	BackupTS     int64  `json:"backup_ts"`     // 备份时间戳
	BackupName   string `json:"backup_name"`   // 备份名称
	BackupSize   int    `json:"backup_size"`   // 备份大小
	IsProcessing bool   `json:"is_processing"` // 是否正在处理
}

func QueryRentalServerBackupInfo(ctx context.Context, gu *g79.G79User, worldID, serverID string) ([]*RentalServerBackupInfo, *gin.Error) {
	reqBody, _ := json.Marshal(map[string]any{
		"world_id": worldID,
		"sid":      serverID,
	})
	reader, protocolErr := gu.CreateHttpClient(ctx).
		SetMethod(http.MethodPost).
		SetUrl(gameinfo.G79Servers.Load().WebServerUrl + "/rental-server-backup/query/search-by-server").
		SetRawBody(reqBody).
		SetTokenMode(g79.TOKEN_MODE_NORMAL).
		Do()
	if protocolErr != nil {
		return nil, giner.NewGinErrorFromProtocolErr(protocolErr)
	}
	// Parse backup list
	var query struct {
		Entities []*RentalServerBackupInfo `json:"entities"`
	}
	if err := json.NewDecoder(reader).Decode(&query); err != nil {
		return nil, giner.NewPrivateGinError(err)
	}
	return query.Entities, nil
}

func CreateRentalServerBackup(ctx context.Context, gu *g79.G79User, name string, backupSlot int, backupName string) *gin.Error {
	// 0. Check backup slot
	if backupSlot < 1 || backupSlot > 5 {
		return giner.NewPublicGinError("备份槽位不在可使用范围内")
	}
	// 1. Query all rental servers
	serverInfos, ginerr := QueryMyRentalServers(ctx, gu)
	if ginerr != nil {
		return ginerr
	}
	// 2. Get entity id by name
	for _, serverInfo := range serverInfos {
		if serverInfo.Name == name {
			// Create backup
			reqBody, _ := json.Marshal(map[string]any{
				"world_id":    serverInfo.WorldID,
				"backup_id":   backupSlot,
				"backup_name": backupName,
				"sid":         serverInfo.EntityID,
			})
			_, protocolErr := gu.CreateHttpClient(ctx).
				SetMethod(http.MethodPost).
				SetUrl(gameinfo.G79Servers.Load().WebServerUrl + "/rental-server-backup-upload/create-async").
				SetRawBody(reqBody).
				SetTokenMode(g79.TOKEN_MODE_NORMAL).
				Do()
			if protocolErr != nil {
				return giner.NewGinErrorFromProtocolErr(protocolErr)
			}
			return nil
		}
	}
	return giner.NewPublicGinError("未查找到服务器")
}

func QueryRentalServerBackup(ctx context.Context, gu *g79.G79User, name string) ([]*RentalServerBackupInfo, *gin.Error) {
	serverInfos, ginerr := QueryMyRentalServers(ctx, gu)
	if ginerr != nil {
		return nil, ginerr
	}
	for _, serverInfo := range serverInfos {
		if serverInfo.Name == name {
			return QueryRentalServerBackupInfo(ctx, gu, serverInfo.WorldID, serverInfo.EntityID)
		}
	}
	return nil, giner.NewPublicGinError("未查找到服务器")
}

func QueryRentalServerPlayerList(ctx context.Context, gu *g79.G79User, status *int, serverID string, orderType int, isOnline *bool, length *int, offset int) ([]map[string]any, *gin.Error) {
	reqBody, _ := json.Marshal(map[string]any{
		"status":     status,
		"server_id":  serverID,
		"order_type": orderType,
		"is_online":  isOnline,
		"length":     length,
		"offset":     offset,
	})
	reader, protocolErr := gu.CreateHttpClient(ctx).
		SetMethod(http.MethodPost).
		SetUrl(gameinfo.G79Servers.Load().WebServerUrl + "/rental-server-player/query/search-by-server").
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

func UpdateRentalServerPlayerBanStatus(ctx context.Context, gu *g79.G79User, serverID string, uid string, status int) *gin.Error {
	reqBody, _ := json.Marshal(map[string]any{
		"sid":    serverID,
		"uid":    uid,
		"status": status,
	})
	_, protocolErr := gu.CreateHttpClient(ctx).
		SetMethod(http.MethodPost).
		SetUrl(gameinfo.G79Servers.Load().WebServerUrl + "/rental-server-player/update-server-player").
		SetRawBody(reqBody).
		SetTokenMode(g79.TOKEN_MODE_NORMAL).
		Do()
	return giner.NewGinErrorFromProtocolErr(protocolErr)
}

func KickRentalServerPlayer(ctx context.Context, gu *g79.G79User, serverID string, uid string) *gin.Error {
	// Query player list
	isOnline := true
	status := 0
	playerList, ginerr := QueryRentalServerPlayerList(ctx, gu, &status, serverID, 0, &isOnline, nil, 0)
	if ginerr != nil {
		return ginerr
	}
	// Find player
	for _, player := range playerList {
		if userId, _ := player["user_id"].(string); userId == uid {
			// get entity id
			entityId, ok := player["entity_id"].(string)
			if !ok {
				return giner.NewPublicGinError("无法获取到有效的 entity id")
			}
			// kick player
			reqBody, _ := json.Marshal(map[string]any{
				"entity_id": entityId,
				"is_online": false,
			})
			_, protocolErr := gu.CreateHttpClient(ctx).
				SetMethod(http.MethodPost).
				SetUrl(gameinfo.G79Servers.Load().WebServerUrl + "/rental-server-player/update").
				SetRawBody(reqBody).
				SetTokenMode(g79.TOKEN_MODE_NORMAL).
				Do()
			return giner.NewGinErrorFromProtocolErr(protocolErr)
		}
	}
	return giner.NewPublicGinError("未查找到玩家, 请确认玩家在线")
}
