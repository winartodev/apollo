package helpers

import (
	"github.com/winartodev/apollo/core"
	"strings"
)

type Paginate struct {
	Limit   *int64
	Offset  *int64
	OrderBy *string
	SortBy  *string
}

func (f *Paginate) Validate() {
	if f.OrderBy == nil || len(*f.OrderBy) == 0 {
		ob := core.DefaultOrder
		f.OrderBy = &ob
	}

	if f.SortBy == nil || len(*f.SortBy) == 0 {
		sb := core.DefaultSort
		f.SortBy = &sb
	} else {
		sb := strings.ToLower(*f.SortBy)
		if sb != "asc" && sb != "desc" {
			sb = core.DefaultSort
			f.SortBy = &sb
		}
	}

	if f.Offset == nil || *f.Offset < 0 {
		o := core.DefaultOffset
		f.Offset = &o
	}

	if f.Limit == nil || *f.Limit < 1 {
		lmt := core.DefaultLimit
		f.Limit = &lmt
	} else {
		if *f.Limit < 1 || *f.Limit > core.MaxLimit {
			lmt := core.DefaultLimit
			f.Limit = &lmt
		}
	}
}

func (f *Paginate) BuildQueryParam() string {
	return ""
}

func (f *Paginate) CalculateOffset(offset *int64, limit *int64) int64 {
	if offset == nil {
		return 0
	}

	if limit == nil {
		return *offset
	}

	result := (*offset - 1) * *limit
	if result < 0 {
		result = 0
	}

	return result
}
