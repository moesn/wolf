package sql

import (
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

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

func (s *Cnd) Cols(selectCols ...string) *Cnd {
	if len(selectCols) > 0 {
		s.SelectCols = append(s.SelectCols, selectCols...)
	}
	return s
}

func (s *Cnd) Eq(column string, args ...interface{}) *Cnd {
	s.Where(column+" = ?", args)
	return s
}

func (s *Cnd) NotEq(column string, args ...interface{}) *Cnd {
	s.Where(column+" <> ?", args)
	return s
}

func (s *Cnd) Gt(column string, args ...interface{}) *Cnd {
	s.Where(column+" > ?", args)
	return s
}

func (s *Cnd) Gte(column string, args ...interface{}) *Cnd {
	s.Where(column+" >= ?", args)
	return s
}

func (s *Cnd) Lt(column string, args ...interface{}) *Cnd {
	s.Where(column+" < ?", args)
	return s
}

func (s *Cnd) Lte(column string, args ...interface{}) *Cnd {
	s.Where(column+" <= ?", args)
	return s
}

func (s *Cnd) Like(columns []string, str string) *Cnd {
	likeCnd := ""
	for i, column := range columns {
		likeCnd += column + " LIKE " + "'%" + str + "%'"
		if i < len(columns)-1 {
			likeCnd += " OR "
		}
	}
	s.Where(likeCnd)
	return s
}

func (s *Cnd) Starting(column string, str string) *Cnd {
	s.Where(column+" LIKE ?", str+"%")
	return s
}

func (s *Cnd) Ending(column string, str string) *Cnd {
	s.Where(column+" LIKE ?", "%"+str)
	return s
}

func (s *Cnd) In(column string, params interface{}) *Cnd {
	s.Where(column+" in (?) ", params)
	return s
}

func (s *Cnd) NotIn(column string, params interface{}) *Cnd {
	s.Where(column+" not in (?) ", params)
	return s
}

func (s *Cnd) Where(query string, args ...interface{}) *Cnd {
	s.Params = append(s.Params, ParamPair{query, args})
	return s
}

func (s *Cnd) Asc(column string) *Cnd {
	s.Orders = append(s.Orders, OrderByCol{Column: column, Asc: true})
	return s
}

func (s *Cnd) Desc(column string) *Cnd {
	s.Orders = append(s.Orders, OrderByCol{Column: column, Asc: false})
	return s
}

func (s *Cnd) Limit(limit int) *Cnd {
	s.Page(1, limit)
	return s
}

func (s *Cnd) Page(page, limit int) *Cnd {
	if s.Paging == nil {
		s.Paging = &Paging{Page: page, Limit: limit}
	} else {
		s.Paging.Page = page
		s.Paging.Limit = limit
	}
	return s
}

func (s *Cnd) Build(db *gorm.DB) *gorm.DB {
	ret := db

	if len(s.SelectCols) > 0 {
		ret = ret.Select(s.SelectCols)
	}

	if len(s.Params) > 0 {
		for _, param := range s.Params {
			ret = ret.Where(param.Query, param.Args...)
		}
	}

	if len(s.Orders) > 0 {
		for _, order := range s.Orders {
			if order.Asc {
				ret = ret.Order(order.Column + " ASC")
			} else {
				ret = ret.Order(order.Column + " DESC")
			}
		}
	}

	if s.Paging != nil && s.Paging.Limit > 0 {
		ret = ret.Limit(s.Paging.Limit)
	}

	if s.Paging != nil && s.Paging.Offset() > 0 {
		ret = ret.Offset(s.Paging.Offset())
	}
	return ret
}

func (s *Cnd) Find(db *gorm.DB, out interface{}) {
	if err := s.Build(db).Find(out).Error; err != nil {
		logrus.Error(err)
	}
}

func (s *Cnd) FindOne(db *gorm.DB, out interface{}) error {
	if err := s.Limit(1).Build(db).First(out).Error; err != nil {
		return err
	}
	return nil
}

func (s *Cnd) Count(db *gorm.DB, model interface{}) int64 {
	ret := db.Model(model)

	if len(s.Params) > 0 {
		for _, query := range s.Params {
			ret = ret.Where(query.Query, query.Args...)
		}
	}

	var count int64
	if err := ret.Count(&count).Error; err != nil {
		logrus.Error(err)
	}
	return count
}
