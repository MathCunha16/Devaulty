package util

import "github.com/mark3labs/mcp-go/mcp"

var ReadOnlyAnnotations = mcp.ToolAnnotation{
	ReadOnlyHint:    boolPtr(true),
	DestructiveHint: boolPtr(false),
	OpenWorldHint:   boolPtr(false),
}

var WriteAnnotations = mcp.ToolAnnotation{
	ReadOnlyHint:    boolPtr(false),
	DestructiveHint: boolPtr(false),
	OpenWorldHint:   boolPtr(false),
}

var DeleteAnnotations = mcp.ToolAnnotation{
	ReadOnlyHint:    boolPtr(false),
	DestructiveHint: boolPtr(true),
	OpenWorldHint:   boolPtr(false),
}

func boolPtr(v bool) *bool {
	return &v
}
