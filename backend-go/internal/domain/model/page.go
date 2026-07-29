package model

import "math"

type Page[T any] struct {
	Content       []T   `json:"content"`
	Size          int   `json:"size"`
	Number        int   `json:"number"`
	TotalElements int64 `json:"totalElements"`
	TotalPages    int   `json:"totalPages"`
}

func NewPage[T any](content []T, pageNumber, PageSize int, totalElements int64) Page[T] {
	totalPages := 0
	if PageSize > 0 {
		totalPages = int(math.Ceil(float64(totalElements) / float64(PageSize)))
	}

	return Page[T]{
		Content:       content,
		Size:          PageSize,
		Number:        pageNumber,
		TotalElements: totalElements,
		TotalPages:    totalPages,
	}
}
