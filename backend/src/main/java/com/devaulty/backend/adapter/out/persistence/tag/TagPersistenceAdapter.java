package com.devaulty.backend.adapter.out.persistence.tag;

import com.devaulty.backend.adapter.out.persistence.project.SpringDataProjectRepository;
import com.devaulty.backend.application.port.out.persistence.TagRepositoryPort;
import com.devaulty.backend.domain.model.Tag;
import org.springframework.stereotype.Component;

import java.util.List;
import java.util.Optional;
import java.util.UUID;

@Component
public class TagPersistenceAdapter implements TagRepositoryPort {

    private final SpringDataTagRepository tagRepository;
    private final SpringDataProjectRepository projectRepository;
    private final TagMapper tagMapper;

    public TagPersistenceAdapter(SpringDataTagRepository tagRepository, SpringDataProjectRepository projectRepository, TagMapper tagMapper) {
        this.tagRepository = tagRepository;
        this.projectRepository = projectRepository;
        this.tagMapper = tagMapper;
    }

    @Override
    public Tag save(Tag tag) {
        TagEntity entity = tagMapper.toEntity(tag);
        entity.setProject(projectRepository.getReferenceById(tag.getProjectId()));
        TagEntity savedEntity = tagRepository.save(entity);
        return tagMapper.toDomain(savedEntity);
    }

    @Override
    public Optional<Tag> findById(UUID id) {
        Optional<TagEntity> entity = tagRepository.findById(id);
        return entity.map(tagMapper::toDomain);
    }

    @Override
    public List<Tag> searchByNameAndProjectId(UUID projectId, String name) {
        List<TagEntity> listResult = tagRepository.searchByNameAndProjectId(name, projectId);
        return listResult.stream().map(tagMapper::toDomain).toList();
    }

    @Override
    public List<Tag> findAllByProject(UUID projectId) {
        List<TagEntity> entities = tagRepository.findAllByProject_Id(projectId);
        return entities.stream().map(tagMapper::toDomain).toList();
    }

    @Override
    public void delete(UUID id) {
        tagRepository.deleteById(id);
    }

    @Override
    public boolean existsByNameAndProjectId(UUID projectId, String name) {
        return tagRepository.existsByNameAndProject_Id(name, projectId);
    }

    @Override
    public boolean existsByIdAndProjectId(UUID id, UUID projectId) {
        return tagRepository.existsByIdAndProject_Id(id, projectId);
    }

}
