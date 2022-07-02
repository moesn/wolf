package db

import (
	"github.com/kataras/iris/v12"
	"github.com/moesn/wolf/sqls"
	"github.com/moesn/wolf/web"
)

// 增加
func Create(ctx iris.Context, model interface{}) *web.JsonResult {
	err := ctx.ReadJSON(&model) // 读取Json请求参数

	if err != nil { // 读取Json错误，返回请求参数格式错误
		return web.JsonErrorMsg(err.Error())
	}

	err = web.Verify(model) // 校验参数合法性
	if err != nil {          // 参数有误，返回参数错误信息
		return web.JsonErrorMsg(err.Error())
	}

	err = sqls.DB().Create(model).Error // 增加数据
	if err != nil {                     // 增加错误，返回异常错误信息
		return web.JsonErrorMsg(err.Error())
	}

	return web.JsonData(model) // 返回成功
}
