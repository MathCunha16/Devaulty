package com.devaulty.backend.adapter.out.persistence.note;

import org.springframework.data.domain.Page;
import org.springframework.data.domain.Pageable;
import org.springframework.data.jpa.repository.JpaRepository;
import org.springframework.data.jpa.repository.Query;
import org.springframework.data.repository.query.Param;

import java.util.List;
import java.util.UUID;

public interface SpringDataNoteRepository extends JpaRepository<NoteEntity, UUID> {
    Page<NoteEntity> findAllByProject_Id(UUID projectId, Pageable pageable);

    boolean existsByIdAndProject_Id(UUID id, UUID projectId);

    @Query("SELECT n.id FROM NoteEntity n WHERE n.id IN :ids AND n.project.id = :projectId")
    List<UUID> findExistingIdsByProjectId(
            @Param("ids") List<UUID> ids,
            @Param("projectId") UUID projectId
    );
}
