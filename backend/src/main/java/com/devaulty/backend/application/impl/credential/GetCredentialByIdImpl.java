package com.devaulty.backend.application.impl.credential;

import com.devaulty.backend.application.exception.MasterPasswordNotConfiguredException;
import com.devaulty.backend.application.exception.ResourceNotFoundException;
import com.devaulty.backend.application.exception.VaultLockedException;
import com.devaulty.backend.application.port.in.credential.DecryptedCredential;
import com.devaulty.backend.application.port.in.credential.GetCredentialByIdUseCase;
import com.devaulty.backend.application.port.in.security.CheckMasterPasswordSetupUseCase;
import com.devaulty.backend.application.port.out.persistence.CredentialRepositoryPort;
import com.devaulty.backend.application.port.out.persistence.ProjectRepositoryPort;
import com.devaulty.backend.application.port.out.security.CryptoPort;
import com.devaulty.backend.application.port.out.security.MasterKeySessionPort;
import com.devaulty.backend.domain.model.Credential;
import org.springframework.transaction.annotation.Transactional;

import javax.crypto.SecretKey;
import java.util.UUID;

public class GetCredentialByIdImpl implements GetCredentialByIdUseCase {

    private final CredentialRepositoryPort credentialRepositoryPort;
    private final ProjectRepositoryPort projectRepositoryPort;
    private final CryptoPort cryptoPort;
    private final MasterKeySessionPort masterKeySessionPort;
    private final CheckMasterPasswordSetupUseCase checkMasterPasswordSetupUseCase;

    public GetCredentialByIdImpl(CredentialRepositoryPort credentialRepositoryPort, ProjectRepositoryPort projectRepositoryPort, CryptoPort cryptoPort, MasterKeySessionPort masterKeySessionPort, CheckMasterPasswordSetupUseCase checkMasterPasswordSetupUseCase) {
        this.credentialRepositoryPort = credentialRepositoryPort;
        this.projectRepositoryPort = projectRepositoryPort;
        this.cryptoPort = cryptoPort;
        this.masterKeySessionPort = masterKeySessionPort;
        this.checkMasterPasswordSetupUseCase = checkMasterPasswordSetupUseCase;
    }

    /**
     * Retrieves a single credential by its ID within a specific project and decrypts its payload.
     *
     * <p><b>⚠ MEMORY OWNERSHIP:</b> the returned {@link DecryptedCredential#decryptedPayload()}
     * is NOT zeroed here — ownership passes to the caller. Whoever consumes it MUST
     * wipe it after use:
     * <pre>{@code Arrays.fill(decryptedCredential.decryptedBytesPayload(), (byte) 0);}</pre>
     * See {@code docs/security/memory-hygiene.md} for the full rule.
     *
     * @param projectId the ID of the project the credential belongs to.
     * @param id the ID of the credential to retrieve.
     * @return the decrypted credential details; the caller takes ownership of the sensitive
     * decrypted byte array and must wipe it after processing.
     */

    @Override
    @Transactional(readOnly = true)
    public DecryptedCredential execute(UUID projectId, UUID id) {

        SecretKey key = masterKeySessionPort.getKey();

        if (key == null) throw new VaultLockedException();
        if (checkMasterPasswordSetupUseCase.isSetupRequired()) throw new MasterPasswordNotConfiguredException();
        if (!projectRepositoryPort.existsById(projectId)) throw new ResourceNotFoundException("Project", projectId);

        Credential credential = credentialRepositoryPort.findById(id)
                .filter(c -> projectId.equals(c.getProjectId()))
                .orElseThrow(() -> new ResourceNotFoundException("Credential", id));

        byte [] decryptedPayloadBytes = cryptoPort.decrypt(
                credential.getPayloadEncrypted(),
                credential.getEncryptionIv(),
                credential.getEncryptionAuthTag(),
                key
        );

        return new DecryptedCredential(
                credential.getId(),
                credential.getProjectId(),
                credential.getTitle(),
                credential.getSecretType(),
                decryptedPayloadBytes,
                credential.getNotes(),
                credential.getRelatedUrl(),
                credential.getCreatedAt(),
                credential.getUpdatedAt()
        );
    }
}
