package mcp

import (
	"testing"

	mcpgo "github.com/mark3labs/mcp-go/mcp"
	"github.com/stretchr/testify/require"
)

func TestToolResultError_UsesTextContent(t *testing.T) {
	result := mcpgo.NewToolResultError("bad input")
	require.True(t, result.IsError)
	require.Contains(t, textFromResult(t, result), "bad input")
}
