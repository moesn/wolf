package db

import (
	"encoding/json"
	"github.com/kataras/iris/v12"
	"github.com/moesn/wolf/sqls"
	"github.com/moesn/wolf/web"
	"github.com/tidwall/gjson"
)

// 查询结构体
type QueryJson struct {
	limit  int
	offset int
	filter interface{}
	sort   interface{}
	fuzzy  struct {
		field   []string
		keyword string
	}
	exact   interface{}
	exclude interface{}
	period  interface{}
}

// 查询信息
func QueryBy(id string, model interface{}) *web.JsonResult {
	if err := sqls.DB().First(model, "id = ?", id).Error; err != nil {
		return web.JsonErrorMsg(err.Error())
	}
	return web.JsonData(model) // 返回数据
}

// 查询列表
func QueryList(ctx iris.Context, model interface{}) *web.JsonResult {
	var params interface{}
	err := ctx.ReadJSON(&params) // 读取Json请求参数

	if err != nil { // 读取Json错误，返回请求参数格式错误
		return web.JsonErrorMsg(err.Error())
	}

	jsonb, _ := json.Marshal(params) // Json转Byte数组
	jsons := string(jsonb)           // Byte数组转字符串

	sql := sqls.NewCnd() // Sql查询条件

	page := gjson.Get(jsons, "page").Int()   // 第几页
	limit := gjson.Get(jsons, "limit").Int() // 每页条数

	if page != 0 && limit != 0 { // 分页参数不为空
		sql.Page(int(page), int(limit))
	}

	field := gjson.Get(jsons, "fuzzy.field").Array()      // 模糊查询字段
	keyword := gjson.Get(jsons, "fuzzy.keyword").String() // 模糊查询关键字

	if len(field) > 0 && keyword != "" { // 有效的模糊查询参数
		for _, column := range field {
			sql.Like(column.String(), keyword)
		}
	}

	sql.Find(sqls.DB(), model) // 查询数据

	count := sql.Count(sqls.DB(), model) //查询条数

	return web.JsonPageData(model, count) // 返回数据
}
