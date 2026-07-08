package com.devaulty.backend.adapter.out.persistence.project;

import org.springframework.data.jpa.repository.JpaRepository;

import java.util.UUID;

public interface SpringDataProjectRepository extends JpaRepository<ProjectEntity, UUID> {

}
