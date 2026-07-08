package com.devaulty.backend.application.impl.project;

import com.devaulty.backend.application.exception.ResourceNotFoundException;
import com.devaulty.backend.application.port.in.project.CreateProjectCommand;
import com.devaulty.backend.application.port.in.project.UpdateProjectUseCase;
import com.devaulty.backend.application.port.out.persistence.ProjectRepositoryPort;
import com.devaulty.backend.domain.model.Project;
import org.springframework.transaction.annotation.Transactional;

import java.time.LocalDateTime;
import java.util.UUID;

public class UpdateProjectImpl implements UpdateProjectUseCase {

    private final ProjectRepositoryPort projectRepositoryPort;

    public UpdateProjectImpl(ProjectRepositoryPort projectRepositoryPort) {
        this.projectRepositoryPort = projectRepositoryPort;
    }

    @Override
    @Transactional
    public Project execute(UUID id, CreateProjectCommand command) {

        Project project = projectRepositoryPort.findById(id)
                .orElseThrow(() -> new ResourceNotFoundException("Project", id));

        if(command.name() != null) project.setName(command.name());
        if(command.description() != null) project.setDescription(command.description());
        if(command.icon() != null) project.setIcon(command.icon());
        if(command.color() != null) project.setColor(command.color());
        project.setUpdatedAt(LocalDateTime.now());

        return projectRepositoryPort.save(project);
    }
}
