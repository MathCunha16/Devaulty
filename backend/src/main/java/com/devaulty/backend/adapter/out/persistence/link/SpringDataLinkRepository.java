package com.devaulty.backend.adapter.out.persistence.link;

import org.springframework.data.domain.Page;
import org.springframework.data.domain.Pageable;
import org.springframework.data.jpa.repository.JpaRepository;

import java.util.UUID;

public interface SpringDataLinkRepository extends JpaRepository<LinkEntity, UUID> {
    Page<LinkEntity> findAllByProject_Id(UUID projectId, Pageable pageable);

    boolean existsByIdAndProject_Id(UUID id, UUID projectId);
}
