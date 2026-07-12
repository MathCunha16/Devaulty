package com.devaulty.backend.adapter.in.web.problem;

import com.devaulty.backend.adapter.in.web.problem.dto.*;
import com.devaulty.backend.adapter.in.web.util.UriLocationBuilderHelper;
import com.devaulty.backend.application.port.in.problem.*;
import com.devaulty.backend.domain.model.Problem;
import jakarta.validation.Valid;
import org.springframework.data.domain.Page;
import org.springframework.http.ResponseEntity;
import org.springframework.web.bind.annotation.*;

import java.net.URI;
import java.util.UUID;

@RestController
@RequestMapping("/api/v1/projects/{projectId}/problems")
public class ProblemController implements ProblemApi{

    private final CreateProblemUseCase createProblemUseCase;
    private final GetAllProblemsByProjectUseCase getAllProblemsByProjectUseCase;
    private final GetProblemByIdUseCase getProblemByIdUseCase;
    private final UpdateProblemUseCase updateProblemUseCase;
    private final UpdateProblemStatusUseCase updateProblemStatusUseCase;
    private final DeleteProblemUseCase deleteProblemUseCase;
    private final ProblemWebMapper mapper;
    private final UriLocationBuilderHelper uriLocationBuilderHelper;

    public ProblemController(CreateProblemUseCase createProblemUseCase, GetAllProblemsByProjectUseCase getAllProblemsByProjectUseCase, GetProblemByIdUseCase getProblemByIdUseCase, UpdateProblemUseCase updateProblemUseCase, UpdateProblemStatusUseCase updateProblemStatusUseCase, DeleteProblemUseCase deleteProblemUseCase, ProblemWebMapper mapper, UriLocationBuilderHelper uriLocationBuilderHelper) {
        this.createProblemUseCase = createProblemUseCase;
        this.getAllProblemsByProjectUseCase = getAllProblemsByProjectUseCase;
        this.getProblemByIdUseCase = getProblemByIdUseCase;
        this.updateProblemUseCase = updateProblemUseCase;
        this.updateProblemStatusUseCase = updateProblemStatusUseCase;
        this.deleteProblemUseCase = deleteProblemUseCase;
        this.mapper = mapper;
        this.uriLocationBuilderHelper = uriLocationBuilderHelper;
    }

    @Override
    @PostMapping
    public ResponseEntity<ProblemViewResponse> create(
            @PathVariable UUID projectId,
            @RequestBody @Valid CreateProblemRequest request)
    {
        CreateProblemCommand command = mapper.toCreateProblemCommand(request, projectId);
        Problem problem = createProblemUseCase.execute(command);
        URI location = uriLocationBuilderHelper.buildLocationUri(problem.getId());
        return ResponseEntity.created(location).body(mapper.toViewResponse(problem));
    }

    @Override
    @GetMapping
    public ResponseEntity<Page<ProblemSummaryResponse>> getAllByProject(
            @PathVariable UUID projectId,
            @RequestParam(defaultValue = "0") int page,
            @RequestParam(defaultValue = "10") int size
    ){
        Page<Problem> problems = getAllProblemsByProjectUseCase.execute(projectId, page, size);
        return ResponseEntity.ok(problems.map(mapper::toSummaryResponse));
    }

    @Override
    @GetMapping("/{problemId}")
    public ResponseEntity<ProblemViewResponse> getById(
            @PathVariable UUID projectId,
            @PathVariable UUID problemId
    ){
        Problem problem = getProblemByIdUseCase.execute(projectId, problemId);
        return ResponseEntity.ok(mapper.toViewResponse(problem));
    }

    @Override
    @PatchMapping("/{problemId}")
    public ResponseEntity<ProblemViewResponse> update(
            @PathVariable UUID projectId,
            @PathVariable UUID problemId,
            @RequestBody @Valid UpdateProblemRequest request
    ){
        UpdateProblemCommand command = mapper.toUpdateProblemCommand(request, projectId, problemId);
        Problem problem = updateProblemUseCase.execute(command);
        return ResponseEntity.ok(mapper.toViewResponse(problem));
    }

    @Override
    @PatchMapping("/{problemId}/status")
    public ResponseEntity<ProblemViewResponse> updateStatus(
            @PathVariable UUID projectId,
            @PathVariable UUID problemId,
            @RequestBody @Valid UpdateProblemStatusRequest request
    ){
        UpdateProblemStatusCommand command = mapper.toUpdateProblemStatusCommand(request, projectId, problemId);
        Problem problem = updateProblemStatusUseCase.execute(command);
        return ResponseEntity.ok(mapper.toViewResponse(problem));
    }

    @Override
    @DeleteMapping("/{problemId}")
    public ResponseEntity<Void> delete(
            @PathVariable UUID projectId,
            @PathVariable UUID problemId
    ){
        deleteProblemUseCase.execute(projectId, problemId);
        return ResponseEntity.noContent().build();
    }
}
