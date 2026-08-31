# MCP Architecture and Conventions

This document explains the conventions used by the Devaulty MCP layer (`backend/internal/adapter/in/mcp`) and the rules that contributors should follow when adding or updating tools.

## What MCP is in this project

The MCP layer is a thin adapter between external clients and the application use cases. It does not contain business logic itself; it validates incoming tool arguments, calls the relevant use case, and returns MCP-native results.

The MCP package is intentionally structured to mirror the existing backend use case architecture:

- `internal/adapter/in/mcp` contains tool definitions and registration
- `internal/usecase` holds business behavior
- `internal/dto` holds request/response contracts
- `internal/domain/model` holds shared enums and domain primitives

## Tool registration pattern

Each entity has its own tool file, for example:

- `project_tools.go`
- `snippet_tools.go`
- `problem_tools.go`
- `card_tools.go`
- `tag_tools.go`
- `item_tag_tools.go`

Each tool type exposes:

- `NewXTools(...)`
- `Register(s *server.MCPServer, opts Options)`
- the method handlers used by the registered MCP tool names

The registration flow follows the same structure across files:

1. Register read-only tools first
2. If `opts.ReadOnly` is set, return early
3. Register write tools
4. If `opts.DisableDelete` is set, skip delete tools

This keeps behavior consistent and predictable across the whole MCP surface.

## Global registration options

The `Options` struct is used to control the visible tool set for a given MCP server instance.

```go
type Options struct {
    ReadOnly      bool
    DisableDelete bool
}
```

Rules:

- `ReadOnly: true` hides all create/update/delete tools
- `DisableDelete: true` hides delete tools while leaving create/update intact
- This is used by clients that need a limited or safer MCP surface

## Tool naming and argument conventions

The project uses a consistent naming pattern across all tools:

- `list_*` for listing collections
- `get_*` for retrieval by ID
- `create_*` for creation
- `update_*` for mutation
- `delete_*` for removal
- `*_status` when a status-specific mutation exists

Common argument conventions:

- project-scoped entities use `projectID`
- board-scoped or board-column-scoped operations include `boardID`
- item IDs are usually called `id` or `itemID`, depending on the entity
- UUIDs are validated via the shared MCP helpers in `internal/adapter/in/mcp/util`

Always prefer the same naming pattern used by the rest of the package so the MCP contract stays stable for AI clients and frontend integrations.

## UUID validation and errors

All request handlers must validate IDs before calling use cases. The shared utility layer provides helpers such as:

- `ExtractUUID(...)`
- `RequireString(...)`
- `ValidateQuery(...)`
- `ValidateProjectQuery(...)`

This is important because invalid UUIDs should return a proper MCP tool error rather than bubbling an unhandled low-level error.

The project convention is to return an MCP error result with a human-readable message, e.g.:

```go
return mcp.NewToolResultError(err.Error()), nil
```

## Read-only vs write annotations

Tool annotations are used to describe their behavior to MCP clients:

- `util.ReadOnlyAnnotations`
- `util.WriteAnnotations`
- `util.DeleteAnnotations`

This matters because the MCP server may expose different tool sets depending on runtime mode and access policy. Tools should be tagged correctly to reflect the allowed actions.

## Markdown references for linked items in cards

Card descriptions support Markdown mentions for linked items. The frontend pattern is:

```text
@[Item Title](item:TYPE:UUID)
```

Example:

```text
@[Login API](item:SNIPPET:123e4567-e89b-12d3-a456-426614174000)
```

Important rule:

- The item must first be added to the card's `linkedItems` collection
- After that, the same `itemType` and `itemId` may be used in the Markdown mention
- The mention syntax is only valid for items that are already linked to the card

This is enforced in the MCP description strings for `create_card` and `update_card`, and is part of the tool contract used by clients building card content.

## Tag system conventions

Tags are cross-cutting and follow the same architecture described in `docs/architecture/tags.md`.

Some rules relevant to MCP:

- `ItemType` values are uppercase domain constants such as `SNIPPET`, `NOTE`, `PROBLEM`, `LINK`, `CREDENTIAL`, `BOARD`, and `CARD`
- Tag association tools must accept the item type and normalize values consistently before the use case layer
- Search support exists in MCP for tag lookup by name using `search_tags_by_name`

## Entity-specific conventions

Each tool file usually follows these patterns:

- `list_*`: returns a list or paginated list of entity views
- `get_*`: returns a single entity by identifier
- `create_*`: accepts DTO fields matching `internal/dto/*Command`
- `update_*`: uses the update command DTO and returns the updated entity
- `delete_*`: returns a simple success text like `"Project Deleted!"`

When adding a new tool, match the naming, argument names, and result style used by neighboring tool files.

## What to avoid

When working on MCP code, avoid these common mistakes:

- creating business logic in the MCP adapter layer
- registering tools without matching the use case contract
- using inconsistent argument names (`project_id` vs `projectID`)
- forgetting `ReadOnly` and `DisableDelete` gates
- returning raw errors instead of MCP `ToolResultError`
- exposing delete tools in read-only mode
- documenting card mention syntax without also documenting the `linkedItems` requirement

## Summary

The MCP layer should stay thin and consistent. The primary responsibilities are:

- validate inputs
- map to DTO commands
- call the use case
- return MCP-native tool result payloads
- follow the same contract conventions across all tools

If a newcomer knows the patterns in this file and the tag architecture doc, they can understand the expected MCP behavior without reading every tool implementation individually.
