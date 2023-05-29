package common

const (
	successCode    = 0
	successMessage = "success"
)

type baseResponse struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}
