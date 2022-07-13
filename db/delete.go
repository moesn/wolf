package db

import (
	"github.com/kataras/iris/v12"
	"github.com/moesn/wolf/common/strs"
	"github.com/moesn/wolf/http"
	"reflect"
	"strings"
)

// 删除
func Delete(ctx iris.Context, model interface{}) *http.JsonResult {
	var ids []string
	err := ctx.ReadJSON(&ids) // 读取Json请求参数

	if err != nil || len(ids) == 0 { // 读取Json错误或为空，返回请求参数格式错误
		return http.JsonErrorMsg(err.Error())
	}

	err = DB().Delete(model, "id in (?)", ids).Error

	if err != nil { // 删除错误，返回异常错误信息
		return http.JsonErrorMsg(err.Error())
	}

	if logger!=nil{
		logMap:=GetLogColumn(model,"-")
		logMap["Id"]=strings.Join(ids,",")

		logger(ctx,logMap,"删除")
	}

	return http.JsonData(nil) // 返回成功
}

func GetLogColumn(obj interface{},column string) map[string]interface{} {
	return GetLog(obj,column)
}

func GetLogMap(obj interface{}) map[string]interface{} {
	return GetLog(obj,"")
}

func GetLog(obj interface{},column string) map[string]interface{} {
	var data = make(map[string]interface{})

	keys := reflect.TypeOf(obj)
	values := reflect.ValueOf(obj)

	if values.Kind() == reflect.Ptr {
		values = values.Elem()
	}
	if keys.Kind() == reflect.Ptr {
		keys = keys.Elem()
	}

	for i := 0; i < keys.NumField(); i++ {
		keyField := keys.Field(i)
		valueField := values.Field(i)
		jsonTag:=keyField.Tag.Get("json")

		if keyField.Name=="_Table"{
			data["_Table"] = keyField.Tag.Get("trans")
		}else if keyField.Name=="Id"{
			data["Id"] = valueField.Interface()
		}else if len(column)==0||jsonTag==column{
			trans__Tag := keyField.Tag.Get("trans__")
			if len(trans__Tag)>0{
				data[trans__Tag] = keyField.Tag.Get(strs.ToString(valueField.Interface()))
			}else{
				trans_Tag := keyField.Tag.Get("trans_")
				if len(trans_Tag)>0{
					data[trans_Tag] = valueField.Interface()
				}
			}
		}
	}

	return data
}

