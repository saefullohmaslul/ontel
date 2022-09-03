package pkg

import (
	"errors"
	"fmt"
	"net/http"
	"reflect"

	"github.com/labstack/echo/v4"
)

const (
	ERROR_RESPONSE   = "error"
	SUCCESS_RESPONSE = "success"
)

type (
	Response struct {
		Code    int         `json:"code"`
		Message string      `json:"message"`
		Data    interface{} `json:"data"`
		Error   interface{} `json:"error"`
	}

	Error struct {
		Data string `json:"data"`
	}

	HTTPError struct {
		Code     int         `json:"-"`
		Message  interface{} `json:"message"`
		Err      Error       `json:"-"`
		Internal error       `json:"-"`
	}
)

func NewHTTPError(code int, message interface{}) *HTTPError {
	err := Error{
		Data: fmt.Sprint(message),
	}

	return &HTTPError{
		Code:     code,
		Message:  message,
		Err:      err,
		Internal: errors.New(err.Data),
	}
}

func (e *HTTPError) Error() string {
	return e.Internal.Error()
}

func NewResponse(data interface{}, err interface{}) Response {
	var (
		code                    = http.StatusOK
		defaultErr  interface{} = struct{}{}
		defaultData interface{}
	)

	rt := reflect.TypeOf(data)

	if rt != nil {
		switch expression := rt.Kind(); expression {
		case reflect.Struct:
			defaultData = make(map[string]interface{}, 0)
		case reflect.Ptr:
			dataInterface := reflect.New(rt.Elem()).Interface()
			dataKind := reflect.TypeOf(dataInterface).Kind()
			switch dataKind {
			case reflect.Struct:
				defaultData = make(map[string]interface{}, 0)
			default:
				defaultData = make([]interface{}, 0)
			}
		default:
			defaultData = make([]interface{}, 0)
		}
	} else {
		defaultData = make(map[string]interface{}, 0)
	}

	if err != nil {
		if errData, ok := err.(*HTTPError); ok {
			code = errData.Code
			defaultErr = errData.Err
		} else {
			if errData, ok := err.(*echo.HTTPError); ok {
				code = errData.Code
				defaultErr = Error{
					Data: fmt.Sprint(errData.Message),
				}
			} else if errData, ok := err.(error); ok {
				code = http.StatusServiceUnavailable
				defaultErr = Error{
					Data: fmt.Sprint(errData.Error()),
				}
			} else {
				code = http.StatusServiceUnavailable
				defaultErr = Error{
					Data: fmt.Sprint(err),
				}
			}
		}
	}

	if code >= 300 {
		return Response{
			Code:    code,
			Message: ERROR_RESPONSE,
			Data:    defaultData,
			Error:   defaultErr,
		}
	} else {
		return Response{
			Code:    code,
			Message: SUCCESS_RESPONSE,
			Data:    data,
			Error:   defaultErr,
		}
	}
}
