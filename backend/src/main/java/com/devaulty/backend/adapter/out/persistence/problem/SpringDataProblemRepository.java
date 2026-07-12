package com.devaulty.backend.adapter.out.persistence.problem;

import org.springframework.data.domain.Page;
import org.springframework.data.domain.Pageable;
import org.springframework.data.jpa.repository.JpaRepository;

import java.util.UUID;

public interface SpringDataProblemRepository extends JpaRepository<ProblemEntity, UUID> {
    Page<ProblemEntity> findAllByProject_Id(UUID projectId, Pageable pageable);
}
