package com.devaulty.backend.application.port.in.snippet;

import java.util.UUID;

public interface DeleteSnippetUseCase {
    void execute(UUID projectId, UUID id);
}
