package com.devaulty.backend.adapter.out.persistence.credential;

import com.devaulty.backend.application.port.in.credential.CredentialSummary;
import org.springframework.data.domain.Page;
import org.springframework.data.domain.Pageable;
import org.springframework.data.jpa.repository.JpaRepository;
import org.springframework.data.jpa.repository.Query;
import org.springframework.data.repository.query.Param;

import java.util.UUID;

public interface SpringDataCredentialRepository extends JpaRepository<CredentialEntity, UUID> {
    @Query("""
        SELECT new com.devaulty.backend.application.port.in.credential.CredentialSummary(
            c.id,
            c.project.id,
            c.title,
            c.secretType,
            c.relatedUrl,
            c.createdAt,
            c.updatedAt
        )
        FROM CredentialEntity c
        WHERE c.project.id = :projectId
    """)
    Page<CredentialSummary> findAllSummaryByProjectId(
            @Param("projectId") UUID projectId,
            Pageable pageable
    );

    boolean existsByIdAndProject_Id(UUID id, UUID projectId);
}