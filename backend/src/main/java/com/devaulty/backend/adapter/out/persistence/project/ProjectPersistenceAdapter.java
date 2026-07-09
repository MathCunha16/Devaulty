package com.devaulty.backend.adapter.out.persistence.project;

import com.devaulty.backend.application.port.out.persistence.ProjectRepositoryPort;
import com.devaulty.backend.domain.model.Project;
import org.springframework.data.domain.Page;
import org.springframework.data.domain.PageRequest;
import org.springframework.data.domain.Pageable;
import org.springframework.data.domain.Sort;
import org.springframework.stereotype.Component;

import java.util.Optional;
import java.util.UUID;

@Component
public class ProjectPersistenceAdapter implements ProjectRepositoryPort {

    private final SpringDataProjectRepository projectRepository;
    private final ProjectMapper projectMapper;

    public ProjectPersistenceAdapter(SpringDataProjectRepository projectRepository, ProjectMapper projectMapper) {
        this.projectRepository = projectRepository;
        this.projectMapper = projectMapper;
    }

    @Override
    public Project save(Project project) {
        ProjectEntity entity = projectRepository.save(projectMapper.toEntity(project));
        return projectMapper.toDomain(entity);
    }

    @Override
    public Page<Project> findAll(int page, int size) {
        Pageable pageable = PageRequest.of(page, size, Sort.by(Sort.Direction.DESC, "createdAt"));
        Page<ProjectEntity> entityPage = projectRepository.findAll(pageable);
        return entityPage.map(projectMapper::toDomain);
    }

    @Override
    public Optional<Project> findById(UUID id) {
        Optional<ProjectEntity> entity = projectRepository.findById(id);
        return entity.map(projectMapper::toDomain);
    }

    @Override
    public void deleteById(UUID id) {
        projectRepository.deleteById(id);
    }

    @Override
    public boolean existsById(UUID id) {
        return projectRepository.existsById(id);
    }
}
