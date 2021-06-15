package db

import (
	"github.com/kataras/iris/v12"
	"github.com/mitchellh/mapstructure"
	"todo/wolf"
	"todo/wolf/http"
)

// 修改多个字段
func Update(ctx iris.Context, model interface{}) *http.Response {
	var columns map[string]interface{}
	err := ctx.ReadJSON(&columns) // 读取Json请求参数

	if err != nil { // 读取Json错误，返回请求参数格式错误
		return http.ResJsonError()
	}

	mapstructure.Decode(columns, &model) // 将Map的值映射进Struct

	err = wolf.Verify(model) // 校验参数合法性
	if err != nil {          // 参数有误，返回参数错误信息
		return http.ResParamError(err.Error())
	}

	err = wolf.DB().Model(model).Where("id = ?", columns["id"]).Updates(columns).Error // 修改数据

	if err != nil { // 修改错误，返回异常错误信息
		return http.ResError(err.Error())
	}

	return http.ResData(nil) // 返回成功
}

// 修改单个字段
func UpdateColumn(ctx iris.Context, model interface{}, column string, value interface{}) *http.Response {
	var ids []string
	err := ctx.ReadJSON(&ids) // 读取Json请求参数

	if err != nil || len(ids) == 0 { // 读取Json错误或为空，返回请求参数格式错误
		return http.ResJsonError()
	}

	err = wolf.DB().Model(model).Where("id = (?)", ids).
		UpdateColumn(column, value).Error

	if err != nil { // 删除错误，返回异常错误信息
		return http.ResError(err.Error())
	}

	return http.ResData(nil) // 返回成功

}