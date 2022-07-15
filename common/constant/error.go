// Package constant
// @author： Boice
// @createTime：2022/7/14 10:10
package constant

type ServiceError struct {
	StatusCode int    `json:"status_code"`
	Code       int    `json:"code"`
	Message    string `json:"message"`
}

func (s ServiceError) Error() string {
	return s.Message
}
