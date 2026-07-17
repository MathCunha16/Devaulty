package com.devaulty.backend.application.port.in.credential;

import org.springframework.data.domain.Page;

import java.util.UUID;

public interface GetAllCredentialsByProjectUseCase {
    Page<CredentialSummary> execute(UUID projectId, int page, int size);
}
