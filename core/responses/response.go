package responses

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/winartodev/apollo/core"
	"github.com/winartodev/apollo/core/helpers"
)

const (
	statusSuccess = "Success"
	statusFailed  = "Failed"
)

type Response struct {
	Status   string      `json:"status"`
	Message  string      `json:"message"`
	Data     interface{} `json:"data,omitempty"`
	Metadata interface{} `json:"metadata,omitempty"`
	Error    string      `json:"error,omitempty"`
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

func SuccessResponse(c *fiber.Ctx, statusCode int, message string, data interface{}, metadata interface{}) error {
	return c.Status(statusCode).JSON(Response{
		Status:   statusSuccess,
		Message:  message,
		Data:     data,
		Metadata: metadata,
	})
}

func FailedResponse(c *fiber.Ctx, statusCode int, message string, err error) error {
	var e string
	if err != nil {
		e = err.Error()
	}

	return c.Status(statusCode).JSON(Response{
		Status:  statusFailed,
		Message: message,
		Error:   e,
	})
}
func generateLink(link string, page int64, limit int64) string {
	if link == "" {
		return ""
	}

	return fmt.Sprintf("%s?page=%d&limit=%d", link, page, limit)
}

func BuildPaginate(totalItems int64, link string, paginate *helpers.Paginate) *PaginateResponse {
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

	firstLink := generateLink(link, 1, limit)
	lastLink := generateLink(link, totalPages, limit)

	var nextLink, prevLink string

	if nextPage != nil {
		nextLink = generateLink(link, *nextPage, limit)
	}

	if prevPage != nil {
		prevLink = generateLink(link, *prevPage, limit)
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

	return &result
}
