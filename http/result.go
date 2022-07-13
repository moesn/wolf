package http

import (
	"github.com/moesn/wolf/sql"
)


type PageResult struct {
	Page    *sql.Paging `json:"page"`
	Results interface{} `json:"results"`
}

type CursorResult struct {
	Results interface{} `json:"results"`
	Cursor  string      `json:"cursor"`
	HasMore bool        `json:"hasMore"`
}
