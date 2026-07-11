package com.devaulty.backend.application.impl.link;

import com.devaulty.backend.application.exception.ResourceNotFoundException;
import com.devaulty.backend.application.port.in.link.DeleteLinkUseCase;
import com.devaulty.backend.application.port.out.persistence.LinkRepositoryPort;
import com.devaulty.backend.application.port.out.persistence.ProjectRepositoryPort;
import org.springframework.transaction.annotation.Transactional;

import java.util.UUID;

public class DeleteLinkImpl implements DeleteLinkUseCase {

    private final LinkRepositoryPort linkRepositoryPort;
    private final ProjectRepositoryPort projectRepositoryPort;

    public DeleteLinkImpl(LinkRepositoryPort linkRepositoryPort, ProjectRepositoryPort projectRepositoryPort) {
        this.linkRepositoryPort = linkRepositoryPort;
        this.projectRepositoryPort = projectRepositoryPort;
    }

    @Override
    @Transactional
    public void execute(UUID projectId, UUID id) {

        if (!projectRepositoryPort.existsById(projectId)) throw new ResourceNotFoundException("Project", projectId);

        if (linkRepositoryPort.findById(id)
                .filter(link -> projectId.equals(link.getProjectId()))
                .isEmpty()) {
            throw new ResourceNotFoundException("Link", id);
        }

        linkRepositoryPort.deleteById(id);
    }
}
