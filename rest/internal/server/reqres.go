package server

type Response struct {
	StatusCode int    `json:"status_code"`
	Ok         string `json:"ok"`
	Error      string `json:"error"`
}
