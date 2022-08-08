package db

import (
	"github.com/kataras/iris/v12"
	"github.com/mitchellh/mapstructure"
	"github.com/moesn/wolf/common/jsons"
	"github.com/moesn/wolf/common/structs"
	"github.com/moesn/wolf/http"
	"github.com/moesn/wolf/http/params"
	"reflect"
	"strings"
)


func Update(ctx iris.Context, model interface{}) *http.JsonResult {
	var jsonParams map[string]interface{}
	err := ctx.ReadJSON(&jsonParams)

	if err != nil {
		return http.JsonErrorMsg(err.Error())
	}

	columns := make(map[string]interface{}, 0)

	for key, val := range jsonParams {
		if reflect.TypeOf(val) == reflect.TypeOf([]interface{}{}) {
			columns[key] = jsons.ToString(val)
		} else {
			columns[key] = val
		}
	}

	mapstructure.Decode(columns, &model)

	err = params.Verify(model)
	if err != nil {
		return http.JsonErrorMsg(err.Error())
	}

	rawData := make(map[string]interface{}, 0)
	if logger != nil {
		rawData = GetLogMap(QueryBy(columns["id"].(string), model).Data)
	}

	errDb := DB().Model(model).Where("id = ?", columns["id"]).Updates(columns).Error

	if errDb != nil {
		errMsg := errDb.Error()

		if strings.Contains(errMsg, "Duplicate entry") {
			return http.JsonErrorMsg(strings.Replace(strings.Split(errMsg, ".")[1], "'", "", 1))
		}
		return http.JsonErrorMsg(errMsg)
	}

	QueryBy(columns["id"].(string),model)

	if logger!=nil{
		logMap:=GetLogMap(model)
		for key, val := range logMap {
			if reflect.TypeOf(val)==reflect.TypeOf(structs.JSON{}){
				if(jsons.ToString(val)==jsons.ToString(rawData[key])){
					delete(logMap, key)
				}
			} else if val == rawData[key] && key != "Id" && key != "_Table" {
				delete(logMap, key)
			}
		}

		logger(ctx, logMap, "修改")
	}

	return http.JsonData(model)
}


func UpdateColumn(ctx iris.Context, model interface{}, column string, value interface{}) *http.JsonResult {
	var ids []string
	err := ctx.ReadJSON(&ids)

	if err != nil || len(ids) == 0 {
		return http.JsonErrorMsg(err.Error())
	}

	errDb := DB().Model(model).Where("id = (?)", ids).
		UpdateColumn(column, value).Error

	if errDb != nil {
		errMsg := errDb.Error()

		if strings.Contains(errMsg,"Duplicate entry"){
			return http.JsonErrorMsg(strings.Replace(strings.Split(errMsg,".")[1],"'","",1))
		}
		return http.JsonErrorMsg(errMsg)
	}

	if logger != nil {
		logMap := GetLogColumn(model, column)
		logMap["Id"] = strings.Join(ids, ",")

		logger(ctx, logMap, "修改")
	}

	return http.JsonData(nil)

}
