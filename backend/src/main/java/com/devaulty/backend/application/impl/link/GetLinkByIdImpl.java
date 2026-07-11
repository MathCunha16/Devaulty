package com.devaulty.backend.application.impl.link;

import com.devaulty.backend.application.exception.ResourceNotFoundException;
import com.devaulty.backend.application.port.in.link.GetLinkByIdUseCase;
import com.devaulty.backend.application.port.out.persistence.LinkRepositoryPort;
import com.devaulty.backend.application.port.out.persistence.ProjectRepositoryPort;
import com.devaulty.backend.domain.model.Link;
import org.springframework.transaction.annotation.Transactional;

import java.util.UUID;

public class GetLinkByIdImpl implements GetLinkByIdUseCase {

    private final LinkRepositoryPort linkRepositoryPort;
    private final ProjectRepositoryPort projectRepositoryPort;

    public GetLinkByIdImpl(LinkRepositoryPort linkRepositoryPort, ProjectRepositoryPort projectRepositoryPort) {
        this.linkRepositoryPort = linkRepositoryPort;
        this.projectRepositoryPort = projectRepositoryPort;
    }

    @Override
    @Transactional(readOnly = true)
    public Link execute(UUID projectId ,UUID id) {

        if (!projectRepositoryPort.existsById(projectId)) throw new ResourceNotFoundException("Project", projectId);

        Link link = linkRepositoryPort.findById(id)
                .orElseThrow(() -> new ResourceNotFoundException("Link", id));

        if (!link.getProjectId().equals(projectId)) throw new ResourceNotFoundException("Link", id);

        return link;
    }
}
