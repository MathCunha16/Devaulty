package com.devaulty.backend.application.port.out.persistence;

import com.devaulty.backend.domain.model.Note;
import org.springframework.data.domain.Page;

import java.util.Optional;
import java.util.UUID;

public interface NoteRepositoryPort extends ProjectScopedRepositoryPort {

    Note save(Note note);

    Optional<Note> findById(UUID id);

    Page<Note> findAllByProject(UUID projectId, int page, int size);

    void deleteById(UUID id);

    boolean existsByIdAndProjectId(UUID id, UUID projectId);
}
