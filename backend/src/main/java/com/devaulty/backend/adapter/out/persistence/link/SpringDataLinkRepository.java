package com.devaulty.backend.adapter.out.persistence.link;

import org.springframework.data.domain.Page;
import org.springframework.data.domain.Pageable;
import org.springframework.data.jpa.repository.JpaRepository;
import org.springframework.data.jpa.repository.Query;
import org.springframework.data.repository.query.Param;

import java.util.List;
import java.util.UUID;

public interface SpringDataLinkRepository extends JpaRepository<LinkEntity, UUID> {
    Page<LinkEntity> findAllByProject_Id(UUID projectId, Pageable pageable);

    boolean existsByIdAndProject_Id(UUID id, UUID projectId);

    @Query("SELECT l.id FROM LinkEntity l WHERE l.id IN :ids AND l.project.id = :projectId")
    List<UUID> findExistingIdsByProjectId(
            @Param("ids") List<UUID> ids,
            @Param("projectId") UUID projectId
    );
}
