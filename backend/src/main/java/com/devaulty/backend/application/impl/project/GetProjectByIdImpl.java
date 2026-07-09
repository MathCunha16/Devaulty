package com.devaulty.backend.application.impl.project;

import com.devaulty.backend.application.exception.ResourceNotFoundException;
import com.devaulty.backend.application.port.in.project.GetProjectByIdUseCase;
import com.devaulty.backend.application.port.out.persistence.ProjectRepositoryPort;
import com.devaulty.backend.domain.model.Project;
import org.springframework.transaction.annotation.Transactional;

import java.util.UUID;

public class GetProjectByIdImpl implements GetProjectByIdUseCase {

    private final ProjectRepositoryPort projectRepository;

    public GetProjectByIdImpl(ProjectRepositoryPort projectRepository) {
        this.projectRepository = projectRepository;
    }

    @Override
    @Transactional(readOnly = true)
    public Project execute(UUID id) {
        return projectRepository.findById(id).
                orElseThrow(() -> new ResourceNotFoundException("Project", id));
    }
}
