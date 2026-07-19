package com.devaulty.backend.adapter.out.persistence.link;

import com.devaulty.backend.adapter.out.persistence.project.SpringDataProjectRepository;
import com.devaulty.backend.application.port.out.persistence.LinkRepositoryPort;
import com.devaulty.backend.domain.model.Link;
import org.springframework.data.domain.Page;
import org.springframework.data.domain.PageRequest;
import org.springframework.data.domain.Pageable;
import org.springframework.data.domain.Sort;
import org.springframework.stereotype.Component;

import java.util.Optional;
import java.util.UUID;

@Component
public class LinkPersistenceAdapter implements LinkRepositoryPort {

    private final SpringDataLinkRepository linkRepository;
    private final SpringDataProjectRepository projectRepository;
    private final LinkMapper linkMapper;

    public LinkPersistenceAdapter(SpringDataLinkRepository linkRepository, SpringDataProjectRepository projectRepository, LinkMapper linkMapper) {
        this.linkRepository = linkRepository;
        this.projectRepository = projectRepository;
        this.linkMapper = linkMapper;
    }

    @Override
    public Link save(Link link) {
        LinkEntity linkEntity = linkMapper.toEntity(link);
        linkEntity.setProject(projectRepository.getReferenceById(link.getProjectId()));
        LinkEntity savedLinkEntity = linkRepository.save(linkEntity);
        return linkMapper.toDomain(savedLinkEntity);
    }

    @Override
    public Optional<Link> findById(UUID id) {
        Optional<LinkEntity> entity = linkRepository.findById(id);
        return entity.map(linkMapper::toDomain);
    }

    @Override
    public Page<Link> findAllByProject(UUID projectId, int page, int size) {
        Pageable pageable = PageRequest.of(page, size, Sort.by(Sort.Direction.DESC, "createdAt"));
        Page<LinkEntity> entities = linkRepository.findAllByProject_Id(projectId, pageable);
        return entities.map(linkMapper::toDomain);
    }

    @Override
    public void deleteById(UUID id) {
        linkRepository.deleteById(id);
    }

    @Override
    public boolean existsByIdAndProjectId(UUID id, UUID projectId) {
        return linkRepository.existsByIdAndProject_Id(id, projectId);
    }
}
