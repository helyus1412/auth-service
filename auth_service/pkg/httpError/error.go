package error

import "net/http"

// CommonError struct
type CommonError struct {
	Code             int         `json:"code"`
	CustomCode       string      `json:"custom_code"`
	ResponseCode     int         `json:"responseCode,omitempty"`
	Message          string      `json:"message"`
	SystemMessage    string      `json:"system_message"`
	Data             interface{} `json:"data"`
	ValidationErrors interface{} `json:"validation_errors,omitempty"`
}

// BadRequest struct
type BadRequest struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

// NewBadRequest
func NewBadRequest(message string) BadRequest {
	errObj := BadRequest{}
	errObj.Message = "Bad Request"
	errObj.Code = http.StatusBadRequest
	if message != "" {
		errObj.Message = message
	}

	return errObj
}

// NotFound struct
type NotFound struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

func NewNotFound(message string) NotFound {
	errObj := NotFound{}
	errObj.Message = "NotFound"
	if message != "" {
		errObj.Message = message
	}
	errObj.Code = http.StatusNotFound

	return errObj
}

// Unauthorized struct
type Unauthorized struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

func NewUnauthorized(message string) Unauthorized {
	errObj := Unauthorized{}
	errObj.Message = "Unauthorized"
	if message != "" {
		errObj.Message = message
	}
	errObj.Code = http.StatusUnauthorized

	return errObj
}

// Forbidden struct
type Forbidden struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

func NewForbidden(message string) Forbidden {
	errObj := Forbidden{}
	errObj.Message = "Forbidden"
	if message != "" {
		errObj.Message = message
	}
	errObj.Code = http.StatusForbidden

	return errObj
}

// Conflict struct
type Conflict struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

func NewConflict(message string) Conflict {
	errObj := Conflict{}
	errObj.Message = "Conflict"
	if message != "" {
		errObj.Message = message
	}
	errObj.Code = http.StatusConflict

	return errObj
}

// InternalServerError struct
type InternalServerError struct {
	Code          int         `json:"code"`
	Message       string      `json:"message"`
	SystemMessage string      `json:"system_message"`
	Data          interface{} `json:"data"`
}

func NewInternalServerError(system_msg string) InternalServerError {
	errObj := InternalServerError{}
	errObj.Message = "Internal Server Error"
	if system_msg != "" {
		errObj.SystemMessage = system_msg
	}
	errObj.Code = http.StatusInternalServerError

	return errObj
}

type CustomError struct {
	Code       int         `json:"code"`
	CustomCode string      `json:"custom_code"`
	Message    string      `json:"message"`
	Data       interface{} `json:"data"`
}

func NewCustomError(code int, customCode, message string, data ...interface{}) CustomError {
	errObj := CustomError{}
	errObj.Message = "Error"
	if message != "" {
		errObj.Message = message
	}
	errObj.CustomCode = customCode
	errObj.Code = code

	if len(data) > 0 {
		errObj.Data = data[0]
	}

	return errObj
}
