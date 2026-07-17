package com.devaulty.backend.application.impl.credential;

import com.devaulty.backend.application.exception.MasterPasswordNotConfiguredException;
import com.devaulty.backend.application.exception.ResourceNotFoundException;
import com.devaulty.backend.application.exception.VaultLockedException;
import com.devaulty.backend.application.port.in.credential.DeleteCredentialUseCase;
import com.devaulty.backend.application.port.in.security.CheckMasterPasswordSetupUseCase;
import com.devaulty.backend.application.port.out.persistence.CredentialRepositoryPort;
import com.devaulty.backend.application.port.out.persistence.ProjectRepositoryPort;
import com.devaulty.backend.application.port.out.security.MasterKeySessionPort;
import org.springframework.transaction.annotation.Transactional;

import java.util.UUID;

public class DeleteCredentialImpl implements DeleteCredentialUseCase {

    private final CredentialRepositoryPort credentialRepositoryPort;
    private final ProjectRepositoryPort projectRepositoryPort;
    private final MasterKeySessionPort masterKeySessionPort;
    private final CheckMasterPasswordSetupUseCase checkMasterPasswordSetupUseCase;

    public DeleteCredentialImpl(CredentialRepositoryPort credentialRepositoryPort, ProjectRepositoryPort projectRepositoryPort, MasterKeySessionPort masterKeySessionPort, CheckMasterPasswordSetupUseCase checkMasterPasswordSetupUseCase) {
        this.credentialRepositoryPort = credentialRepositoryPort;
        this.projectRepositoryPort = projectRepositoryPort;
        this.masterKeySessionPort = masterKeySessionPort;
        this.checkMasterPasswordSetupUseCase = checkMasterPasswordSetupUseCase;
    }

    @Override
    @Transactional
    public void execute(UUID projectId, UUID id) {

        if (checkMasterPasswordSetupUseCase.isSetupRequired()) throw new MasterPasswordNotConfiguredException();
        if (masterKeySessionPort.getKey() == null) throw new VaultLockedException();
        if (!projectRepositoryPort.existsById(projectId)) throw new ResourceNotFoundException("Project", projectId);

        if(credentialRepositoryPort.findById(id)
                .filter(c -> projectId.equals(c.getProjectId()))
                .isEmpty()) {
            throw new ResourceNotFoundException("Credential", id);
        }

        credentialRepositoryPort.deleteById(id);
    }
}
