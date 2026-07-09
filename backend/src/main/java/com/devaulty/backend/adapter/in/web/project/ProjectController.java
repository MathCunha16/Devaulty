package com.devaulty.backend.adapter.in.web.project;

import com.devaulty.backend.adapter.in.web.project.dto.CreateProjectRequest;
import com.devaulty.backend.adapter.in.web.project.dto.ProjectViewResponse;
import com.devaulty.backend.adapter.in.web.util.UriLocationBuilderHelper;
import com.devaulty.backend.application.port.in.project.*;
import com.devaulty.backend.domain.model.Project;
import jakarta.validation.Valid;
import org.springframework.data.domain.Page;
import org.springframework.http.ResponseEntity;
import org.springframework.web.bind.annotation.*;

import java.net.URI;
import java.util.UUID;

@RestController
@RequestMapping("/api/v1/projects")
public class ProjectController implements ProjectApi {

    private final CreateProjectUseCase createProjectUseCase;
    private final GetAllProjectsUseCase getAllProjectsUseCase;
    private final GetProjectByIdUseCase getProjectByIdUseCase;
    private final UpdateProjectUseCase updateProjectUseCase;
    private final ArchiveProjectUseCase archiveProjectUseCase;
    private final UnarchiveProjectUseCase unarchiveProjectUseCase;
    private final DeleteProjectUseCase deleteProjectUseCase;
    private final ProjectWebMapper webMapper;
    private final UriLocationBuilderHelper uriLocationBuilderHelper;

    public ProjectController(CreateProjectUseCase createProjectUseCase, GetAllProjectsUseCase getAllProjectsUseCase, GetProjectByIdUseCase getProjectByIdUseCase, UpdateProjectUseCase updateProjectUseCase, ArchiveProjectUseCase archiveProjectUseCase, UnarchiveProjectUseCase unarchiveProjectUseCase, DeleteProjectUseCase deleteProjectUseCase, ProjectWebMapper webMapper, UriLocationBuilderHelper uriLocationBuilderHelper) {
        this.createProjectUseCase = createProjectUseCase;
        this.getAllProjectsUseCase = getAllProjectsUseCase;
        this.getProjectByIdUseCase = getProjectByIdUseCase;
        this.updateProjectUseCase = updateProjectUseCase;
        this.archiveProjectUseCase = archiveProjectUseCase;
        this.unarchiveProjectUseCase = unarchiveProjectUseCase;
        this.deleteProjectUseCase = deleteProjectUseCase;
        this.webMapper = webMapper;
        this.uriLocationBuilderHelper = uriLocationBuilderHelper;
    }

    @Override
    @PostMapping
    public ResponseEntity<ProjectViewResponse> create(@RequestBody @Valid CreateProjectRequest request){
        CreateProjectCommand command = webMapper.toCreateProjectCommand(request);
        Project project = createProjectUseCase.execute(command);
        URI location = uriLocationBuilderHelper.buildLocationUri(project.getId());
        return ResponseEntity.created(location).body(webMapper.toViewResponse(project));
    }

    @Override
    @GetMapping
    public ResponseEntity<Page<ProjectViewResponse>> getAll(
            @RequestParam(defaultValue = "0") int page,
            @RequestParam(defaultValue = "10") int size)
    {
        Page<Project> projects = getAllProjectsUseCase.execute(page, size);
        return ResponseEntity.ok(projects.map(webMapper::toViewResponse));
    }

    @Override
    @GetMapping("/{id}")
    public ResponseEntity<ProjectViewResponse> getById(@PathVariable UUID id){
        Project project = getProjectByIdUseCase.execute(id);
        return ResponseEntity.ok(webMapper.toViewResponse(project));
    }

    @PatchMapping("/{id}")
    @Override
    public ResponseEntity<ProjectViewResponse> update(@PathVariable UUID id,
                                                      @RequestBody @Valid CreateProjectRequest request)
    {
        CreateProjectCommand command = webMapper.toCreateProjectCommand(request);
        Project project = updateProjectUseCase.execute(id, command);
        return ResponseEntity.ok(webMapper.toViewResponse(project));
    }

    @PatchMapping("/{id}/archive")
    @Override
    public ResponseEntity<Void> archive(@PathVariable UUID id){
        archiveProjectUseCase.execute(id);
        return ResponseEntity.noContent().build();
    }

    @PatchMapping("/{id}/unarchive")
    @Override
    public ResponseEntity<Void> unarchive(@PathVariable UUID id){
        unarchiveProjectUseCase.execute(id);
        return ResponseEntity.noContent().build();
    }

    @DeleteMapping("/{id}")
    @Override
    public ResponseEntity<Void> delete(@PathVariable UUID id){
        deleteProjectUseCase.execute(id);
        return ResponseEntity.noContent().build();
    }
}