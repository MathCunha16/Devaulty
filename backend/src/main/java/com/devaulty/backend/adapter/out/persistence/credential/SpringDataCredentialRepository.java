package com.devaulty.backend.adapter.out.persistence.credential;

import org.springframework.data.domain.Page;
import org.springframework.data.domain.Pageable;
import org.springframework.data.jpa.repository.JpaRepository;

import java.util.UUID;

public interface SpringDataCredentialRepository extends JpaRepository<CredentialEntity, UUID> {
    Page<CredentialEntity> findAllByProject_Id(UUID projectId, Pageable pageable);
}
