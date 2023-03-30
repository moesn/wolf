package db

import (
	"github.com/kataras/iris/v12"
	"github.com/moesn/wolf/http"
	"github.com/moesn/wolf/http/params"
	"github.com/moesn/wolf/sql"
)

func QueryById(id string, model interface{}) *http.JsonResult {
	if err := DB().First(model, "id = ?", id).Error; err != nil {
		return http.JsonErrorMsg(err.Error())
	}
	return http.JsonData(model)
}

func QueryList(ctx iris.Context, model interface{}) *http.JsonResult {
	return QueryListPrePost(ctx,model,nil,nil)
}

func QueryListPre(ctx iris.Context, model interface{},preProcess func(sqlCnd * sql.Cnd, jsonParams string)) *http.JsonResult {
	return QueryListPrePost(ctx,model,preProcess,nil)
}

func QueryListPost(ctx iris.Context, model interface{},postProcess func() interface{}) *http.JsonResult {
	return QueryListPrePost(ctx,model,nil,postProcess)
}

func QueryListPrePost(ctx iris.Context, model interface{},preProcess func(sqlCnd * sql.Cnd, jsonParams string) ,postProcess func() interface{}) *http.JsonResult {
	json, err := params.ReadJson(ctx)

	if err != nil {
		return http.JsonError(err)
	}

	var (
		page    = params.GetInt("page", json)
		limit   = params.GetInt("limit", json)
		field   = params.GetArray("fuzzy.field", json)
		keyword = params.GetString("fuzzy.keyword", json)
		exact   = params.GetMap("exact", json)
		filter    = params.GetMap("filter", json)
		sort    = params.GetMap("sort", json)
	)

	sql := sql.NewCnd()

	if field != nil && keyword != "" {
			sql.Like(field, keyword)
	}

	if exact != nil {
		for column, val := range exact {
			sql.Eq(column, val)
		}
	}

	if filter != nil {
		for column, val := range filter {
			sql.In(column,val)
		}
	}

	if sort != nil {
		for column, val := range sort {
			if val.(string) == "asc" {
				sql.Asc(column)
			} else if val.(string) == "desc" {
				sql.Desc(column)
			}
		}
	}

	if page != 0 && limit != 0 {
		sql.Page(int(page), int(limit))
	}

	if preProcess != nil {
		preProcess(sql,json)
	}

	if len(sql.Orders)==0{
		sql.Desc("id")
	}

	sql.Find(DB(), model)

	count := sql.Count(DB(), model)

	if postProcess != nil {
		return http.JsonPageData(postProcess(), count)
	}

	return http.JsonPageData(model, count)
}
