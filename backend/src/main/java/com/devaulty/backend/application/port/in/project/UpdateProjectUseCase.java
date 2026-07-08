package com.devaulty.backend.application.port.in.project;

import com.devaulty.backend.domain.model.Project;

import java.util.UUID;

public interface UpdateProjectUseCase {
    Project execute(UUID id, CreateProjectCommand command);
}
