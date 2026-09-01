package util

import (
	"fmt"

	"github.com/google/uuid"
	"github.com/mark3labs/mcp-go/mcp"
)

type MCPPaginationQuery struct {
	PageNumber int `json:"page,omitempty"`
	PageSize   int `json:"page_size,omitempty"`
}

type ProjectPaginationQuery struct {
	ProjectID uuid.UUID `json:"projectID"`
	MCPPaginationQuery
}

// ToStrings converts a slice of any ~string type to []string.
func ToStrings[T ~string](values []T) []string {
	result := make([]string, len(values))
	for i, v := range values {
		result[i] = string(v)
	}
	return result
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

func ValidateProjectQuery(query *ProjectPaginationQuery) {
	ValidateQuery(&query.MCPPaginationQuery)
}

func ExtractUUID(req mcp.CallToolRequest, paramName string) (uuid.UUID, error) {
	idStr, err := req.RequireString(paramName)
	if err != nil {
		return uuid.Nil, err
	}

	id, err := uuid.Parse(idStr)
	if err != nil {
		return uuid.Nil, fmt.Errorf("%s: invalid UUID format", paramName)
	}

	return id, nil
}
