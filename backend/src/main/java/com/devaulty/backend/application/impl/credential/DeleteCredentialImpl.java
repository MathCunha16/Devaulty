package com.devaulty.backend.application.impl.credential;

import com.devaulty.backend.application.exception.MasterPasswordNotConfiguredException;
import com.devaulty.backend.application.exception.ResourceNotFoundException;
import com.devaulty.backend.application.exception.VaultLockedException;
import com.devaulty.backend.application.port.in.credential.DeleteCredentialUseCase;
import com.devaulty.backend.application.port.in.security.CheckMasterPasswordSetupUseCase;
import com.devaulty.backend.application.port.out.persistence.CredentialRepositoryPort;
import com.devaulty.backend.application.port.out.persistence.ItemTagRepositoryPort;
import com.devaulty.backend.application.port.out.persistence.ProjectRepositoryPort;
import com.devaulty.backend.application.port.out.security.MasterKeySessionPort;
import org.springframework.transaction.annotation.Transactional;

import java.util.UUID;

public class DeleteCredentialImpl implements DeleteCredentialUseCase {

    private final CredentialRepositoryPort credentialRepositoryPort;
    private final ProjectRepositoryPort projectRepositoryPort;
    private final MasterKeySessionPort masterKeySessionPort;
    private final CheckMasterPasswordSetupUseCase checkMasterPasswordSetupUseCase;
    private final ItemTagRepositoryPort itemTagRepositoryPort;

    public DeleteCredentialImpl(CredentialRepositoryPort credentialRepositoryPort, ProjectRepositoryPort projectRepositoryPort, MasterKeySessionPort masterKeySessionPort, CheckMasterPasswordSetupUseCase checkMasterPasswordSetupUseCase, ItemTagRepositoryPort itemTagRepositoryPort) {
        this.credentialRepositoryPort = credentialRepositoryPort;
        this.projectRepositoryPort = projectRepositoryPort;
        this.masterKeySessionPort = masterKeySessionPort;
        this.checkMasterPasswordSetupUseCase = checkMasterPasswordSetupUseCase;
        this.itemTagRepositoryPort = itemTagRepositoryPort;
    }

    @Override
    @Transactional
    public void execute(UUID projectId, UUID id) {

        if (checkMasterPasswordSetupUseCase.isSetupRequired()) throw new MasterPasswordNotConfiguredException();
        if (masterKeySessionPort.getKey() == null) throw new VaultLockedException();
        if (!projectRepositoryPort.existsById(projectId)) throw new ResourceNotFoundException("Project", projectId);
        if(!credentialRepositoryPort.existsByIdAndProjectId(id, projectId)) throw new ResourceNotFoundException("Credential", id);

        itemTagRepositoryPort.removeAllTagsFromItem("credential", id);
        credentialRepositoryPort.deleteById(id);
    }
}
