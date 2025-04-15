package responses

import (
	"github.com/gofiber/fiber/v2"
	"github.com/winartodev/apollo/core/errors"
	"github.com/winartodev/apollo/core/helpers"
)

const (
	statusSuccess = "Success"
	statusFailed  = "Failed"
)

type Header struct {
	Status    string `json:"status,omitempty"`
	Message   string `json:"message"`
	Reason    string `json:"reason,omitempty"`
	ErrorCode string `json:"error_code,omitempty"`
}

type Response struct {
	Header   Header      `json:"header"`
	Data     interface{} `json:"data"`
	Metadata interface{} `json:"metadata"`
}

type UpdateResponseData struct {
	Message         string        `json:"message,omitempty"`
	TotalData       int64         `json:"total_data"`
	SuccessData     int64         `json:"success_data"`
	FailData        int64         `json:"fail_data"`
	SuccessRowsData []interface{} `json:"success_rows_data,omitempty"`
	FailedRowsData  []interface{} `json:"failed_rows_data,omitempty"`
}

func SuccessResponse(c *fiber.Ctx, statusCode int, message string, data interface{}, metadata interface{}) error {
	return c.Status(statusCode).JSON(Response{
		Header: Header{
			Status:  statusSuccess,
			Message: "Your request has been processed successfully",
		},
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
		Header: Header{
			Status:  statusFailed,
			Message: message,
			Reason:  e,
		},
		Data: nil,
	})
}

func SuccessResponseV2(c *fiber.Ctx, statusCode int, data interface{}, metadata interface{}) error {
	return c.Status(statusCode).JSON(Response{
		Header: Header{
			Status:  statusSuccess,
			Message: "Your request has been processed successfully",
		},
		Data:     data,
		Metadata: metadata,
	})
}

func SuccessResponseWithPaginate(c *fiber.Ctx, statusCode int, data interface{}, totalItem int64, paginate *helpers.Paginate, metadata interface{}) error {
	var paginateResponse helpers.PaginateResponse
	paginationMetadata := paginateResponse.NewFromContext(c, totalItem, paginate)
	finalMetadata := map[string]interface{}{
		"pagination": paginationMetadata,
	}

	if metadata != nil {
		if existingMetadata, ok := metadata.(map[string]interface{}); ok {
			for key, value := range existingMetadata {
				finalMetadata[key] = value
			}
		} else {
			// Handle the case where the provided metadata is not a map
			finalMetadata["additional_metadata"] = metadata
		}
	}

	return SuccessResponseV2(c, statusCode, data, finalMetadata)
}

func FailedResponseV2(c *fiber.Ctx, err errors.Errors) error {
	var data = err.Error()
	return c.Status(data.StatusCode).JSON(Response{
		Header: Header{
			Message:   data.Message,
			Reason:    data.Reason,
			ErrorCode: data.ErrorCode,
		},
		Data: nil,
	})
}
