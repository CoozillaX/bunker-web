package giner

import (
	"bunker-core/protocol/defines"
	"errors"

	"github.com/gin-gonic/gin"
)

const (
	metaKeyPublic      = "public"
	metaKeyTranslation = "translation"
)

func NewGinError(publicErrStr string, privateErr error) *gin.Error {
	return &gin.Error{
		Err:  privateErr,
		Type: gin.ErrorTypePrivate,
		Meta: map[string]any{
			metaKeyPublic: publicErrStr,
		},
	}
}

func NewGinErrorFromProtocolErr(protocolErr *defines.ProtocolError) *gin.Error {
	if protocolErr == nil {
		return nil
	}
	newGinErr := NewPublicGinError(protocolErr.Error())
	if protocolErr.VerifyUrl != "" {
		newGinErr = SetVerifyUrl(newGinErr, protocolErr.VerifyUrl)
	}
	return newGinErr
}

func NewPublicGinError(publicErrStr string) *gin.Error {
	return &gin.Error{
		Err:  errors.New(publicErrStr),
		Type: gin.ErrorTypePublic,
		Meta: map[string]any{
			metaKeyPublic: publicErrStr,
		},
	}
}

func NewPrivateGinError(privateErr error) *gin.Error {
	if privateErr == nil {
		return nil
	}
	return &gin.Error{
		Err:  privateErr,
		Type: gin.ErrorTypePrivate,
		Meta: map[string]any{
			metaKeyPublic: "验证服务器内部错误",
		},
	}
}

func GetPrivateErrorString(ginerr *gin.Error) string {
	if ginerr == nil || ginerr.Err == nil {
		return ""
	}
	if !ginerr.IsType(gin.ErrorTypePrivate) {
		return ""
	}
	return ginerr.Err.Error()
}

func GetPublicErrorString(ginerr *gin.Error) string {
	if ginerr == nil {
		return ""
	}
	if ginerr.IsType(gin.ErrorTypePublic) {
		return ginerr.Error()
	}
	if ginerr.Meta == nil {
		return ""
	}
	metaMap, ok := ginerr.Meta.(map[string]any)
	if !ok {
		return ""
	}
	pub, ok := metaMap[metaKeyPublic]
	if !ok {
		return ""
	}
	return pub.(string)
}

func GetTranslationCode(ginerr *gin.Error) int {
	if ginerr == nil {
		return 0
	}
	if ginerr.Meta == nil {
		return 0
	}
	metaMap, ok := ginerr.Meta.(map[string]any)
	if !ok {
		return 0
	}
	trs, ok := metaMap[metaKeyTranslation]
	if !ok {
		return 0
	}
	return trs.(int)
}

func SetTranslationCode(ginerr *gin.Error, code int) *gin.Error {
	if ginerr.Meta == nil {
		ginerr.Meta = map[string]any{}
	}
	metaMap, ok := ginerr.Meta.(map[string]any)
	if !ok {
		return ginerr
	}
	metaMap[metaKeyTranslation] = code
	return ginerr
}

func GetVerifyUrl(ginerr *gin.Error) string {
	if ginerr == nil {
		return ""
	}
	if ginerr.Meta == nil {
		return ""
	}
	metaMap, ok := ginerr.Meta.(map[string]any)
	if !ok {
		return ""
	}
	verifyUrl, ok := metaMap["verifyUrl"]
	if !ok {
		return ""
	}
	return verifyUrl.(string)
}

func SetVerifyUrl(ginerr *gin.Error, verifyUrl string) *gin.Error {
	if ginerr.Meta == nil {
		ginerr.Meta = map[string]any{}
	}
	metaMap, ok := ginerr.Meta.(map[string]any)
	if !ok {
		return ginerr
	}
	metaMap["verifyUrl"] = verifyUrl
	return ginerr
}
