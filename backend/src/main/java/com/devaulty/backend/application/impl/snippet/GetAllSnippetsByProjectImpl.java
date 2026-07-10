package com.devaulty.backend.application.impl.snippet;

import com.devaulty.backend.application.exception.ResourceNotFoundException;
import com.devaulty.backend.application.port.in.snippet.GetAllSnippetsByProjectUseCase;
import com.devaulty.backend.application.port.out.persistence.ProjectRepositoryPort;
import com.devaulty.backend.application.port.out.persistence.SnippetRepositoryPort;
import com.devaulty.backend.domain.model.Snippet;
import org.springframework.data.domain.Page;
import org.springframework.transaction.annotation.Transactional;

import java.util.UUID;

public class GetAllSnippetsByProjectImpl implements GetAllSnippetsByProjectUseCase {

    private final SnippetRepositoryPort snippetRepositoryPort;
    private final ProjectRepositoryPort projectRepositoryPort;

    public GetAllSnippetsByProjectImpl(SnippetRepositoryPort snippetRepositoryPort, ProjectRepositoryPort projectRepositoryPort) {
        this.snippetRepositoryPort = snippetRepositoryPort;
        this.projectRepositoryPort = projectRepositoryPort;
    }

    @Override
    @Transactional(readOnly = true)
    public Page<Snippet> execute(UUID projectId, int page, int size) {

        if(!projectRepositoryPort.existsById(projectId)){
            throw new ResourceNotFoundException("Project", projectId);
        }

        return snippetRepositoryPort.findAllByProject(projectId, page, size);
    }
}
