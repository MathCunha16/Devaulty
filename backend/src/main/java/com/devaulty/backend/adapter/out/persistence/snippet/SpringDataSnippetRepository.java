package com.devaulty.backend.adapter.out.persistence.snippet;

import org.springframework.data.domain.Page;
import org.springframework.data.domain.Pageable;
import org.springframework.data.jpa.repository.JpaRepository;

import java.util.UUID;

public interface SpringDataSnippetRepository extends JpaRepository<SnippetEntity, UUID> {
    Page<SnippetEntity> findAllByProject_Id(UUID projectId, Pageable pageable);

    boolean existsByIdAndProject_Id(UUID id, UUID projectId);
}
