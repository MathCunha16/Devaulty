package com.devaulty.backend.adapter.out.persistence.tag;

import org.springframework.data.jpa.repository.JpaRepository;
import org.springframework.data.jpa.repository.Query;
import org.springframework.data.repository.query.Param;

import java.util.List;
import java.util.UUID;

public interface SpringDataTagRepository extends JpaRepository<TagEntity, UUID> {

    @Query("""
                SELECT t FROM TagEntity t
                WHERE t.project.id = :projectId
                  AND (:term IS NULL OR :term = '' OR LOWER(t.name) LIKE LOWER(CONCAT('%', :term, '%')))
                ORDER BY t.name ASC
            """)
    List<TagEntity> searchByNameAndProjectId(@Param("term") String term, @Param("projectId") UUID projectId);

    List<TagEntity> findAllByProject_Id(UUID projectId);

    boolean existsByNameAndProject_Id(String name, UUID projectId);

    boolean existsByIdAndProject_Id(UUID id, UUID projectId);
}
