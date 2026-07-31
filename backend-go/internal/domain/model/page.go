package model

type Page[T any] struct {
	Content       []T   `json:"content"`
	Size          int   `json:"size"`
	Number        int   `json:"number"`
	TotalElements int64 `json:"totalElements"`
	TotalPages    int   `json:"totalPages"`
}

func NewPage[T any](content []T, pageNumber, pageSize int, totalElements int64) Page[T] {
	totalPages := 0
	if pageSize > 0 {
		totalPages = int((totalElements + int64(pageSize) - 1) / int64(pageSize))
	}

	return Page[T]{
		Content:       content,
		Size:          pageSize,
		Number:        pageNumber,
		TotalElements: totalElements,
		TotalPages:    totalPages,
	}
}
