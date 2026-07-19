package com.devaulty.backend.adapter.out.persistence.problem;

import org.springframework.data.domain.Page;
import org.springframework.data.domain.Pageable;
import org.springframework.data.jpa.repository.JpaRepository;
import org.springframework.data.jpa.repository.Query;
import org.springframework.data.repository.query.Param;

import java.util.List;
import java.util.UUID;

public interface SpringDataProblemRepository extends JpaRepository<ProblemEntity, UUID> {
    Page<ProblemEntity> findAllByProject_Id(UUID projectId, Pageable pageable);

    boolean existsByIdAndProject_Id(UUID id, UUID projectId);

    @Query("SELECT p.id FROM ProblemEntity p WHERE p.id IN :ids AND p.project.id = :projectId")
    List<UUID> findExistingIdsByProjectId(
            @Param("ids") List<UUID> ids,
            @Param("projectId") UUID projectId
    );
}
