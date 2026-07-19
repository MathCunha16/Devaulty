package com.devaulty.backend.application.impl.snippet;

import com.devaulty.backend.application.exception.ResourceNotFoundException;
import com.devaulty.backend.application.port.in.snippet.DeleteSnippetUseCase;
import com.devaulty.backend.application.port.out.persistence.ItemTagRepositoryPort;
import com.devaulty.backend.application.port.out.persistence.ProjectRepositoryPort;
import com.devaulty.backend.application.port.out.persistence.SnippetRepositoryPort;
import org.springframework.transaction.annotation.Transactional;

import java.util.UUID;

public class DeleteSnippetImpl implements DeleteSnippetUseCase {

    private final SnippetRepositoryPort snippetRepositoryPort;
    private final ProjectRepositoryPort projectRepositoryPort;
    private final ItemTagRepositoryPort itemTagRepositoryPort;

    public DeleteSnippetImpl(SnippetRepositoryPort snippetRepositoryPort, ProjectRepositoryPort projectRepositoryPort, ItemTagRepositoryPort itemTagRepositoryPort) {
        this.snippetRepositoryPort = snippetRepositoryPort;
        this.projectRepositoryPort = projectRepositoryPort;
        this.itemTagRepositoryPort = itemTagRepositoryPort;
    }

    @Override
    @Transactional
    public void execute(UUID projectId, UUID id) {
        if(!projectRepositoryPort.existsById(projectId)) throw new ResourceNotFoundException("Project", projectId);
        if (!snippetRepositoryPort.existsByIdAndProjectId(id, projectId)) throw new ResourceNotFoundException("Snippet", id);

        itemTagRepositoryPort.removeAllTagsFromItem("snippet", id);
        snippetRepositoryPort.deleteById(id);
    }
}
