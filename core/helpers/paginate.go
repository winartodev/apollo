package helpers

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/winartodev/apollo/core"
	"github.com/winartodev/apollo/core/errors"
	"strings"
)

var (
	allowedSortBy = map[string]bool{
		"asc":  true,
		"desc": true,
	}
)

type Paginate struct {
	Limit   *int64  `json:"limit" query:"limit"`
	Offset  *int64  `json:"page" query:"page"`
	OrderBy *string `json:"order_by" query:"order_by"`
	SortBy  *string `json:"sort_by" query:"sort_direction"`

	ValidOrderOptions map[string]bool
}

func (f *Paginate) BuildDefault() *Paginate {
	limitValue := core.DefaultLimit
	offsetValue := core.DefaultOffset
	orderByValue := core.DefaultOrder
	sortByValue := core.DefaultSort

	return &Paginate{
		Limit:   &limitValue,
		Offset:  &offsetValue,
		OrderBy: &orderByValue,
		SortBy:  &sortByValue,
	}
}

func (f *Paginate) NewFromContext(ctx *fiber.Ctx) (res *Paginate, err errors.Errors) {
	res = f
	ctxErr := ctx.QueryParser(res)
	if ctxErr != nil {
		return nil, errors.BadRequestErr.WithReason("failed to decode query parameters")
	}

	res.Validate()

	return res, nil
}

func (f *Paginate) Validate() {
	defaultOrder := core.DefaultOrder
	if !f.validateOrderColumns() {
		f.OrderBy = &defaultOrder
	}

	defaultSort := core.DefaultSort
	if !f.validateSortBy() {
		f.SortBy = &defaultSort
	}

	defaultOffset := core.DefaultOffset
	if !f.validateOffset() {
		f.Offset = &defaultOffset
	}

	defaultLimit := core.DefaultLimit
	if !f.validateLimit() {
		f.Limit = &defaultLimit
	}
}

func (f *Paginate) validateOrderColumns() bool {
	if f.OrderBy == nil || len(*f.OrderBy) == 0 {
		return false
	}

	if f.ValidOrderOptions == nil {
		return false
	}

	orderByLower := strings.ToLower(*f.OrderBy)
	if isAllowed := f.ValidOrderOptions[orderByLower]; !isAllowed {
		return false
	}

	return true
}

func (f *Paginate) validateSortBy() bool {
	if f.SortBy == nil || len(*f.SortBy) == 0 {
		return false
	}

	sortDirection := strings.ToLower(*f.SortBy)
	if isAllowed := allowedSortBy[sortDirection]; !isAllowed {
		return false
	}

	return true
}

func (f *Paginate) validateOffset() bool {
	if f.Offset == nil || *f.Offset < 0 {
		return false
	}

	offset := f.calculateOffset(f.Offset, f.Limit)
	f.Offset = &offset

	return true
}

func (f *Paginate) validateLimit() bool {
	if f.Limit == nil || *f.Limit < 1 {
		return false
	} else if *f.Limit > core.MaxLimit {
		return false
	}

	return true
}

func (f *Paginate) calculateOffset(offset *int64, limit *int64) int64 {
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

func (f *Paginate) BuildQueryParams() string {
	var queries []string

	if *f.Offset != core.DefaultOffset {
		queries = append(queries, fmt.Sprintf("page=%d", *f.Offset))
	}

	if *f.Limit != core.DefaultLimit {
		queries = append(queries, fmt.Sprintf("limit=%d", *f.Limit))
	}

	if *f.OrderBy != core.DefaultOrder {
		queries = append(queries, fmt.Sprintf("order_by=%s", *f.OrderBy))
	}

	if *f.SortBy != core.DefaultSort {
		queries = append(queries, fmt.Sprintf("sort_direction=%s", *f.SortBy))
	}

	return fmt.Sprintf("%s", strings.Join(queries, "&"))
}

func (f *Paginate) BuildSQLQuery() string {
	return fmt.Sprintf("ORDER BY %s %s LIMIT $1 OFFSET $2", *f.OrderBy, strings.ToUpper(*f.SortBy))
}

type PaginateResponse struct {
	TotalItems  int64  `json:"total_items"`
	TotalPages  int64  `json:"total_pages"`
	CurrentPage int64  `json:"current_page"`
	NextPage    *int64 `json:"next_page,omitempty"`
	PrevPage    *int64 `json:"prev_page,omitempty"`
	Links       struct {
		First    string `json:"first,omitempty"`
		Last     string `json:"last,omitempty"`
		Next     string `json:"next,omitempty"`
		Previous string `json:"previous,omitempty"`
	} `json:"links,omitempty"`
}

func (pr *PaginateResponse) NewFromContext(ctx *fiber.Ctx, totalItems int64, paginate *Paginate) PaginateResponse {
	var link = fmt.Sprintf("%s%s", ctx.BaseURL(), ctx.Path())
	return pr.Build(totalItems, link, paginate)
}

func (pr *PaginateResponse) Build(totalItems int64, link string, paginate *Paginate) PaginateResponse {
	var limit = *paginate.Limit
	var page = *paginate.Offset

	if page <= 0 {
		page = core.DefaultPage
	}

	if limit < 1 || limit > core.MaxLimit {
		limit = core.DefaultLimit
	}

	totalPages := (totalItems + limit - 1) / limit

	var nextPage, prevPage *int64
	if page < totalPages {
		nextPageVal := page + 1
		nextPage = &nextPageVal
	}

	if page > 1 {
		prevPageVal := page - 1
		prevPage = &prevPageVal
	}

	firstLink := generateLink(link, 1, paginate)
	lastLink := generateLink(link, totalPages, paginate)

	var nextLink, prevLink string

	if nextPage != nil {
		nextLink = generateLink(link, *nextPage, paginate)
	}

	if prevPage != nil {
		prevLink = generateLink(link, *prevPage, paginate)
	}

	result := PaginateResponse{
		TotalItems:  totalItems,
		TotalPages:  totalPages,
		CurrentPage: page,
		NextPage:    nextPage,
		PrevPage:    prevPage,
		Links: struct {
			First    string `json:"first,omitempty"`
			Last     string `json:"last,omitempty"`
			Next     string `json:"next,omitempty"`
			Previous string `json:"previous,omitempty"`
		}{
			First:    firstLink,
			Last:     lastLink,
			Next:     nextLink,
			Previous: prevLink,
		},
	}

	return result
}

func generateLink(link string, page int64, paginate *Paginate) string {
	if link == "" {
		return ""
	}

	paginate.Offset = &page

	return fmt.Sprintf("%s?%s", link, paginate.BuildQueryParams())
}
