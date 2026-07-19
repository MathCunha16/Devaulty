package com.devaulty.backend.application.port.out.persistence;

import java.util.UUID;

public interface ProjectScopedRepositoryPort {
    String getSupportedType();
    boolean existsByIdAndProjectId(UUID id, UUID projectId);
}
