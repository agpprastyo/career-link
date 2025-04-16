package pagination

import (
	"math"
	"strconv"

	"github.com/gofiber/fiber/v2"
)

// Pagination represents pagination parameters
type Pagination struct {
	Page     int `json:"page"`
	PageSize int `json:"page_size"`
	Total    int `json:"total"`
	Pages    int `json:"pages"`
}

// PageResponse wraps any data with pagination metadata
type PageResponse struct {
	Data       interface{} `json:"data"`
	Pagination Pagination  `json:"pagination"`
}

// DefaultPage is the default page number
const DefaultPage = 1

// DefaultPageSize is the default items per page
const DefaultPageSize = 10

// MaxPageSize is the maximum allowed page size
const MaxPageSize = 100

// ExtractFromRequest extracts pagination parameters from request query params
func ExtractFromRequest(c *fiber.Ctx) Pagination {
	// Extract page and page_size from query params
	pageStr := c.Query("page", strconv.Itoa(DefaultPage))
	pageSizeStr := c.Query("page_size", strconv.Itoa(DefaultPageSize))

	// Parse to integers with defaults
	page, err := strconv.Atoi(pageStr)
	if err != nil || page < 1 {
		page = DefaultPage
	}

	pageSize, err := strconv.Atoi(pageSizeStr)
	if err != nil || pageSize < 1 {
		pageSize = DefaultPageSize
	}

	// Enforce maximum page size
	if pageSize > MaxPageSize {
		pageSize = MaxPageSize
	}

	return Pagination{
		Page:     page,
		PageSize: pageSize,
	}
}

// CalculateOffset calculates the offset for SQL queries
func (p Pagination) CalculateOffset() int {
	return (p.Page - 1) * p.PageSize
}

// WithTotal creates a complete pagination object with total items count
func (p Pagination) WithTotal(totalItems int) Pagination {
	return Pagination{
		Page:     p.Page,
		PageSize: p.PageSize,
		Total:    totalItems,
		Pages:    int(math.Ceil(float64(totalItems) / float64(p.PageSize))),
	}
}

// NewResponse creates a paginated response with data and metadata
func NewResponse(data interface{}, pagination Pagination, totalItems int) PageResponse {
	return PageResponse{
		Data:       data,
		Pagination: pagination.WithTotal(totalItems),
	}
}

// GetSQLLimitOffset returns the LIMIT and OFFSET SQL clause for pagination

func (p Pagination) GetSQLLimitOffset() string {
	return " LIMIT " + strconv.Itoa(p.PageSize) + " OFFSET " + strconv.Itoa(p.CalculateOffset())
}
