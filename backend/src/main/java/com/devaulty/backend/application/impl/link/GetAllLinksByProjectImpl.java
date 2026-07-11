package com.devaulty.backend.application.impl.link;

import com.devaulty.backend.application.exception.ResourceNotFoundException;
import com.devaulty.backend.application.port.in.link.GetAllLinksByProjectUseCase;
import com.devaulty.backend.application.port.out.persistence.LinkRepositoryPort;
import com.devaulty.backend.application.port.out.persistence.ProjectRepositoryPort;
import com.devaulty.backend.domain.model.Link;
import org.springframework.data.domain.Page;
import org.springframework.transaction.annotation.Transactional;

import java.util.UUID;

public class GetAllLinksByProjectImpl implements GetAllLinksByProjectUseCase {

    private final LinkRepositoryPort linkRepositoryPort;
    private final ProjectRepositoryPort projectRepositoryPort;

    public GetAllLinksByProjectImpl(LinkRepositoryPort linkRepositoryPort, ProjectRepositoryPort projectRepositoryPort) {
        this.linkRepositoryPort = linkRepositoryPort;
        this.projectRepositoryPort = projectRepositoryPort;
    }

    @Override
    @Transactional(readOnly = true)
    public Page<Link> execute(UUID projectId, int page, int size) {

        if(!projectRepositoryPort.existsById(projectId)) throw new ResourceNotFoundException("Project", projectId);

        return linkRepositoryPort.findAllByProject(projectId, page, size);
    }
}
