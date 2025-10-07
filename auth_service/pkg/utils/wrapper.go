package utils

import (
	"net/http"
	"strconv"
	"time"

	httpError "github.com/helyus1412/auth-service/pkg/httpError"
	"github.com/labstack/echo/v4"
)

// Result common output
type Result struct {
	Data     interface{}
	MetaData interface{}
	Error    interface{}
	Count    int64
}

type MetaData struct {
	Page      int     `json:"page"`
	Quantity  int64   `json:"quantity"`
	TotalPage float64 `json:"totalPage"`
	TotalData int64   `json:"totalData"`
}

type ResultCount struct {
	Data     int64
	MetaData interface{}
	Error    interface{}
}

// BaseWrapperModel data structure
type BaseWrapperModel struct {
	Code          string      `json:"code"`
	Message       string      `json:"message"`
	SystemMessage string      `json:"system_message"`
	Data          interface{} `json:"data"`
	Meta          interface{} `json:"meta,omitempty"`
	Timestamp     int64       `json:"timestamp"`
}

// Response function
func Response(data interface{}, message string, code int, c echo.Context) error {
	result := BaseWrapperModel{
		Data:      data,
		Message:   message,
		Code:      getRespCode(code, ""),
		Timestamp: time.Now().Unix(),
	}

	return c.JSON(code, result)
}

// PaginationResponse function
func PaginationResponse(data interface{}, meta interface{}, message string, code int, c echo.Context) error {
	result := BaseWrapperModel{
		Data:      data,
		Meta:      meta,
		Message:   message,
		Code:      getRespCode(code, ""),
		Timestamp: time.Now().Unix(),
	}

	return c.JSON(code, result)
}

// ResponseError function
func ResponseError(err interface{}, c echo.Context) error {
	errObj := getErrorStatusCode(err)
	result := BaseWrapperModel{
		Data:          errObj.Data,
		Message:       errObj.Message,
		SystemMessage: errObj.SystemMessage,
		Code:          getRespCode(errObj.Code, errObj.CustomCode),
		Timestamp:     time.Now().Unix(),
	}
	return c.JSON(errObj.ResponseCode, result)
}

func getErrorStatusCode(err interface{}) httpError.CommonError {
	errData := httpError.CommonError{}

	switch obj := err.(type) {
	case httpError.BadRequest:
		errData.ResponseCode = http.StatusBadRequest
		errData.Code = obj.Code
		errData.Data = obj.Data
		errData.Message = obj.Message
		return errData
	case httpError.Unauthorized:
		errData.ResponseCode = http.StatusUnauthorized
		errData.Code = obj.Code
		errData.Data = obj.Data
		errData.Message = obj.Message
		return errData
	case httpError.Forbidden:
		errData.ResponseCode = http.StatusForbidden
		errData.Code = obj.Code
		errData.Data = obj.Data
		errData.Message = obj.Message
		return errData
	case httpError.Conflict:
		errData.ResponseCode = http.StatusConflict
		errData.Code = obj.Code
		errData.Data = obj.Data
		errData.Message = obj.Message
		return errData
	case httpError.NotFound:
		errData.ResponseCode = http.StatusNotFound
		errData.Code = obj.Code
		errData.Data = obj.Data
		errData.Message = obj.Message
		return errData
	case httpError.InternalServerError:
		errData.ResponseCode = http.StatusInternalServerError
		errData.Code = obj.Code
		errData.Data = obj.Data
		errData.Message = obj.Message
		errData.SystemMessage = obj.SystemMessage
		return errData
	case httpError.CustomError:
		errData.ResponseCode = obj.Code
		errData.CustomCode = obj.CustomCode
		errData.Code = obj.Code
		errData.Data = obj.Data
		errData.Message = obj.Message
		return errData
	default:
		errData.Code = http.StatusConflict
		return errData
	}
}

func getRespCode(code int, customCode string) string {
	var resCode string
	switch customCode {
	case "":
		if code == 200 {
			resCode = "RES-000"
		} else {
			resCode = "RES-" + strconv.Itoa(code)
		}
	default:
		resCode = customCode
	}

	return resCode
}
