package com.devaulty.backend.adapter.in.web.snippet;

import com.devaulty.backend.adapter.in.web.snippet.dto.CreateSnippetRequest;
import com.devaulty.backend.adapter.in.web.snippet.dto.SnippetViewResponse;
import com.devaulty.backend.adapter.in.web.util.UriLocationBuilderHelper;
import com.devaulty.backend.application.port.in.snippet.CreateSnippetCommand;
import com.devaulty.backend.application.port.in.snippet.CreateSnippetUseCase;
import com.devaulty.backend.application.port.in.snippet.GetAllSnippetsByProjectUseCase;
import com.devaulty.backend.application.port.in.snippet.GetSnippetByIdUseCase;
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
    private final SnippetWebMapper webMapper;
    private final UriLocationBuilderHelper uriLocationBuilderHelper;

    public SnippetController(CreateSnippetUseCase createSnippetUseCase, GetAllSnippetsByProjectUseCase getAllSnippetsByProjectUseCase, GetSnippetByIdUseCase getSnippetByIdUseCase, SnippetWebMapper webMapper, UriLocationBuilderHelper uriLocationBuilderHelper) {
        this.createSnippetUseCase = createSnippetUseCase;
        this.getAllSnippetsByProjectUseCase = getAllSnippetsByProjectUseCase;
        this.getSnippetByIdUseCase = getSnippetByIdUseCase;
        this.webMapper = webMapper;
        this.uriLocationBuilderHelper = uriLocationBuilderHelper;
    }

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

    @GetMapping
    public ResponseEntity<Page<SnippetViewResponse>> getAllByProject (
            @PathVariable UUID projectId,
            @RequestParam(defaultValue = "0") int page,
            @RequestParam(defaultValue = "10") int size
    ){
        Page<Snippet> snippets = getAllSnippetsByProjectUseCase.execute(projectId, page, size);
        return ResponseEntity.ok(snippets.map(webMapper::toViewResponse));
    }

    @GetMapping("/{snippetId}")
    public ResponseEntity<SnippetViewResponse> getById(
            @PathVariable UUID projectId,
            @PathVariable UUID snippetId
    ){
        Snippet snippet = getSnippetByIdUseCase.execute(projectId, snippetId);
        return ResponseEntity.ok(webMapper.toViewResponse(snippet));
    }
}
