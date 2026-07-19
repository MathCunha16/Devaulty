package com.devaulty.backend.application.impl.tag;

import com.devaulty.backend.application.exception.ResourceNotFoundException;
import com.devaulty.backend.application.port.in.tag.SearchTagByNameUseCase;
import com.devaulty.backend.application.port.out.persistence.ProjectRepositoryPort;
import com.devaulty.backend.application.port.out.persistence.TagRepositoryPort;
import com.devaulty.backend.domain.model.Tag;
import org.springframework.transaction.annotation.Transactional;

import java.util.List;
import java.util.UUID;

public class SearchTagByNameImpl implements SearchTagByNameUseCase {

    private final TagRepositoryPort tagRepositoryPort;
    private final ProjectRepositoryPort projectRepositoryPort;

    public SearchTagByNameImpl(TagRepositoryPort tagRepositoryPort, ProjectRepositoryPort projectRepositoryPort) {
        this.tagRepositoryPort = tagRepositoryPort;
        this.projectRepositoryPort = projectRepositoryPort;
    }

    @Override
    @Transactional(readOnly = true)
    public List<Tag> execute(UUID projectId, String name) {

        if(!projectRepositoryPort.existsById(projectId)) throw new ResourceNotFoundException("Project", projectId);

        return tagRepositoryPort.searchByNameAndProjectId(projectId, name);
    }
}
