package g79

import (
	"bunker-core/protocol/g79"
	"bunker-core/protocol/gameinfo"
	"bunker-web/pkg/giner"
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
	RentalServerStatusInitializing = 7
	RentalServerStatusUpdating     = 8
)

type RentalServerInfo struct {
	EntityID    string `json:"entity_id"`
	Name        string `json:"name"`
	OwnerID     int    `json:"owner_id"`
	Visibility  int    `json:"visibility"`
	Status      int    `json:"status"`
	IconIndex   int    `json:"icon_index"`
	Capacity    int    `json:"capacity"`
	MCVersion   string `json:"mc_version"`
	PlayerCount int    `json:"player_count"`
	LikeNum     int    `json:"like_num"`
	ServerType  string `json:"server_type"`
	Offset      any    `json:"offset"` // unknown type
	HasPwd      string `json:"has_pwd"`
	ImageUrl    string `json:"image_url"`
	WorldID     string `json:"world_id"`
	MinLevel    string `json:"min_level"`
	PVP         bool   `json:"pvp"`
	ServerName  string `json:"server_name"`
	IPAddress   string `json:"ip_address,omitempty"`
	ChainInfo   string `json:"chainInfo,omitempty"`
}

func QueryRentalServer(gu *g79.G79User, name string) (*RentalServerInfo, *gin.Error) {
	// 1. Init server info
	serverInfo := &RentalServerInfo{}
	// 2. Get entity id
	{
		reqBody, _ := json.Marshal(map[string]any{
			"server_name": name,
			"offset":      0,
		})
		reader, protocolErr := gu.CreateHttpClient().
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

func QueryMyRentalServers(gu *g79.G79User) (result []*RentalServerInfo, ginerr *gin.Error) {
	offset := 0
	for {
		reqBody, _ := json.Marshal(map[string]any{
			"offset": offset,
			"length": 100,
		})
		reader, protocolErr := gu.CreateHttpClient().
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

func GetRentalServerIP(gu *g79.G79User, serverInfo *RentalServerInfo, passcode string) *gin.Error {
	reqBody, _ := json.Marshal(map[string]any{
		"server_id": serverInfo.EntityID,
		"pwd":       passcode,
	})
	reader, protocolErr := gu.CreateHttpClient().
		SetMethod(http.MethodPost).
		SetUrl(gameinfo.G79Servers.Load().WebServerUrl + "/rental-server-world-enter/get").
		SetRawBody(reqBody).
		SetTokenMode(g79.TOKEN_MODE_NORMAL).
		Do()
	if protocolErr != nil {
		return giner.NewGinErrorFromProtocolErr(protocolErr)
	}
	var query struct {
		Entity *struct {
			MCServerHost string `json:"mcserver_host"`
			MCServerPort int    `json:"mcserver_port"`
			State        int    `json:"state"`
		} `json:"entity"`
	}
	if err := json.NewDecoder(reader).Decode(&query); err != nil {
		return giner.NewPrivateGinError(err)
	}
	if query.Entity == nil || query.Entity.State != 1 {
		return giner.NewPublicGinError("请求服务器信息失败, 请确认是否满足等级限制与密码后重试")
	}
	serverInfo.IPAddress = fmt.Sprintf("%v:%v", query.Entity.MCServerHost, query.Entity.MCServerPort)
	return nil
}

func updateRentalServerAndWaitFinished(gu *g79.G79User, sid string, newStatus int, timeout time.Duration) (ginerr *gin.Error) {
	// 1. Send command
	{
		reqBody, _ := json.Marshal(map[string]any{
			"status":    newStatus,
			"server_id": sid,
		})
		_, protocolErr := gu.CreateHttpClient().
			SetMethod(http.MethodPost).
			SetUrl(gameinfo.G79Servers.Load().WebServerUrl + "/rental-server-control/update").
			SetRawBody(reqBody).
			SetTokenMode(g79.TOKEN_MODE_NORMAL).
			Do()
		if protocolErr != nil {
			return giner.NewGinErrorFromProtocolErr(protocolErr)
		}
	}
	// 2. Wait rental server finished
	{
		waitChan := make(chan struct{})
		isFinished := false
		// Polling with 1 second interval
		go func() {
			defer close(waitChan)
			for {
				time.Sleep(time.Second)
				if isFinished {
					return
				}
				reqBody, _ := json.Marshal(map[string]any{
					"server_id": sid,
				})
				reader, protocolErr := gu.CreateHttpClient().
					SetMethod(http.MethodPost).
					SetUrl(gameinfo.G79Servers.Load().WebServerUrl + "/rental-server-control/get-status").
					SetRawBody(reqBody).
					SetTokenMode(g79.TOKEN_MODE_NORMAL).
					Do()
				if protocolErr != nil {
					ginerr = giner.NewGinErrorFromProtocolErr(protocolErr)
					return
				}
				var query struct {
					Entity *struct {
						Status int `json:"status"`
					} `json:"entity"`
				}
				if err := json.NewDecoder(reader).Decode(&query); err != nil {
					ginerr = giner.NewPrivateGinError(err)
					return
				}
				if query.Entity == nil {
					ginerr = giner.NewPublicGinError("获取服务器状态失败")
					return
				}
				if query.Entity.Status == newStatus {
					return
				}
			}
		}()
		// Wait finished or timeout
		select {
		case <-waitChan:
		case <-time.After(timeout):
			ginerr = giner.NewPublicGinError("服务器状态转换未能在指定时间内完成，请耐心等待")
		}
		isFinished = true
		return ginerr
	}

}

func TurnOnRentalServer(gu *g79.G79User, name string, timeout time.Duration) *gin.Error {
	// 1. Query all rental servers
	serverInfos, ginerr := QueryMyRentalServers(gu)
	if ginerr != nil {
		return ginerr
	}
	// 2. Get entity id by name
	for _, serverInfo := range serverInfos {
		if serverInfo.Name == name {
			return updateRentalServerAndWaitFinished(gu, serverInfo.EntityID, RentalServerStatusRunning, timeout)
		}
	}
	return giner.NewPublicGinError("未查找到服务器")

}

func TurnOffRentalServer(gu *g79.G79User, name string, timeout time.Duration) *gin.Error {
	// 1. Query all rental servers
	serverInfos, ginerr := QueryMyRentalServers(gu)
	if ginerr != nil {
		return ginerr
	}
	// 2. Get entity id by name
	for _, serverInfo := range serverInfos {
		if serverInfo.Name == name {
			return updateRentalServerAndWaitFinished(gu, serverInfo.EntityID, RentalServerStatusStopped, timeout)
		}
	}
	return giner.NewPublicGinError("未查找到服务器")
}

func SetRentalServerLevelLimitation(gu *g79.G79User, name string, level int) *gin.Error {
	// 0. Check level
	if level < 0 || level > 50 {
		return giner.NewPublicGinError("不在可设置范围内")
	}
	// 1. Query all rental servers
	serverInfos, ginerr := QueryMyRentalServers(gu)
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
			_, protocolErr := gu.CreateHttpClient().
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

func SetRentalServerPasswordCode(gu *g79.G79User, name string, password string) *gin.Error {
	// 0. Check password
	if len(password) != 6 {
		return giner.NewPublicGinError("密码必须为6位数字")
	}
	if _, err := strconv.Atoi(password); err != nil {
		return giner.NewPublicGinError("密码必须为6位数字")
	}
	// 1. Query all rental servers
	serverInfos, ginerr := QueryMyRentalServers(gu)
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
			_, protocolErr := gu.CreateHttpClient().
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

func SetRentalServerVisibility(gu *g79.G79User, name string, visibility int) *gin.Error {
	// 0. Check visibility
	if visibility < 0 || visibility > 2 {
		return giner.NewPublicGinError("不在可设置范围内")
	}
	// 1. Query all rental servers
	serverInfos, ginerr := QueryMyRentalServers(gu)
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
			_, protocolErr := gu.CreateHttpClient().
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
	BackupID     int    `json:"backup_id"`
	BackupTS     int64  `json:"backup_ts"`
	BackupName   string `json:"backup_name"`
	BackupSize   int    `json:"backup_size"`
	IsProcessing bool   `json:"is_processing"`
}

func QueryRentalServerBackupInfo(gu *g79.G79User, worldID, serverID string) ([]*RentalServerBackupInfo, *gin.Error) {
	reqBody, _ := json.Marshal(map[string]any{
		"world_id": worldID,
		"sid":      serverID,
	})
	reader, protocolErr := gu.CreateHttpClient().
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

func CreateRentalServerBackup(gu *g79.G79User, name string, backupSlot int, backupName string) *gin.Error {
	// 0. Check backup slot
	if backupSlot < 1 || backupSlot > 5 {
		return giner.NewPublicGinError("备份槽位不在可使用范围内")
	}
	// 1. Query all rental servers
	serverInfos, ginerr := QueryMyRentalServers(gu)
	if ginerr != nil {
		return ginerr
	}
	// 2. Get entity id by name
	for _, serverInfo := range serverInfos {
		if serverInfo.Name == name {
			// Check status
			if serverInfo.Status != RentalServerStatusStopped {
				return giner.NewPublicGinError("请关闭服务器后再进行备份")
			}
			// Create backup
			reqBody, _ := json.Marshal(map[string]any{
				"world_id":    serverInfo.WorldID,
				"backup_id":   backupSlot,
				"backup_name": backupName,
				"sid":         serverInfo.EntityID,
			})
			_, protocolErr := gu.CreateHttpClient().
				SetMethod(http.MethodPost).
				SetUrl(gameinfo.G79Servers.Load().WebServerUrl + "/rental-server-backup-upload/create-async").
				SetRawBody(reqBody).
				SetTokenMode(g79.TOKEN_MODE_NORMAL).
				Do()
			if protocolErr != nil {
				return giner.NewGinErrorFromProtocolErr(protocolErr)
			}
			// Polling backup status
			for {
				time.Sleep(time.Second)
				backupInfoList, ginerr := QueryRentalServerBackupInfo(gu, serverInfo.WorldID, serverInfo.EntityID)
				if ginerr != nil {
					return ginerr
				}
				// Check backup status
				isProcessing := false
				for _, backupInfo := range backupInfoList {
					if backupInfo.BackupID == backupSlot {
						isProcessing = backupInfo.IsProcessing
						break
					}
				}
				if !isProcessing {
					return nil
				}
			}

		}
	}
	return giner.NewPublicGinError("未查找到服务器")
}

func QueryRentalServerPlayerList(gu *g79.G79User, status *int, serverID string, orderType int, isOnline *bool, length *int, offset int) ([]map[string]any, *gin.Error) {
	reqBody, _ := json.Marshal(map[string]any{
		"status":     status,
		"server_id":  serverID,
		"order_type": orderType,
		"is_online":  isOnline,
		"length":     length,
		"offset":     offset,
	})
	reader, protocolErr := gu.CreateHttpClient().
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

func UpdateRentalServerPlayerBanStatus(gu *g79.G79User, serverID string, uid string, status int) *gin.Error {
	reqBody, _ := json.Marshal(map[string]any{
		"sid":    serverID,
		"uid":    uid,
		"status": status,
	})
	_, protocolErr := gu.CreateHttpClient().
		SetMethod(http.MethodPost).
		SetUrl(gameinfo.G79Servers.Load().WebServerUrl + "/rental-server-player/update-server-player").
		SetRawBody(reqBody).
		SetTokenMode(g79.TOKEN_MODE_NORMAL).
		Do()
	return giner.NewGinErrorFromProtocolErr(protocolErr)
}

func KickRentalServerPlayer(gu *g79.G79User, serverID string, uid string) *gin.Error {
	// Query player list
	isOnline := true
	status := 0
	playerList, ginerr := QueryRentalServerPlayerList(gu, &status, serverID, 0, &isOnline, nil, 0)
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
			_, protocolErr := gu.CreateHttpClient().
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
