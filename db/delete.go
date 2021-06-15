package db

import (
	"github.com/kataras/iris/v12"
	"github.com/moesn/wolf"
	"github.com/moesn/wolf/http"
)

// 删除
func Delete(ctx iris.Context, model interface{}) *http.Response {
	var ids []string
	err := ctx.ReadJSON(&ids) // 读取Json请求参数

	if err != nil || len(ids) == 0 { // 读取Json错误或为空，返回请求参数格式错误
		return http.ResJsonError()
	}

	err = wolf.DB().Delete(model, "id in (?)", ids).Error

	if err != nil { // 删除错误，返回异常错误信息
		return http.ResError(err.Error())
	}

	return http.ResData(nil) // 返回成功
}
