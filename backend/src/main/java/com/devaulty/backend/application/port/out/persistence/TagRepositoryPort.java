package com.devaulty.backend.application.port.out.persistence;

import com.devaulty.backend.domain.model.Tag;

import java.util.List;
import java.util.Optional;
import java.util.UUID;

public interface TagRepositoryPort {
    Tag save(Tag tag);

    Optional<Tag> findById(UUID id);

    List<Tag> searchByNameAndProjectId(UUID projectId ,String name);

    List<Tag> findAllByProject(UUID projectId);

    void delete(UUID id);

    boolean existsByNameAndProjectId(UUID projectId ,String name);

    boolean existsByIdAndProjectId(UUID id, UUID projectId);
}
