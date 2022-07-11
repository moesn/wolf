package db

import (
	"github.com/kataras/iris/v12"
	"github.com/moesn/wolf/http"
	"github.com/moesn/wolf/http/params"
	"github.com/moesn/wolf/sql"
	"github.com/samber/lo"
)

func QueryBy(id string, model interface{}) *http.JsonResult {
	if err := DB().First(model, "id = ?", id).Error; err != nil {
		return http.JsonErrorMsg(err.Error())
	}
	return http.JsonData(model)
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
		page    = params.GetInt("page", json)
		limit   = params.GetInt("limit", json)
		field   = params.GetArray("fuzzy.field", json, extendParams.Fuzzy.Field)
		keyword = params.GetString("fuzzy.keyword", json, extendParams.Fuzzy.Keyword)
		exact   = lo.Assign(params.GetMap("exact", json), extendParams.Exact)
		filter    = lo.Assign(params.GetMap("filter", json), extendParams.Filter)
		sort    = lo.Assign(params.GetMap("sort", json), extendParams.Sort)
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

	sql.Find(DB(), model)

	count := sql.Count(DB(), model)

	if process != nil {
		return http.JsonPageData(process(), count)
	}

	return http.JsonPageData(model, count)
}
