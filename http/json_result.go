package http

import (
	"github.com/moesn/wolf/common/structs"
)

const OkCode = 0

type JsonResult struct {
	ErrorCode int         `json:"code"`
	Message   string      `json:"msg"`
	Error     string      `json:"error"`
	Data      interface{} `json:"data"`
	Type      string      `json:"type"`
	Count     int64       `json:"total"`
}

func NewJsonResult() *JsonResult {
	return &JsonResult{
		ErrorCode: OkCode,
		Message: "",
		Error: "",
		Data: "",
		Type: "",
		Count: 0,
	}
}

func Json(code int, message string, data interface{}) *JsonResult {
	return &JsonResult{
		ErrorCode: code,
		Message:   message,
		Data:      data,
	}
}

func JsonData(data interface{}) *JsonResult {
	return &JsonResult{
		ErrorCode: OkCode,
		Data:      data,
	}
}

func JsonItemList(data []interface{}) *JsonResult {
	return &JsonResult{
		ErrorCode: OkCode,
		Data:      data,
	}
}

func JsonPageData(data interface{}, count int64) *JsonResult {
	return &JsonResult{
		ErrorCode: OkCode,
		Data:      data,
		Count:     count,
	}
}

func JsonCursorData(results interface{}, cursor string, hasMore bool) *JsonResult {
	return JsonData(&CursorResult{
		Results: results,
		Cursor:  cursor,
		HasMore: hasMore,
	})
}

func JsonSuccess() *JsonResult {
	return &JsonResult{
		ErrorCode: OkCode,
		Data:      nil,
	}
}

func JsonError(err *CodeError) *JsonResult {
	return &JsonResult{
		ErrorCode: err.Code,
		Message:   err.Message,
		Data:      err.Data,
	}
}

func JsonSuccessMsg(message string) *JsonResult {
	return &JsonResult{
		ErrorCode: OkCode,
		Message:   message,
		Data:      nil,
	}
}

func JsonErrorMsg(message string) *JsonResult {
	return &JsonResult{
		ErrorCode: 1,
		Message:   message,
		Data:      nil,
	}
}

func JsonErrorCode(code int, message string) *JsonResult {
	return &JsonResult{
		ErrorCode: code,
		Message:   message,
		Data:      nil,
	}
}

func JsonErrorData(code int, message string, data interface{}) *JsonResult {
	return &JsonResult{
		ErrorCode: code,
		Message:   message,
		Data:      data,
	}
}

type RspBuilder struct {
	Data map[string]interface{}
}

func NewEmptyRspBuilder() *RspBuilder {
	return &RspBuilder{Data: make(map[string]interface{})}
}

func NewRspBuilder(obj interface{}) *RspBuilder {
	return NewRspBuilderExcludes(obj)
}

func NewRspBuilderExcludes(obj interface{}, excludes ...string) *RspBuilder {
	return &RspBuilder{Data: structs.StructToMap(obj, "json",excludes...)}
}

func (builder *RspBuilder) Put(key string, value interface{}) *RspBuilder {
	builder.Data[key] = value
	return builder
}

func (builder *RspBuilder) Build() map[string]interface{} {
	return builder.Data
}

func (builder *RspBuilder) JsonResult() *JsonResult {
	return JsonData(builder.Data)
}
