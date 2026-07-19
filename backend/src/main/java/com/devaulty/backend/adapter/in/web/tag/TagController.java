package com.devaulty.backend.adapter.in.web.tag;

import com.devaulty.backend.adapter.in.web.tag.dto.CreateTagRequest;
import com.devaulty.backend.adapter.in.web.tag.dto.TagViewResponse;
import com.devaulty.backend.adapter.in.web.tag.dto.UpdateTagRequest;
import com.devaulty.backend.adapter.in.web.util.UriLocationBuilderHelper;
import com.devaulty.backend.application.port.in.tag.*;
import jakarta.validation.Valid;
import org.springframework.http.ResponseEntity;
import org.springframework.web.bind.annotation.*;

import java.net.URI;
import java.util.List;
import java.util.UUID;

@RestController
@RequestMapping("/ap1/v1/projects/{projectId}/tags")
public class TagController implements TagApi {

    private final CreateTagUseCase createTagUseCase;
    private final GetAllTagsByProjectUseCase getAllTagsByProjectUseCase;
    private final GetTagByIdUseCase getTagByIdUseCase;
    private final SearchTagByNameUseCase searchTagByNameUseCase;
    private final UpdateTagUseCase updateTagUseCase;
    private final DeleteTagUseCase deleteTagUseCase;
    private final TagWebMapper mapper;
    private final UriLocationBuilderHelper uriLocationBuilderHelper;

    public TagController(CreateTagUseCase createTagUseCase, GetAllTagsByProjectUseCase getAllTagsByProjectUseCase, GetTagByIdUseCase getTagByIdUseCase, SearchTagByNameUseCase searchTagByNameUseCase, UpdateTagUseCase updateTagUseCase, DeleteTagUseCase deleteTagUseCase, TagWebMapper mapper, UriLocationBuilderHelper uriLocationBuilderHelper) {
        this.createTagUseCase = createTagUseCase;
        this.getAllTagsByProjectUseCase = getAllTagsByProjectUseCase;
        this.getTagByIdUseCase = getTagByIdUseCase;
        this.searchTagByNameUseCase = searchTagByNameUseCase;
        this.updateTagUseCase = updateTagUseCase;
        this.deleteTagUseCase = deleteTagUseCase;
        this.mapper = mapper;
        this.uriLocationBuilderHelper = uriLocationBuilderHelper;
    }

    @Override
    @PostMapping
    public ResponseEntity<TagViewResponse> create(
            @PathVariable UUID projectId,
            @RequestBody @Valid CreateTagRequest request)
    {
        CreateTagCommand command = mapper.toCreateTagCommand(request, projectId);
        TagViewResponse response = mapper.toViewResponse(createTagUseCase.execute(command));
        URI location = uriLocationBuilderHelper.buildLocationUri(response.id());
        return ResponseEntity.created(location).body(response);
    }

    @Override
    @GetMapping
    public ResponseEntity<List<TagViewResponse>> getAll(
            @PathVariable UUID projectId
    ){
        List<TagViewResponse> responses = getAllTagsByProjectUseCase.execute(projectId)
                .stream()
                .map(mapper::toViewResponse)
                .toList();
        return ResponseEntity.ok(responses);
    }

    @Override
    @GetMapping("/{tagId}")
    public ResponseEntity<TagViewResponse> getById(
            @PathVariable UUID projectId,
            @PathVariable UUID tagId
    ){
        TagViewResponse response = mapper.toViewResponse(getTagByIdUseCase.execute(projectId, tagId));
        return ResponseEntity.ok(response);
    }

    @Override
    @GetMapping("/search")
    public ResponseEntity<List<TagViewResponse>> searchByName(
            @PathVariable UUID projectId,
            @RequestParam String name)
    {
        List<TagViewResponse> responses = searchTagByNameUseCase.execute(projectId, name)
                .stream()
                .map(mapper::toViewResponse)
                .toList();
        return ResponseEntity.ok(responses);
    }

    @Override
    @PatchMapping("/{tagId}")
    public ResponseEntity<TagViewResponse> update(
            @PathVariable UUID projectId,
            @PathVariable UUID tagId,
            @RequestBody @Valid UpdateTagRequest request
    ){
        UpdateTagCommand command = mapper.toUpdateTagCommand(request, projectId, tagId);
        TagViewResponse response = mapper.toViewResponse(updateTagUseCase.execute(command));
        return ResponseEntity.ok(response);
    }

    @Override
    @DeleteMapping("/{tagId}")
    public ResponseEntity<Void> delete(
            @PathVariable UUID projectId,
            @PathVariable UUID tagId
    ){
        deleteTagUseCase.execute(projectId, tagId);
        return ResponseEntity.noContent().build();
    }
}
