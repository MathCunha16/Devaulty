package util

import (
	"fmt"

	"github.com/google/uuid"
	"github.com/mark3labs/mcp-go/mcp"
)

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
