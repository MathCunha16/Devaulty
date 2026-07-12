package com.devaulty.backend.adapter.out.persistence.problem;

import com.devaulty.backend.adapter.out.persistence.project.SpringDataProjectRepository;
import com.devaulty.backend.application.port.out.persistence.ProblemRepositoryPort;
import com.devaulty.backend.domain.model.Problem;
import org.springframework.data.domain.Page;
import org.springframework.data.domain.PageRequest;
import org.springframework.data.domain.Pageable;
import org.springframework.data.domain.Sort;
import org.springframework.stereotype.Component;

import java.util.Optional;
import java.util.UUID;

@Component
public class ProblemPersistenceAdapter implements ProblemRepositoryPort {

    private final SpringDataProblemRepository problemRepository;
    private final SpringDataProjectRepository projectRepository;
    private final ProblemMapper problemMapper;

    public ProblemPersistenceAdapter(SpringDataProblemRepository problemRepository, SpringDataProjectRepository projectRepository, ProblemMapper problemMapper) {
        this.problemRepository = problemRepository;
        this.projectRepository = projectRepository;
        this.problemMapper = problemMapper;
    }

    @Override
    public Problem save(Problem problem) {
        ProblemEntity entity = problemMapper.toEntity(problem);
        entity.setProject(projectRepository.getReferenceById(problem.getProjectId()));
        ProblemEntity savedEntity = problemRepository.save(entity);
        return problemMapper.toDomain(savedEntity);
    }

    @Override
    public Optional<Problem> findById(UUID id) {
        Optional<ProblemEntity> entity = problemRepository.findById(id);
        return entity.map(problemMapper::toDomain);
    }

    @Override
    public Page<Problem> findAllByProject(UUID projectId, int page, int size) {
        Pageable pageable = PageRequest.of(page, size, Sort.by(Sort.Direction.DESC, "createdAt"));
        Page<ProblemEntity> entities = problemRepository.findAllByProject_Id(projectId, pageable);
        return entities.map(problemMapper::toDomain);
    }

    @Override
    public void deleteById(UUID id) {
        problemRepository.deleteById(id);
    }
}
