package web

import (
	"github.com/moesn/wolf/common/structs"
	"github.com/moesn/wolf/sqls"
)

type JsonResult struct {
	ErrorCode int         `json:"code"`
	Message   string      `json:"msg"`
	Data      interface{} `json:"data"`
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
		ErrorCode: 0,
		Data:      data,
	}
}

func JsonItemList(data []interface{}) *JsonResult {
	return &JsonResult{
		ErrorCode: 0,
		Data:      data,
	}
}

func JsonPageData(results interface{}, page *sqls.Paging) *JsonResult {
	return JsonData(&PageResult{
		Results: results,
		Page:    page,
	})
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
		ErrorCode: 0,
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

func JsonErrorMsg(message string) *JsonResult {
	return &JsonResult{
		ErrorCode: 0,
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
	return &RspBuilder{Data: structs.StructToMap(obj, excludes...)}
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
