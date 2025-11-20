package phoenix

type Phoenix struct{}

type SkinInfo struct {
	EntityID string `json:"entity_id"`
	ResUrl   string `json:"res_url"`
	IsSlim   bool   `json:"is_slim"`
}