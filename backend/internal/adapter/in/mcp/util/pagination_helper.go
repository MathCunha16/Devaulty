package util

type MCPPaginationQuery struct {
	PageNumber int `json:"page,omitempty"`
	PageSize   int `json:"page_size,omitempty"`
}

func ValidateQuery(query *MCPPaginationQuery) {
	if query.PageNumber <= 0 {
		query.PageNumber = 0
	}
	if query.PageSize <= 0 {
		query.PageSize = 20
	}
	if query.PageSize > 100 {
		query.PageSize = 100
	}
}
