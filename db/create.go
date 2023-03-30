package db

import (
	"github.com/kataras/iris/v12"
	"github.com/moesn/wolf/http"
	"github.com/moesn/wolf/http/params"
	"strings"
)

func Create(ctx iris.Context, model interface{}) *http.JsonResult {
	return Creater(ctx, model, nil, nil)
}

func CreatePre(ctx iris.Context, model interface{}, preProcess func()) *http.JsonResult {
	return Creater(ctx, model, preProcess, nil)
}

func CreatePost(ctx iris.Context, model interface{}, postProcess func()) *http.JsonResult {
	return Creater(ctx, model, nil, postProcess)
}

func CreatePrePost(ctx iris.Context, model interface{}, preProcess, postProcess func()) *http.JsonResult {
	return Creater(ctx, model, preProcess, postProcess)
}

func Creater(ctx iris.Context, model interface{}, preProcess, postProcess func()) *http.JsonResult {
	err := ctx.ReadJSON(&model)

	if err != nil {
		return http.JsonErrorMsg(err.Error())
	}

	err = params.Verify(model)
	if err != nil {
		return http.JsonErrorMsg(err.Error())
	}

	if preProcess != nil {
		preProcess()
	}

	errDb := DB().Create(model).Error
	if errDb != nil {
		errMsg := errDb.Error()

		if strings.Contains(errMsg, "Duplicate entry") {
			return http.JsonErrorMsg(strings.Replace(strings.Split(errMsg, ".")[1], "'", "", 1))
		}
		return http.JsonErrorMsg(errMsg)
	}

	logMap := GetLogMap(model)

	QueryById(logMap["Id"].(string), model)

	if logger != nil {
		logger(ctx, logMap, "新增")
	}

	if postProcess != nil {
		postProcess()
	}

	return http.JsonData(model)
}
