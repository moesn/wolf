package db

import (
	"github.com/kataras/iris/v12"
	"github.com/moesn/wolf/sqls"
	"github.com/moesn/wolf/web"
)

// 删除
func Delete(ctx iris.Context, model interface{}) *web.JsonResult {
	var ids []string
	err := ctx.ReadJSON(&ids) // 读取Json请求参数

	if err != nil || len(ids) == 0 { // 读取Json错误或为空，返回请求参数格式错误
		return web.JsonErrorMsg(err.Error())
	}

	err = sqls.DB().Delete(model, "id in (?)", ids).Error

	if err != nil { // 删除错误，返回异常错误信息
		return web.JsonErrorMsg(err.Error())
	}

	return web.JsonData(nil) // 返回成功
}
