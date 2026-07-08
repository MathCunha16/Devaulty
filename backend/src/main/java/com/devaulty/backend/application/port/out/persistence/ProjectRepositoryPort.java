package com.devaulty.backend.application.port.out.persistence;

import com.devaulty.backend.domain.model.Project;
import org.springframework.data.domain.Page;

import java.util.Optional;
import java.util.UUID;

public interface ProjectRepositoryPort {

    Project save(Project project);

    Page<Project> findAll(int page, int size);

    Optional<Project> findById(UUID id);

    void deleteById(UUID id);
}
