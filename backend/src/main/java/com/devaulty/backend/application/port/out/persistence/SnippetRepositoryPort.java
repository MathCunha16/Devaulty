package com.devaulty.backend.application.port.out.persistence;

import com.devaulty.backend.domain.model.Snippet;
import org.springframework.data.domain.Page;

import java.util.Optional;
import java.util.UUID;

public interface SnippetRepositoryPort {
    Snippet save(Snippet snippet);

    Optional<Snippet> findById(UUID id);

    Page<Snippet> findAllByProject(UUID projectId, int page, int size);

    void deleteById(UUID id);
}
