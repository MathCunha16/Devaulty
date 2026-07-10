package com.devaulty.backend.application.impl.snippet;

import com.devaulty.backend.application.exception.ResourceNotFoundException;
import com.devaulty.backend.application.port.in.snippet.GetSnippetByIdUseCase;
import com.devaulty.backend.application.port.out.persistence.ProjectRepositoryPort;
import com.devaulty.backend.application.port.out.persistence.SnippetRepositoryPort;
import com.devaulty.backend.domain.model.Snippet;
import org.springframework.transaction.annotation.Transactional;

import java.util.UUID;

public class GetSnippetByIdImpl implements GetSnippetByIdUseCase {

    private final SnippetRepositoryPort snippetRepositoryPort;
    private final ProjectRepositoryPort projectRepositoryPort;

    public GetSnippetByIdImpl(SnippetRepositoryPort snippetRepositoryPort, ProjectRepositoryPort projectRepositoryPort) {
        this.snippetRepositoryPort = snippetRepositoryPort;
        this.projectRepositoryPort = projectRepositoryPort;
    }

    @Override
    @Transactional(readOnly = true)
    public Snippet execute(UUID projectId, UUID id) {

        if(!projectRepositoryPort.existsById(projectId)) throw new ResourceNotFoundException("Project", projectId);

        Snippet snippet = snippetRepositoryPort.findById(id)
                .orElseThrow(() -> new ResourceNotFoundException("Snippet", id));

        if(!snippet.getProjectId().equals(projectId)) throw new ResourceNotFoundException("Snippet", id);

        return snippet;
    }
}
