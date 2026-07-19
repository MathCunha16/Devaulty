package com.devaulty.backend.application.impl.tag;

import com.devaulty.backend.application.exception.ResourceNotFoundException;
import com.devaulty.backend.application.port.in.tag.GetTagByIdUseCase;
import com.devaulty.backend.application.port.out.persistence.ProjectRepositoryPort;
import com.devaulty.backend.application.port.out.persistence.TagRepositoryPort;
import com.devaulty.backend.domain.model.Tag;
import org.springframework.transaction.annotation.Transactional;

import java.util.UUID;

public class GetTagByIdImpl implements GetTagByIdUseCase {

    private final TagRepositoryPort tagRepositoryPort;
    private final ProjectRepositoryPort projectRepositoryPort;

    public GetTagByIdImpl(TagRepositoryPort tagRepositoryPort, ProjectRepositoryPort projectRepositoryPort) {
        this.tagRepositoryPort = tagRepositoryPort;
        this.projectRepositoryPort = projectRepositoryPort;
    }

    @Override
    @Transactional(readOnly = true)
    public Tag execute(UUID projectId, UUID id) {

        if (!projectRepositoryPort.existsById(projectId)) throw new ResourceNotFoundException("Project", projectId);

        return tagRepositoryPort.findById(id)
                .filter(t -> projectId.equals(t.getProjectId()))
                .orElseThrow(() -> new ResourceNotFoundException("Tag", id));

    }
}
