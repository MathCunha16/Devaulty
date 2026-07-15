package com.devaulty.backend.application.port.out.persistence;

import com.devaulty.backend.application.port.in.credential.CredentialSummary;
import com.devaulty.backend.domain.model.Credential;
import org.springframework.data.domain.Page;

import java.util.Optional;
import java.util.UUID;

public interface CredentialRepositoryPort {

    Credential save(Credential credential);

    Optional<Credential> findById(UUID id);

    Page<CredentialSummary> findAllByProject(UUID projectId, int page, int size);

    void deleteById(UUID id);
}
