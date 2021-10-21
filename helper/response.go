package helper

import "strings"

//Response is used for static shape json return
type Response struct {
	Status bool `json:"status"`
	Message string `json:"message"`
	Errors interface{} `json:"error,omitempty"`
	Data interface{} `json:"data"` //dynamic
}

type DataToken struct {
	Type string `json:"type"`
	Token string `json:"token"`
}

//EmptyObj is used when data doesn't want to be null on json
type EmptyObj struct {}

//SuccessResponse method is to inject data value to dynamic success response
func SuccessResponse(status bool, message string, data interface{}) Response {
	res := Response{
		Status:  status,
		Message: message,
		Data:    data,
	}
	return res
}

//ErrorResponse method is to inject data value to dynamic error response
func ErrorResponse(message string, err string, data interface{}) Response {
	splittedError := strings.Split(err, "\n")
	res := Response{
		Status:  false,
		Message: message,
		Errors:  splittedError,
		Data:    data,
	}

	return res
}
