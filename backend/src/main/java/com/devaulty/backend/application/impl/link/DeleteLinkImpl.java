package com.devaulty.backend.application.impl.link;

import com.devaulty.backend.application.exception.ResourceNotFoundException;
import com.devaulty.backend.application.port.in.link.DeleteLinkUseCase;
import com.devaulty.backend.application.port.out.persistence.ItemTagRepositoryPort;
import com.devaulty.backend.application.port.out.persistence.LinkRepositoryPort;
import com.devaulty.backend.application.port.out.persistence.ProjectRepositoryPort;
import org.springframework.transaction.annotation.Transactional;

import java.util.UUID;

public class DeleteLinkImpl implements DeleteLinkUseCase {

    private final LinkRepositoryPort linkRepositoryPort;
    private final ProjectRepositoryPort projectRepositoryPort;
    private final ItemTagRepositoryPort itemTagRepositoryPort;

    public DeleteLinkImpl(LinkRepositoryPort linkRepositoryPort, ProjectRepositoryPort projectRepositoryPort, ItemTagRepositoryPort itemTagRepositoryPort) {
        this.linkRepositoryPort = linkRepositoryPort;
        this.projectRepositoryPort = projectRepositoryPort;
        this.itemTagRepositoryPort = itemTagRepositoryPort;
    }

    @Override
    @Transactional
    public void execute(UUID projectId, UUID id) {

        if (!projectRepositoryPort.existsById(projectId)) throw new ResourceNotFoundException("Project", projectId);
        if (!linkRepositoryPort.existsByIdAndProjectId(id, projectId)) throw new ResourceNotFoundException("Link", id);

        itemTagRepositoryPort.removeAllTagsFromItem("link", id);
        linkRepositoryPort.deleteById(id);
    }
}
