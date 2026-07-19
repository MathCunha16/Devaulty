package com.devaulty.backend.adapter.out.persistence.note;

import org.springframework.data.domain.Page;
import org.springframework.data.domain.Pageable;
import org.springframework.data.jpa.repository.JpaRepository;

import java.util.UUID;

public interface SpringDataNoteRepository extends JpaRepository<NoteEntity, UUID> {
    Page<NoteEntity> findAllByProject_Id(UUID projectId, Pageable pageable);

    boolean existsByIdAndProject_Id(UUID id, UUID projectId);
}
