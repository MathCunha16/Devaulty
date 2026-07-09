package com.devaulty.backend.application.port.in.project;

import com.devaulty.backend.domain.model.Project;
import org.springframework.data.domain.Page;

public interface GetAllProjectsUseCase {
    Page<Project> execute(int page, int size);
}
