package db

import (
	"github.com/kataras/iris/v12"
	"github.com/mitchellh/mapstructure"
	"github.com/moesn/wolf/common/jsons"
	"github.com/moesn/wolf/common/structs"
	"github.com/moesn/wolf/http"
	"reflect"
	"strings"
)

// 修改多个字段
func Update(ctx iris.Context, model interface{}) *http.JsonResult {
	var params map[string]interface{}
	err := ctx.ReadJSON(&params) // 读取Json请求参数

	if err != nil { // 读取Json错误，返回请求参数格式错误
		return http.JsonErrorMsg(err.Error())
	}

	columns:= make(map[string]interface{},0)

	for key, val := range params {
		if reflect.TypeOf(val)== reflect.TypeOf([]interface{}{}){
			columns[key]=jsons.ToJsonStr(val)
		}else{
			columns[key]=val
		}
	}

	mapstructure.Decode(columns, &model) // 将Map的值映射进Struct

	err = http.Verify(model) // 校验参数合法性
	if err != nil {          // 参数有误，返回参数错误信息
		return http.JsonErrorMsg(err.Error())
	}

	rawData := structs.StructToMap(QueryBy(columns["id"].(string), model).Data, "trans")
	errDb := DB().Model(model).Where("id = ?", columns["id"]).Updates(columns).Error // 修改数据

	if errDb != nil {
		errMsg := errDb.Error()

		if strings.Contains(errMsg, "Duplicate entry") {
			return http.JsonErrorMsg(strings.Replace(strings.Split(errMsg, ".")[1], "'", "", 1))
		}
		return http.JsonErrorMsg(errMsg)
	}

	QueryBy(columns["id"].(string), model)

	logMap := structs.StructToMap(model, "trans")
	for key, val := range logMap {
		if reflect.TypeOf(val) == reflect.TypeOf(structs.JSON{}) {
			if jsons.ToJsonStr(val) == jsons.ToJsonStr(rawData[key]) {
				delete(logMap, key)
			}
		} else if val == rawData[key] && key != "ID" && key != "_Table" {
			delete(logMap, key)
		}
	}

	if logger != nil {
		logger(ctx, logMap, "修改")
	}

	return http.JsonData(model) // 返回成功
}

// 修改单个字段
func UpdateColumn(ctx iris.Context, model interface{}, column string, value interface{}) *http.JsonResult {
	var ids []string
	err := ctx.ReadJSON(&ids) // 读取Json请求参数

	if err != nil || len(ids) == 0 { // 读取Json错误或为空，返回请求参数格式错误
		return http.JsonErrorMsg(err.Error())
	}

	errDb := DB().Model(model).Where("id = (?)", ids).
		UpdateColumn(column, value).Error

	if errDb != nil {
		errMsg := errDb.Error()

		if strings.Contains(errMsg, "Duplicate entry") {
			return http.JsonErrorMsg(strings.Replace(strings.Split(errMsg, ".")[1], "'", "", 1))
		}
		return http.JsonErrorMsg(errMsg)
	}

	modelMap := structs.StructToMap(model, "trans")
	logMap := make(map[string]interface{}, 0)

	logMap["_Table"] = modelMap["_Table"]
	logMap[column] = value
	logMap["ID"] = strings.Join(ids, ",")

	if logger != nil {
		logger(ctx, logMap, "修改")
	}

	return http.JsonData(nil) // 返回成功

}
