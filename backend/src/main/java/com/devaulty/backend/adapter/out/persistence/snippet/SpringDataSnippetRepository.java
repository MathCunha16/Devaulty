package com.devaulty.backend.adapter.out.persistence.snippet;

import org.springframework.data.domain.Page;
import org.springframework.data.domain.Pageable;
import org.springframework.data.jpa.repository.JpaRepository;
import org.springframework.data.jpa.repository.Query;
import org.springframework.data.repository.query.Param;

import java.util.List;
import java.util.UUID;

public interface SpringDataSnippetRepository extends JpaRepository<SnippetEntity, UUID> {
    Page<SnippetEntity> findAllByProject_Id(UUID projectId, Pageable pageable);

    boolean existsByIdAndProject_Id(UUID id, UUID projectId);

    @Query("SELECT s.id FROM SnippetEntity s WHERE s.id IN :ids AND s.project.id = :projectId")
    List<UUID> findExistingIdsByProjectId(
            @Param("ids") List<UUID> ids,
            @Param("projectId") UUID projectId
    );
}
