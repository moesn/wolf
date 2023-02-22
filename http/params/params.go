package params

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/moesn/wolf/http"
	"github.com/moesn/wolf/sql"
	"github.com/tidwall/gjson"
	"strconv"
	"strings"
	"time"

	"github.com/iris-contrib/schema"
	"github.com/kataras/iris/v12"
	"github.com/moesn/wolf/common/dates"
	"github.com/moesn/wolf/common/strs"
)

type Fuzzy  struct {
	Field   []string
	Keyword string
}

type HttpParams struct {
	Page  int
	Limit  int
	Total int
	Filter map[string]interface{}
	Fuzzy map[string]interface{}
	Sort   map[string]interface{}
	Exact   map[string]interface{}
	Exclude map[string]interface{}
	Range  map[string]interface{}
}

var (
	decoder = schema.NewDecoder()
)

func init() {
	decoder.AddAliasTag("form", "json")
	decoder.ZeroEmpty(true)
}

func ReadJson(ctx iris.Context) (string, *http.CodeError) {
	if ctx==nil {
		return 	"",nil
	}

	var params interface{}
	err := ctx.ReadJSON(&params)

	if err != nil {
		return "", http.NewError(1, "错误的JSON请求参数")
	}

	jsonb, _ := json.Marshal(params)
	jsons := string(jsonb)

	return jsons, nil
}

func GetString(name string, json string, deft ...string) string {
	str := gjson.Get(json, name).String()
	if str == "" && len(deft) != 0 {
		str = deft[0]
	}
	return str
}

func GetInt(name string, json string) int64 {
	return gjson.Get(json, name).Int()
}

func GetMap(name string, json string) map[string]interface{} {
	resMap := map[string]interface{}{}

	for key, val := range gjson.Get(json, name).Map() {
		resMap[key] = val.Value()
	}
	return resMap
}

func GetArray(name string, json string, deft ...[]string) []string {
	resArray := make([]string, 0)

	for _, val := range gjson.Get(json, name).Array() {
		resArray = append(resArray, val.String())
	}

	if len(resArray) == 0 && len(deft) != 0 {
		resArray = deft[0]
	}
	return resArray
}

func GetResult(name string, json string) gjson.Result {
	return gjson.Get(json, name)
}

func paramError(name string) error {
	return errors.New(fmt.Sprintf("无法找到参数值 '%s'", name))
}

func ReadForm(ctx iris.Context, obj interface{}) error {
	values := ctx.FormValues()
	if len(values) == 0 {
		return nil
	}
	return decoder.Decode(obj, values)
}

func FormValue(ctx iris.Context, name string) string {
	return ctx.FormValue(name)
}

func FormValueRequired(ctx iris.Context, name string) (string, error) {
	str := FormValue(ctx, name)
	if len(str) == 0 {
		return "", errors.New("参数：" + name + "不能为空")
	}
	return str, nil
}

func FormValueDefault(ctx iris.Context, name, def string) string {
	return ctx.FormValueDefault(name, def)
}

func FormValueInt(ctx iris.Context, name string) (int, error) {
	str := ctx.FormValue(name)
	if str == "" {
		return 0, paramError(name)
	}
	return strconv.Atoi(str)
}

func FormValueIntDefault(ctx iris.Context, name string, def int) int {
	if v, err := FormValueInt(ctx, name); err == nil {
		return v
	}
	return def
}

func FormValueInt64(ctx iris.Context, name string) (int64, error) {
	str := ctx.FormValue(name)
	if str == "" {
		return 0, paramError(name)
	}
	return strconv.ParseInt(str, 10, 64)
}

func FormValueInt64Default(ctx iris.Context, name string, def int64) int64 {
	if v, err := FormValueInt64(ctx, name); err == nil {
		return v
	}
	return def
}

func FormValueInt64Array(ctx iris.Context, name string) []int64 {
	str := ctx.FormValue(name)
	if str == "" {
		return nil
	}
	ss := strings.Split(str, ",")
	if len(ss) == 0 {
		return nil
	}
	var ret []int64
	for _, v := range ss {
		item, err := strconv.ParseInt(v, 10, 64)
		if err != nil {
			continue
		}
		ret = append(ret, item)
	}
	return ret
}

func FormValueStringArray(ctx iris.Context, name string) []string {
	str := ctx.FormValue(name)
	if len(str) == 0 {
		return nil
	}
	ss := strings.Split(str, ",")
	if len(ss) == 0 {
		return nil
	}
	var ret []string
	for _, s := range ss {
		s = strings.TrimSpace(s)
		if len(s) == 0 {
			continue
		}
		ret = append(ret, s)
	}
	return ret
}

func FormValueBool(ctx iris.Context, name string) (bool, error) {
	str := ctx.FormValue(name)
	if str == "" {
		return false, paramError(name)
	}
	return strconv.ParseBool(str)
}

func FormValueBoolDefault(ctx iris.Context, name string, def bool) bool {
	str := ctx.FormValue(name)
	if str == "" {
		return def
	}
	value, err := strconv.ParseBool(str)
	if err != nil {
		return def
	}
	return value
}

func FormDate(ctx iris.Context, name string) *time.Time {
	value := FormValue(ctx, name)
	if strs.IsBlank(value) {
		return nil
	}
	layouts := []string{dates.FmtDateTime, dates.FmtDate, dates.FmtDateTimeNoSeconds}
	for _, layout := range layouts {
		if ret, err := dates.Parse(value, layout); err == nil {
			return &ret
		}
	}
	return nil
}

func GetPaging(ctx iris.Context) *sql.Paging {
	page := FormValueIntDefault(ctx, "page", 1)
	limit := FormValueIntDefault(ctx, "limit", 20)
	if page <= 0 {
		page = 1
	}
	if limit <= 0 {
		limit = 20
	}
	return &sql.Paging{Page: page, Limit: limit}
}
