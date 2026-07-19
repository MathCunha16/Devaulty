package com.devaulty.backend.adapter.out.persistence.snippet;

import com.devaulty.backend.adapter.out.persistence.project.SpringDataProjectRepository;
import com.devaulty.backend.application.port.out.persistence.SnippetRepositoryPort;
import com.devaulty.backend.domain.model.Snippet;
import org.springframework.data.domain.Page;
import org.springframework.data.domain.PageRequest;
import org.springframework.data.domain.Pageable;
import org.springframework.data.domain.Sort;
import org.springframework.stereotype.Component;

import java.util.Optional;
import java.util.UUID;

@Component
public class SnippetPersistenceAdapter implements SnippetRepositoryPort {

    private final SpringDataSnippetRepository snippetRepository;
    private final SpringDataProjectRepository projectRepository;
    private final SnippetMapper snippetMapper;

    public SnippetPersistenceAdapter(SpringDataSnippetRepository snippetRepository, SpringDataProjectRepository projectRepository, SnippetMapper snippetMapper) {
        this.snippetRepository = snippetRepository;
        this.projectRepository = projectRepository;
        this.snippetMapper = snippetMapper;
    }

    @Override
    public Snippet save(Snippet snippet) {
        SnippetEntity entity = snippetMapper.toEntity(snippet);
        entity.setProject(projectRepository.getReferenceById(snippet.getProjectId()));
        SnippetEntity savedEntity = snippetRepository.save(entity);
        return snippetMapper.toDomain(savedEntity);
    }

    @Override
    public Optional<Snippet> findById(UUID id) {
        Optional<SnippetEntity> entity = snippetRepository.findById(id);
        return entity.map(snippetMapper::toDomain);
    }

    @Override
    public Page<Snippet> findAllByProject(UUID projectId, int page, int size) {
        Pageable pageable = PageRequest.of(page, size, Sort.by(Sort.Direction.DESC, "createdAt"));
        Page<SnippetEntity> snippets = snippetRepository.findAllByProject_Id(projectId, pageable);
        return snippets.map(snippetMapper::toDomain);
    }

    @Override
    public void deleteById(UUID id) {
        snippetRepository.deleteById(id);
    }

    @Override
    public boolean existsByIdAndProjectId(UUID id, UUID projectId) {
        return snippetRepository.existsByIdAndProject_Id(id, projectId);
    }

    @Override
    public String getSupportedType() {
        return "snippet";
    }
}
