package common

type PaginationQuery struct {
	PageNumber int `form:"page,default=0"`
	PageSize   int `form:"size,default=10"`
}
