package com.devaulty.backend.application.port.in.snippet;

import com.devaulty.backend.domain.model.Snippet;

import java.util.UUID;

public interface GetSnippetByIdUseCase {
    Snippet execute(UUID projectId, UUID id);
}
