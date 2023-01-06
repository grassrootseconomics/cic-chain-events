package api

const (
	INTERNAL_ERROR = "ERR_INTERNAL"
)

type okResp struct {
	Ok   bool        `json:"ok"`
	Data interface{} `json:"data"`
}

type errResp struct {
	Ok    bool   `json:"ok"`
	Error string `json:"error"`
}
