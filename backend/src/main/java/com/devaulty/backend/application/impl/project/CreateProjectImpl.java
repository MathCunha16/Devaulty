package com.devaulty.backend.application.impl.project;

import com.devaulty.backend.application.port.in.project.CreateProjectCommand;
import com.devaulty.backend.application.port.in.project.CreateProjectUseCase;
import com.devaulty.backend.application.port.out.persistence.ProjectRepositoryPort;
import com.devaulty.backend.domain.model.Project;
import org.springframework.transaction.annotation.Transactional;

import java.time.LocalDateTime;
import java.util.UUID;

public class CreateProjectImpl implements CreateProjectUseCase {

    private final ProjectRepositoryPort projectRepository;

    public CreateProjectImpl(ProjectRepositoryPort projectRepository) {
        this.projectRepository = projectRepository;
    }

    @Override
    @Transactional
    public Project execute(CreateProjectCommand command) {
        Project project = new Project();
        project.setId(UUID.randomUUID());
        project.setName(command.name());
        project.setDescription(command.description());
        project.setIcon(command.icon());
        project.setColor(command.color());
        project.setArchived(false);
        project.setCreatedAt(LocalDateTime.now());
        return projectRepository.save(project);
    }
}
