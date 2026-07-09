package com.devaulty.backend.application.impl.project;

import com.devaulty.backend.application.port.in.project.GetAllProjectsUseCase;
import com.devaulty.backend.application.port.out.persistence.ProjectRepositoryPort;
import com.devaulty.backend.domain.model.Project;
import org.springframework.data.domain.Page;
import org.springframework.transaction.annotation.Transactional;

public class GetAllProjectsImpl implements GetAllProjectsUseCase {

    private final ProjectRepositoryPort projectRepositoryPort;

    public GetAllProjectsImpl(ProjectRepositoryPort projectRepositoryPort) {
        this.projectRepositoryPort = projectRepositoryPort;
    }

    @Override
    @Transactional(readOnly = true)
    public Page<Project> execute(int page, int size) {
        return projectRepositoryPort.findAll(page, size);
    }
}
