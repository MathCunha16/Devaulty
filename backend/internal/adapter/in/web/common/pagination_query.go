package common

type PaginationQuery struct {
	PageNumber int `form:"page,default=0" binding:"gte=0"`
	PageSize   int `form:"size,default=10" binding:"gte=1,lte=100"`
}
