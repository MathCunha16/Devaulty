package com.devaulty.backend.application.impl.project;

import com.devaulty.backend.application.exception.ResourceNotFoundException;
import com.devaulty.backend.application.port.in.project.DeleteProjectUseCase;
import com.devaulty.backend.application.port.out.persistence.ProjectRepositoryPort;
import org.springframework.transaction.annotation.Transactional;

import java.util.UUID;

public class DeleteProjectImpl implements DeleteProjectUseCase {

    private final ProjectRepositoryPort projectRepositoryPort;

    public DeleteProjectImpl(ProjectRepositoryPort projectRepositoryPort) {
        this.projectRepositoryPort = projectRepositoryPort;
    }

    @Override
    @Transactional
    public void execute(UUID id) {

        if (!projectRepositoryPort.existsById(id)){
            throw new ResourceNotFoundException("Project", id);
        }

        projectRepositoryPort.deleteById(id);
    }
}
