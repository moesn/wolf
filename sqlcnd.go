package wolf

import (
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

// 分页请求参数
type Paging struct {
	Page  int `json:"page"`  // 页码
	Limit int `json:"limit"` // 每页条数
}

// 计算查询偏移量
func (p *Paging) Offset() int {
	offset := 0
	if p.Page > 0 {
		offset = (p.Page - 1) * p.Limit
	}
	return offset
}

// 查询参数
type ParamPair struct {
	Query string        // 查询
	Args  []interface{} // 参数
}

// 排序
type OrderByCol struct {
	Column string // 排序字段
	Asc    bool   // 是否正序
}

// SQL查询条件
type SqlCnd struct {
	SelectCols []string     // 要查询的字段，如果为空，表示查询所有字段
	Params     []ParamPair  // 查询参数
	Orders     []OrderByCol // 排序
	Paging     *Paging      // 分页
}

func NewSqlCnd() *SqlCnd {
	return &SqlCnd{}
}

func (s *SqlCnd) Like(column string, str string) *SqlCnd {
	s.Where(column+" LIKE ?", "%"+str+"%")
	return s
}

func (s *SqlCnd) Where(query string, args ...interface{}) *SqlCnd {
	s.Params = append(s.Params, ParamPair{query, args})
	return s
}

// 升序
func (s *SqlCnd) Asc(column string) *SqlCnd {
	s.Orders = append(s.Orders, OrderByCol{Column: column, Asc: true})
	return s
}

// 倒序
func (s *SqlCnd) Desc(column string) *SqlCnd {
	s.Orders = append(s.Orders, OrderByCol{Column: column, Asc: false})
	return s
}

// 分页
func (s *SqlCnd) Page(page, limit int) *SqlCnd {
	if s.Paging == nil {
		s.Paging = &Paging{Page: page, Limit: limit}
	} else {
		s.Paging.Page = page
		s.Paging.Limit = limit
	}
	return s
}

func (s *SqlCnd) Build(db *gorm.DB, model interface{}) *gorm.DB {
	ret := db.Model(model)

	if len(s.SelectCols) > 0 { // 设置查询字段
		ret = ret.Select(s.SelectCols)
	}

	if len(s.Params) > 0 { // 设置查询参数
		for _, param := range s.Params {
			ret = ret.Where(param.Query, param.Args...)
		}
	}

	if len(s.Orders) > 0 { // 设置排序
		for _, order := range s.Orders {
			if order.Asc {
				ret = ret.Order(order.Column + " ASC")
			} else {
				ret = ret.Order(order.Column + " DESC")
			}
		}
	}

	if s.Paging != nil && s.Paging.Limit > 0 { // 条数
		ret = ret.Limit(s.Paging.Limit)
	}

	if s.Paging != nil && s.Paging.Offset() > 0 { // 偏移
		ret = ret.Offset(s.Paging.Offset())
	}

	return ret
}

func (s *SqlCnd) Find(db *gorm.DB, out interface{}) {
	if err := s.Build(db, out).Find(out).Error; err != nil {
		logrus.Error(err)
	}
}

func (s *SqlCnd) Count(db *gorm.DB, model interface{}) int64 {
	ret := db.Model(model)

	if len(s.Params) > 0 { // 设置查询参数
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
