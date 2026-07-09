package com.devaulty.backend.application.impl.snippet;

import com.devaulty.backend.application.port.in.snippet.DeleteSnippetUseCase;
import com.devaulty.backend.application.port.out.persistence.ProjectRepositoryPort;
import com.devaulty.backend.application.port.out.persistence.SnippetRepositoryPort;
import org.springframework.transaction.annotation.Transactional;

import java.util.UUID;

public class DeleteSnippetImpl implements DeleteSnippetUseCase {

    private final SnippetRepositoryPort snippetRepositoryPort;
    private final ProjectRepositoryPort projectRepositoryPort;

    public DeleteSnippetImpl(SnippetRepositoryPort snippetRepositoryPort, ProjectRepositoryPort projectRepositoryPort) {
        this.snippetRepositoryPort = snippetRepositoryPort;
        this.projectRepositoryPort = projectRepositoryPort;
    }

    @Override
    @Transactional
    public void execute(UUID projectId, UUID id) {

    }
}
