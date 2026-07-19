package com.devaulty.backend.application.port.out.persistence;

import com.devaulty.backend.domain.model.Problem;
import org.springframework.data.domain.Page;

import java.util.Optional;
import java.util.UUID;

public interface ProblemRepositoryPort {
    Problem save(Problem problem);

    Optional<Problem> findById(UUID id);

    Page<Problem> findAllByProject(UUID projectId, int page, int size);

    void deleteById(UUID id);

    boolean existsByIdAndProjectId(UUID id, UUID projectId);
}
