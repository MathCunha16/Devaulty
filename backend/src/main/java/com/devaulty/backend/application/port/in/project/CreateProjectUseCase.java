package com.devaulty.backend.application.port.in.project;

import com.devaulty.backend.domain.model.Project;

public interface CreateProjectUseCase {
    Project execute(CreateProjectCommand command);
}
