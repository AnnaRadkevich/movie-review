package contracts

import "strconv"

type PaginationSetter interface {
	SetPage(page int)
	SetSize(size int)
}

var _ PaginationSetter = (*PaginatedRequest)(nil)

func (r *PaginatedRequest) SetPage(page int) {
	r.Page = page
}

func (r *PaginatedRequest) SetSize(size int) {
	r.Size = size
}

type PaginatedRequest struct {
	Page int `json:"page" query:"page"`
	Size int `json:"size" query:"size"`
}

type PaginatedResponse[T any] struct {
	Page  int  `json:"page" validate:"min=0"`
	Size  int  `json:"size" validate:"min=0"`
	Total int  `json:"total"`
	Items []*T `json:"items"`
}

func (r *PaginatedRequest) ToQueryParams() map[string]string {
	params := make(map[string]string, 2)
	if r.Page > 0 {
		params["page"] = strconv.Itoa(r.Page)
	}
	if r.Size > 0 {
		params["size"] = strconv.Itoa(r.Size)
	}
	return params
}
