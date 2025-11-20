package giner

const (
	C_Auth_HelperNotCreated = 7
	C_Auth_InvalidToken     = 10
	C_Auth_InvalidUser      = 11
)

type BasicResponse struct {
	Success bool   `json:"success"` // 是否成功
	Message string `json:"message"` // 成功或失败原因
}

type HTTPResponse struct {
	BasicResponse
	Data        any `json:"data,omitempty"`        // 返回的数据
	Translation int `json:"translation,omitempty"` // 提供给客户端的翻译序号, 可能不存在
}

func MakeHTTPResponse(success bool) *HTTPResponse {
	resp := &HTTPResponse{
		BasicResponse: BasicResponse{
			Success: success,
		},
	}
	if success {
		resp.Message = "ok"
	} else {
		resp.Message = "fail"
	}
	return resp
}

func (r *HTTPResponse) SetMessage(message string) *HTTPResponse {
	r.Message = message
	return r
}

func (r *HTTPResponse) SetData(data any) *HTTPResponse {
	r.Data = data
	return r
}

func (r *HTTPResponse) SetTranslation(translation int) *HTTPResponse {
	r.Translation = translation
	return r
}
