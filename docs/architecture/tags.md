# Tag System Architecture

This document explains how tags work across Devaulty, why they're built the
way they are, and what to do when adding tag support to a new entity.

## The problem tags solve

Tags are a **cross-cutting concern**: a single tag (`"docker"`, `"urgent"`)
can be attached to items of completely different types — a snippet, a
problem, a link, a credential, a note. Modeling this the "obvious" way (a
`List<Tag> tags` field on every domain entity) would mean every domain model,
every use case, and every persistence adapter would need to know about tags.
That's the wrong trade-off: it turns an add-on feature into a dependency
baked into everything else.

Instead, tags are kept **out of the domain models entirely** and are
attached/read as a separate step, orchestrated at the web layer.

## The core design decision

> Domain models (`Snippet`, `Problem`, `Link`, `Credential`, `Note`, ...)
> never carry a `tags` field. Tags are looked up on demand, only when
> building an HTTP response, and merged in at that point.

This keeps every existing use case (`CreateSnippetImpl`, `UpdateProblemImpl`,
etc.) completely untouched. Tag support is additive, not invasive.

## The moving parts

### 1. Storage: polymorphic association, not a domain object

`item_tags` is a join table with a composite key (`tag_id`, `item_type`,
`item_id`) and **no real foreign key to the tagged item** — `item_type` +
`item_id` is just a loose pointer, because a single join table can't have a
proper FK to five different tables.

Consequence: the database can't cascade-delete `item_tags` rows when a
snippet/problem/link/credential/note is deleted. Every `Delete*Impl` for a
taggable entity **must** explicitly clean up its tag associations
(`ItemTagRepositoryPort#removeAllTagsFromItem`) after deleting the item
itself. This is the one place ownership-by-database breaks down and the
application has to enforce it manually.

`ItemTagEntity` is a JPA implementation detail — nothing outside
`adapter/out/persistence/tag` ever sees it. The rest of the app works with
the domain `Tag` type only.

### 2. Reading tags: two use cases, chosen by how many items you're showing

- `GetTagsForItemUseCase` — one item, one call. Used on `getById` endpoints.
- `GetTagsForItemsUseCase` — many items, one batched call, returns
  `Map<UUID, List<Tag>>` keyed by item id. Used on `getAll`/paginated
  endpoints.

**Never call the singular use case in a loop over a list of items** — that's
an N+1 query bug. If you're building a list response, use the plural,
batched version.

### 3. Writing tags: a separate controller, not part of create/update

Associating or removing a tag is its own endpoint, decoupled from the
entity's own create/update flow:

```
PUT    /api/v1/projects/{projectId}/items/{itemType}/{itemId}/tags/{tagId}
DELETE /api/v1/projects/{projectId}/items/{itemType}/{itemId}/tags/{tagId}
```

`PUT` (not `POST`) because attaching a tag is idempotent — attaching the
same tag twice should be a no-op, not an error or a duplicate.

A newly created item never has tags yet, so `create` endpoints skip the tag
lookup entirely and pass an empty list to the mapper.

### 4. Merging tags into the response: mapper takes an extra parameter

DTOs for taggable entities carry a `List<TagSummaryResponse> tags` field.
The relevant `*WebMapper` gets a second method parameter for the tag list,
and delegates the `Tag → TagSummaryResponse` conversion to `TagWebMapper` via
MapStruct's `uses`:

```java
@Mapper(componentModel = "spring", uses = TagWebMapper.class)
public interface SnippetWebMapper {

    @Mapping(target = "tags", source = "tags")
    SnippetViewResponse toViewResponse(Snippet snippet, List<Tag> tags);
}
```

The controller is responsible for fetching the tags and passing them in —
mappers never reach into a repository themselves.

## Full flow example: `GET /projects/{projectId}/snippets`

```
SnippetController.getAllByProject
  → getAllSnippetsByProjectUseCase.execute(...)      // fetch the page of snippets
  → collect snippet IDs from that page
  → getTagsForItemsUseCase.execute("snippet", projectId, ids)  // ONE batched query
  → snippets.map(s -> webMapper.toSummaryResponse(s, tagsByItem.get(s.getId())))
  → return
```

## Checklist: adding tag support to a new taggable entity

- [ ] Domain model stays untouched — no `tags` field.
- [ ] Add `List<TagSummaryResponse> tags` to its `*ViewResponse` and
      `*SummaryResponse` DTOs.
- [ ] `*WebMapper`: add `uses = TagWebMapper.class` and a second parameter
      (`List<Tag> tags`) to the `toViewResponse`/`toSummaryResponse` methods.
- [ ] Controller: inject `GetTagsForItemUseCase` and
      `GetTagsForItemsUseCase`.
  - `getById` → call the singular use case.
  - `getAll` → collect IDs, call the plural/batched use case once.
  - `create` → pass `List.of()`, no lookup needed.
  - `update` → call the singular use case (the item already existed, it may
    already have tags).
- [ ] `Delete*Impl`: after deleting the item, call
      `itemTagRepositoryPort.removeAllTagsFromItem(itemType, id)` to avoid
      orphaned `item_tags` rows (see storage section above — the database
      won't do this for you).
- [ ] Use the entity's lowercase name as the `item_type` string consistently
      (e.g. `"snippet"`, `"credential"`) — this string is the only thing
      tying `item_tags` rows back to the right table, so a typo here
      silently breaks the association. Prefer a `private static final
      String ITEM_TYPE = "..."` constant in the controller over inlining
      the literal at each call site.

## Why not a `Taggable` interface on domain models?

It was considered and rejected. Giving every taggable domain model a
`List<Tag> tags` field (via a shared interface) would mean every existing
`Create*Impl`/`Update*Impl` would need to manage that list, every
persistence adapter would need to populate it on every read, and the
"cheap" read-only nature of tag lookups (batch them only when building a
response) would be lost. The chosen approach keeps tag support entirely
additive: it touches the web layer only, and existing application/domain
code for every entity remains exactly as it was before tags existed.
