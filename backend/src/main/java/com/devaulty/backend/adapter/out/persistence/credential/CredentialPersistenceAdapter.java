package com.devaulty.backend.adapter.out.persistence.credential;

import com.devaulty.backend.adapter.out.persistence.project.SpringDataProjectRepository;
import com.devaulty.backend.application.port.in.credential.CredentialSummary;
import com.devaulty.backend.application.port.out.persistence.CredentialRepositoryPort;
import com.devaulty.backend.domain.model.Credential;
import org.springframework.data.domain.Page;
import org.springframework.data.domain.PageRequest;
import org.springframework.data.domain.Pageable;
import org.springframework.data.domain.Sort;
import org.springframework.stereotype.Component;

import java.util.List;
import java.util.Optional;
import java.util.UUID;

@Component
public class CredentialPersistenceAdapter implements CredentialRepositoryPort {

    private final SpringDataCredentialRepository credentialRepository;
    private final SpringDataProjectRepository projectRepository;
    private final CredentialMapper credentialMapper;

    public CredentialPersistenceAdapter(SpringDataCredentialRepository credentialRepository, SpringDataProjectRepository projectRepository, CredentialMapper credentialMapper) {
        this.credentialRepository = credentialRepository;
        this.projectRepository = projectRepository;
        this.credentialMapper = credentialMapper;
    }

    @Override
    public Credential save(Credential credential) {
        CredentialEntity entity = credentialMapper.toEntity(credential);
        entity.setProject(projectRepository.getReferenceById(credential.getProjectId()));
        CredentialEntity savedEntity = credentialRepository.save(entity);
        return credentialMapper.toDomain(savedEntity);
    }

    @Override
    public Optional<Credential> findById(UUID id) {
        Optional<CredentialEntity> entity = credentialRepository.findById(id);
        return entity.map(credentialMapper::toDomain);
    }

    @Override
    public Page<CredentialSummary> findAllByProject(UUID projectId, int page, int size) {
        Pageable pageable = PageRequest.of(page, size, Sort.by(Sort.Direction.DESC, "createdAt"));
        return credentialRepository.findAllSummaryByProjectId(projectId, pageable);
    }

    @Override
    public void deleteById(UUID id) {
        credentialRepository.deleteById(id);
    }

    @Override
    public boolean existsByIdAndProjectId(UUID id, UUID projectId) {
        return credentialRepository.existsByIdAndProject_Id(id, projectId);
    }

    @Override
    public String getSupportedType() {
        return "credential";
    }

    @Override
    public List<UUID> findExistingIdsByProject(List<UUID> ids, UUID projectId) {
        if (ids == null || ids.isEmpty()) {
            return List.of();
        }
        return credentialRepository.findExistingIdsByProjectId(ids, projectId);
    }
}
