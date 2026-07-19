package com.devaulty.backend.application.port.out.persistence;

import com.devaulty.backend.domain.model.Link;
import org.springframework.data.domain.Page;

import java.util.Optional;
import java.util.UUID;

public interface LinkRepositoryPort extends ProjectScopedRepositoryPort {
    Link save(Link link);

    Optional<Link> findById(UUID id);

    Page<Link> findAllByProject(UUID projectId, int page, int size);

    void deleteById(UUID id);

    boolean existsByIdAndProjectId(UUID id, UUID projectId);
}
