package com.devaulty.backend.adapter.in.web.snippet;

import com.devaulty.backend.adapter.in.web.snippet.dto.CreateSnippetRequest;
import com.devaulty.backend.adapter.in.web.snippet.dto.SnippetViewResponse;
import com.devaulty.backend.adapter.in.web.snippet.dto.UpdateSnippetRequest;
import com.devaulty.backend.adapter.in.web.util.UriLocationBuilderHelper;
import com.devaulty.backend.application.port.in.snippet.*;
import com.devaulty.backend.domain.model.Snippet;
import jakarta.validation.Valid;
import org.springframework.data.domain.Page;
import org.springframework.http.ResponseEntity;
import org.springframework.web.bind.annotation.*;

import java.net.URI;
import java.util.UUID;

@RestController
@RequestMapping("/api/v1/projects/{projectId}/snippets")
public class SnippetController implements SnippetApi{

    private final CreateSnippetUseCase createSnippetUseCase;
    private final GetAllSnippetsByProjectUseCase getAllSnippetsByProjectUseCase;
    private final GetSnippetByIdUseCase getSnippetByIdUseCase;
    private final UpdateSnippetUseCase updateSnippetUseCase;
    private final DeleteSnippetUseCase deleteSnippetUseCase;
    private final SnippetWebMapper webMapper;
    private final UriLocationBuilderHelper uriLocationBuilderHelper;

    public SnippetController(CreateSnippetUseCase createSnippetUseCase, GetAllSnippetsByProjectUseCase getAllSnippetsByProjectUseCase, GetSnippetByIdUseCase getSnippetByIdUseCase, UpdateSnippetUseCase updateSnippetUseCase, DeleteSnippetUseCase deleteSnippetUseCase, SnippetWebMapper webMapper, UriLocationBuilderHelper uriLocationBuilderHelper) {
        this.createSnippetUseCase = createSnippetUseCase;
        this.getAllSnippetsByProjectUseCase = getAllSnippetsByProjectUseCase;
        this.getSnippetByIdUseCase = getSnippetByIdUseCase;
        this.updateSnippetUseCase = updateSnippetUseCase;
        this.deleteSnippetUseCase = deleteSnippetUseCase;
        this.webMapper = webMapper;
        this.uriLocationBuilderHelper = uriLocationBuilderHelper;
    }

    @Override
    @PostMapping
    public ResponseEntity<SnippetViewResponse> create(
            @PathVariable UUID projectId,
            @RequestBody @Valid CreateSnippetRequest request

    ) {
        CreateSnippetCommand command = webMapper.toCreateSnippetCommand(request, projectId);
        Snippet snippet = createSnippetUseCase.execute(command);
        URI location = uriLocationBuilderHelper.buildLocationUri(snippet.getId());
        return ResponseEntity.created(location).body(webMapper.toViewResponse(snippet));
    }

    @Override
    @GetMapping
    public ResponseEntity<Page<SnippetViewResponse>> getAllByProject (
            @PathVariable UUID projectId,
            @RequestParam(defaultValue = "0") int page,
            @RequestParam(defaultValue = "10") int size
    ){
        Page<Snippet> snippets = getAllSnippetsByProjectUseCase.execute(projectId, page, size);
        return ResponseEntity.ok(snippets.map(webMapper::toViewResponse));
    }

    @Override
    @GetMapping("/{snippetId}")
    public ResponseEntity<SnippetViewResponse> getById(
            @PathVariable UUID projectId,
            @PathVariable UUID snippetId
    ){
        Snippet snippet = getSnippetByIdUseCase.execute(projectId, snippetId);
        return ResponseEntity.ok(webMapper.toViewResponse(snippet));
    }

    @Override
    @PatchMapping("/{snippetId}")
    public ResponseEntity<SnippetViewResponse> update(
            @PathVariable UUID projectId,
            @PathVariable UUID snippetId,
            @RequestBody @Valid UpdateSnippetRequest request
    ){
        UpdateSnippetCommand command = webMapper.toUpdateSnippetCommand(request, projectId, snippetId);
        Snippet snippet = updateSnippetUseCase.execute(command);
        return ResponseEntity.ok(webMapper.toViewResponse(snippet));
    }

    @Override
    @DeleteMapping("/{snippetId}")
    public ResponseEntity<Void> delete(
            @PathVariable UUID projectId,
            @PathVariable UUID snippetId
    ){
        deleteSnippetUseCase.execute(projectId, snippetId);
        return ResponseEntity.noContent().build();
    }
}
