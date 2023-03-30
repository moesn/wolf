package sql

import (
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

var NoConcat bool

type Cnd struct {
	SelectCols []string
	Params     []ParamPair
	Orders     []OrderByCol
	Paging     *Paging
}

type ParamPair struct {
	Query string
	Args  []interface{}
}

type OrderByCol struct {
	Column string
	Asc    bool
}

func NewCnd() *Cnd {
	return &Cnd{}
}

func (cnd *Cnd) Cols(selectCols ...string) *Cnd {
	if len(selectCols) > 0 {
		cnd.SelectCols = append(cnd.SelectCols, selectCols...)
	}
	return cnd
}

func (cnd *Cnd) Eq(column string, args ...interface{}) *Cnd {
	cnd.Where(column+" = ?", args)
	return cnd
}

func (cnd *Cnd) NotEq(column string, args ...interface{}) *Cnd {
	cnd.Where(column+" <> ?", args)
	return cnd
}

func (cnd *Cnd) Gt(column string, args ...interface{}) *Cnd {
	cnd.Where(column+" > ?", args)
	return cnd
}

func (cnd *Cnd) Gte(column string, args ...interface{}) *Cnd {
	cnd.Where(column+" >= ?", args)
	return cnd
}

func (cnd *Cnd) Lt(column string, args ...interface{}) *Cnd {
	cnd.Where(column+" < ?", args)
	return cnd
}

func (cnd *Cnd) Lte(column string, args ...interface{}) *Cnd {
	cnd.Where(column+" <= ?", args)
	return cnd
}

func (cnd *Cnd) Like(columns []string, str string) *Cnd {
	likeCnd := ""

	if NoConcat {
		for i, column := range columns {
			likeCnd += column + " LIKE " + "'%" + str + "%'"
			if i < len(columns)-1 {
				likeCnd += " OR "
			}
		}

		cnd.Where(likeCnd)
	} else {
		var likeArgs []interface{}

		for i, column := range columns {
			likeCnd += column + " LIKE CONCAT('%',CONCAT(?,'%'))"
			if i < len(columns)-1 {
				likeCnd += " OR "
			}

			likeArgs = append(likeArgs, str)
		}

		cnd.Where(likeCnd, likeArgs...)
	}

	return cnd
}

func (cnd *Cnd) Starting(column string, str string) *Cnd {
	cnd.Where(column+" LIKE ?", str+"%")
	return cnd
}

func (cnd *Cnd) Ending(column string, str string) *Cnd {
	cnd.Where(column+" LIKE ?", "%"+str)
	return cnd
}

func (cnd *Cnd) In(column string, params interface{}) *Cnd {
	cnd.Where(column+" in (?) ", params)
	return cnd
}

func (cnd *Cnd) NotIn(column string, params interface{}) *Cnd {
	cnd.Where(column+" not in (?) ", params)
	return cnd
}

func (cnd *Cnd) Where(query string, args ...interface{}) *Cnd {
	cnd.Params = append(cnd.Params, ParamPair{query, args})
	return cnd
}

func (cnd *Cnd) Asc(column string) *Cnd {
	cnd.Orders = append(cnd.Orders, OrderByCol{Column: column, Asc: true})
	return cnd
}

func (cnd *Cnd) Desc(column string) *Cnd {
	cnd.Orders = append(cnd.Orders, OrderByCol{Column: column, Asc: false})
	return cnd
}

func (cnd *Cnd) Limit(limit int) *Cnd {
	cnd.Page(1, limit)
	return cnd
}

func (cnd *Cnd) Page(page, limit int) *Cnd {
	if cnd.Paging == nil {
		cnd.Paging = &Paging{Page: page, Limit: limit}
	} else {
		cnd.Paging.Page = page
		cnd.Paging.Limit = limit
	}
	return cnd
}

func (cnd *Cnd) Build(db *gorm.DB) *gorm.DB {
	ret := db

	if len(cnd.SelectCols) > 0 {
		ret = ret.Select(cnd.SelectCols)
	}

	if len(cnd.Params) > 0 {
		for _, param := range cnd.Params {
			ret = ret.Where(param.Query, param.Args...)
		}
	}

	if len(cnd.Orders) > 0 {
		for _, order := range cnd.Orders {
			if order.Asc {
				ret = ret.Order(order.Column + " ASC")
			} else {
				ret = ret.Order(order.Column + " DESC")
			}
		}
	}

	if cnd.Paging != nil && cnd.Paging.Limit > 0 {
		ret = ret.Limit(cnd.Paging.Limit)
	}

	if cnd.Paging != nil && cnd.Paging.Offset() > 0 {
		ret = ret.Offset(cnd.Paging.Offset())
	}
	return ret
}

func (cnd *Cnd) Find(db *gorm.DB, out interface{}) {
	if err := cnd.Build(db).Find(out).Error; err != nil {
		logrus.Error(err)
	}
}

func (cnd *Cnd) FindOne(db *gorm.DB, out interface{}) error {
	if err := cnd.Limit(1).Build(db).First(out).Error; err != nil {
		return err
	}
	return nil
}

func (cnd *Cnd) Count(db *gorm.DB, model interface{}) int64 {
	ret := db.Model(model)

	if len(cnd.Params) > 0 {
		for _, query := range cnd.Params {
			ret = ret.Where(query.Query, query.Args...)
		}
	}

	var count int64
	if err := ret.Count(&count).Error; err != nil {
		logrus.Error(err)
	}
	return count
}
