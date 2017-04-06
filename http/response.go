package http

type response struct {
	Success bool `json:"success"`
	Message string `json:"msg,omitempty"`
	Data interface{} `json:"data,omitempty"`
	Code int `json:"-"`
	MessageCode MessageCode `json:"code,omitempty"`
}