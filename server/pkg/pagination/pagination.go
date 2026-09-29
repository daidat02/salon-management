package pagination

// Params tham số phân trang dùng chung cho các API GetAll.
type Params struct {
	Page     int
	PageSize int
}

const (
	DefaultPage     = 1
	DefaultPageSize = 20
	MaxPageSize     = 100
)

// Normalize kẹp giá trị về khoảng hợp lệ.
func Normalize(page, pageSize int) Params {
	if page < 1 {
		page = DefaultPage
	}
	if pageSize < 1 {
		pageSize = DefaultPageSize
	}
	if pageSize > MaxPageSize {
		pageSize = MaxPageSize
	}
	return Params{Page: page, PageSize: pageSize}
}

// Offset dùng cho SQL OFFSET.
func (p Params) Offset() int {
	return (p.Page - 1) * p.PageSize
}

// Result envelope phân trang trả về cho client.
type Result[T any] struct {
	Data       []T `json:"data"`
	Page       int `json:"page"`
	PageSize   int `json:"page_size"`
	Total      int64 `json:"total"`
	TotalPages int64 `json:"total_pages"`
}

// NewResult dựng Result, đảm bảo Data không nil để JSON ra [] thay vì null.
func NewResult[T any](data []T, p Params, total int64) *Result[T] {
	if data == nil {
		data = make([]T, 0)
	}
	totalPages := int64(0)
	if p.PageSize > 0 && total > 0 {
		totalPages = (total + int64(p.PageSize) - 1) / int64(p.PageSize)
	}
	return &Result[T]{
		Data:       data,
		Page:       p.Page,
		PageSize:   p.PageSize,
		Total:      total,
		TotalPages: totalPages,
	}
}
