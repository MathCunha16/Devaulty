package common

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func ExtractUUIDParam(c *gin.Context, paramName string) (uuid.UUID, error) {
	paramVal := c.Param(paramName)
	if paramVal == "" {
		return uuid.Nil, fmt.Errorf("missing %s param", paramName)
	}

	id, err := uuid.Parse(paramVal)
	if err != nil {
		return uuid.Nil, fmt.Errorf("invalid %s format", paramName)
	}

	return id, nil
}
