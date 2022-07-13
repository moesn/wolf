package db

import (
	"github.com/kataras/iris/v12"
	"github.com/moesn/wolf/http"
	"strings"
)

func Create(ctx iris.Context, model interface{}) *http.JsonResult {
	return Creater(ctx,model,nil,nil)
}

func CreatePre(ctx iris.Context, model interface{},preProcess func()) *http.JsonResult {
	return Creater(ctx,model,preProcess,nil)
}

func CreatePost(ctx iris.Context, model interface{},postProcess func()) *http.JsonResult {
	return Creater(ctx,model,nil,postProcess)
}

func CreatePrePost(ctx iris.Context, model interface{},preProcess func(),postProcess func()) *http.JsonResult {
	return Creater(ctx,model,preProcess,postProcess)
}

func Creater(ctx iris.Context, model interface{},preProcess func(),postProcess func()) *http.JsonResult {
	err := ctx.ReadJSON(&model) // 读取Json请求参数

	if err != nil { // 读取Json错误，返回请求参数格式错误
		return http.JsonErrorMsg(err.Error())
	}

	err = http.Verify(model) // 校验参数合法性
	if err != nil {          // 参数有误，返回参数错误信息
		return http.JsonErrorMsg(err.Error())
	}

	if preProcess!=nil{
		preProcess()
	}

	errDb := DB().Create(model).Error // 增加数据
	if errDb != nil {                    // 增加错误，返回异常错误信息
		errMsg:=errDb.Error()

		if strings.Contains(errMsg,"Duplicate entry"){
			return http.JsonErrorMsg(strings.Replace(strings.Split(errMsg,".")[1],"'","",1))
		}
		return http.JsonErrorMsg(errMsg)
	}

	logMap:=GetLogMap(model)

	QueryBy(logMap["Id"].(string),model)

	if logger!=nil{
		logger(ctx,logMap,"新增")
	}

	if postProcess!=nil{
		postProcess()
	}

	return http.JsonData(model) // 返回成功
}
