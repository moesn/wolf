package db

import (
	"github.com/kataras/iris/v12"
	"github.com/moesn/wolf/http"
	"strings"
)

// 增加
func Create(ctx iris.Context, model interface{}) *http.JsonResult {
	err := ctx.ReadJSON(&model) // 读取Json请求参数

	if err != nil { // 读取Json错误，返回请求参数格式错误
		return http.JsonErrorMsg(err.Error())
	}

	//err = http.Verify(model) // 校验参数合法性
	//if err != nil {          // 参数有误，返回参数错误信息
	//	return http.JsonErrorMsg(err.Error())
	//}

	errDb := DB().Create(model).Error // 增加数据
	if errDb != nil {                    // 增加错误，返回异常错误信息
		errMsg:=errDb.Error()

		if strings.Contains(errMsg,"Duplicate entry"){
			return http.JsonErrorMsg(strings.Replace(strings.Split(errMsg,".")[1],"'","",1))
		}
		return http.JsonErrorMsg(errMsg)
	}

	return http.JsonData(model) // 返回成功
}

// 增加
func Insert(model interface{}) *http.JsonResult {
	//err = http.Verify(model) // 校验参数合法性
	//if err != nil {          // 参数有误，返回参数错误信息
	//	return http.JsonErrorMsg(err.Error())
	//}

	errDb := DB().Create(model).Error // 增加数据
	if errDb != nil {                    // 增加错误，返回异常错误信息
		errMsg:=errDb.Error()

		if strings.Contains(errMsg,"Duplicate entry"){
			return http.JsonErrorMsg(strings.Replace(strings.Split(errMsg,".")[1],"'","",1))
		}
		return http.JsonErrorMsg(errMsg)
	}

	return http.JsonData(model) // 返回成功
}
