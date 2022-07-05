package db

import (
	"github.com/kataras/iris/v12"
	"github.com/moesn/wolf/http"
	"github.com/moesn/wolf/http/params"
	"github.com/moesn/wolf/sql"
	"github.com/tidwall/gjson"
)

// 查询信息
func QueryBy(id string, model interface{}) *http.JsonResult {
	if err := DB().First(model, "id = ?", id).Error; err != nil {
		return http.JsonErrorMsg(err.Error())
	}
	return http.JsonData(model) // 返回数据
}

func QueryList(ctx iris.Context, model interface{}) *http.JsonResult {
	return QueryListExtendProcess(ctx,model,params.HttpParams{},nil)
}

func QueryListExtend(ctx iris.Context, model interface{},extendParams params.HttpParams) *http.JsonResult {
	return QueryListExtendProcess(ctx,model,extendParams,nil)
}

func QueryListProcess(ctx iris.Context, model interface{},process func() interface{}) *http.JsonResult {
	return QueryListExtendProcess(ctx,model,params.HttpParams{},process)
}

func QueryListExtendProcess(ctx iris.Context, model interface{},extendParams params.HttpParams,process func() interface{}) *http.JsonResult {
	json, err := params.ReadJson(ctx)

	if err != nil {
		return http.JsonError(err)
	}

	var (
		page = params.GetInt("page", json)
		limit = params.GetInt("limit", json)
		field = params.GetResult("fuzzy.field", json)
		keyword = params.GetString("fuzzy.keyword", json)

		exact = extendParams.Exact
		sort = extendParams.Sort
		fuzzy = extendParams.Fuzzy
	)


	sql := sql.NewCnd()


	if page != 0 && limit != 0 {
		sql.Page(int(page), int(limit))
	}

	if exact!=nil{
		for column, val := range exact {
			sql.Eq(column,val)
		}
	}

	if fuzzy.Field!=nil&&len(fuzzy.Keyword)!=0{
		for _, column := range fuzzy.Field {
			sql.Like(column,fuzzy.Keyword)
		}
	}

	if sort!=nil{
		for column, val := range sort {
			if(val=="asc"){
				sql.Asc(column)
			}else if(val=="desc"){
				sql.Desc(column)
			}
		}
	}

	if field.IsArray() && keyword != "" {
		field.ForEach(func(key, value gjson.Result) bool {
			sql.Like(value.String(), keyword)
			return true
		})
	}

	sql.Find(DB(), model) // 查询数据

	count := sql.Count(DB(), model) //查询条数

	if(process!=nil){
		return http.JsonPageData(process(), count) // 返回数据
	}

	return http.JsonPageData(model, count) // 返回数据
}
