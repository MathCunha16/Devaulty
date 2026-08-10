# Tag System Architecture (Go Backend)

This document explains how tags work across Devaulty's Go backend (`backend-go`), why they're built the way they are, and what to do when adding tag support to a new entity.

## The problem tags solve

Tags are a **cross-cutting concern**: a single tag (`"docker"`, `"urgent"`) can be attached to items of completely different types — a snippet, a problem, a link, a credential, a note. Modeling this with a direct field in every domain model would pollute every core domain entity, use case, and repository interface with tag logic.

Instead, domain entities (`model.Snippet`, `model.Link`, `model.Problem`, etc.) remain clean and tag-agnostic in `internal/domain/model`. Tags are stored via a polymorphic association table (`item_tags`) and attached to response DTOs (`internal/dto`) inside the application use cases.

## The core design decision

> Core domain models (`model.Snippet`, `model.Problem`, `model.Link`, etc.) never carry a `Tags` field.
> Tags are fetched on demand by use cases using an injected `port.ItemTagRepository` and populated into response DTOs (`dto.SnippetView`, `dto.LinkView`, `dto.ProblemView`, `dto.ProblemSummary`) before returning to the web layer.

This keeps domain rules and database models simple while centralizing application DTOs in package `internal/dto`.

## The moving parts

### 1. Storage: polymorphic association (`item_tags`)

`item_tags` is a join table with a composite primary key (`tag_id`, `item_type`, `item_id`) and no foreign key constraint on `item_id` (since `item_id` can refer to multiple tables):

```sql
CREATE TABLE item_tags (
    tag_id VARCHAR(36) NOT NULL,
    item_type VARCHAR(20) NOT NULL,
    item_id VARCHAR(36) NOT NULL,
    PRIMARY KEY (tag_id, item_type, item_id),
    FOREIGN KEY (tag_id) REFERENCES tags(id) ON DELETE CASCADE
);
```

**Cascading Cleanup:** Because SQLite cannot enforce FK cascades on `item_id` across different tables, every `Delete` method for a taggable entity in `internal/usecase/` (e.g. `SnippetUseCase.Delete`, `LinkUseCase.Delete`, `ProblemUseCase.Delete`) explicitly calls `itemTagRepo.RemoveAllTagsFromItem(ctx, itemType, itemID)`.

### 2. DTO Centralization (`internal/dto/`)

All request commands and response views live in package `devaulty-backend/internal/dto`:

- **`dto.TagSummary`**: Lightweight DTO (`ID`, `Name`, `Color`) embedded in item response views.
- **`dto.TagView`**: Full DTO (`ID`, `ProjectID`, `Name`, `Color`, `CreatedAt`, `UpdatedAt`) used for tag management responses.
- **Item Views**: `SnippetView`, `LinkView`, `ProblemView`, and `ProblemSummary` contain `Tags []dto.TagSummary`.

### 3. Reading tags: two query patterns in Use Cases

Use cases inject `port.ItemTagRepository` and fetch tags when building DTO views:

- **Single Item (`GetByID`, `Update`)**: Calls `itemTagRepo.FindTagsForItem(ctx, itemType, projectID, itemID)`.
- **Paginated List (`GetAllByProjectID`)**: Collects item IDs from the page content and calls `itemTagRepo.FindTagsForItems(ctx, itemType, projectID, itemIDs)` in a **single batched SQL query** returning `map[uuid.UUID][]model.Tag`.

> **Rule:** Never execute single-item tag queries in a loop over paginated results to avoid N+1 query performance bottlenecks.

### 4. Writing tags: dedicated `ItemTagHandler` endpoints

Tag association and disassociation are decoupled from item creation and updates:

```text
PUT    /api/v1/projects/:project_id/items/:item_type/:item_id/tags/:tag_id
DELETE /api/v1/projects/:project_id/items/:item_type/:item_id/tags/:tag_id
```

- `PUT` associates a tag to an item (`204 No Content`). It is idempotent.
- `DELETE` disassociates a tag from an item (`204 No Content`).
- Query path param `:item_type` is parsed in uppercase (`model.ItemType(strings.ToUpper(c.Param("item_type")))`) to match registered domain constants (`SNIPPET`, `LINK`, `PROBLEM`).

### 5. Tag Management: `TagHandler` endpoints

Full CRUD and search for project tags:

```text
POST   /api/v1/projects/:project_id/tags          (Create tag)
GET    /api/v1/projects/:project_id/tags          (Get all tags)
GET    /api/v1/projects/:project_id/tags/search   (Search tags by name prefix/substring)
GET    /api/v1/projects/:project_id/tags/:tag_id  (Get tag by ID)
PATCH  /api/v1/projects/:project_id/tags/:tag_id  (Update tag)
DELETE /api/v1/projects/:project_id/tags/:tag_id  (Delete tag)
```

Tag names are unique per project (`strings.TrimSpace(cmd.Name)`), returning `400 Bad Request` if a tag with the same name already exists.

## Full flow example: `GET /projects/{project_id}/snippets`

```text
SnippetHandler.GetAll
  → SnippetUseCase.GetAllByProjectID(ctx, projectID, page, size)
    → snippetRepo.FindAllByProjectID(...)             // Fetch page of snippets
    → collect snippet IDs
    → itemTagRepo.FindTagsForItems(...)               // ONE batched query for tags
    → assemble dto.SnippetView slice with []dto.TagSummary
  → returns model.Page[dto.SnippetView] as JSON
```

## Checklist: adding tag support to a new entity

- [ ] Core domain model remains clean (no `Tags` field in `internal/domain/model`).
- [ ] Define its DTO views in `internal/dto/` and include `Tags []dto.TagSummary`.
- [ ] Inject `port.ItemTagRepository` into its UseCase constructor.
- [ ] In UseCase `GetByID` & `Update`: call `itemTagRepo.FindTagsForItem` and attach tags.
- [ ] In UseCase `GetAllByProjectID`: collect IDs, call `itemTagRepo.FindTagsForItems` (batched), and attach tags.
- [ ] In UseCase `Create`: pass empty tag slice `[]dto.TagSummary{}`.
- [ ] In UseCase `Delete`: call `itemTagRepo.RemoveAllTagsFromItem(ctx, itemType, id)` after deleting the item.
- [ ] Register repository in `NewItemTagUseCase` and add corresponding `model.ItemType` enum constant (e.g. `SNIPPET`, `LINK`, `PROBLEM`, `NOTE`).
