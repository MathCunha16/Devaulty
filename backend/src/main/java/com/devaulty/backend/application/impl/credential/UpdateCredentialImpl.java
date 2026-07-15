package com.devaulty.backend.application.impl.credential;

import com.devaulty.backend.application.exception.MasterPasswordNotConfiguredException;
import com.devaulty.backend.application.exception.ResourceNotFoundException;
import com.devaulty.backend.application.exception.VaultLockedException;
import com.devaulty.backend.application.port.in.credential.DecryptedCredential;
import com.devaulty.backend.application.port.in.credential.UpdateCredentialCommand;
import com.devaulty.backend.application.port.in.credential.UpdateCredentialUseCase;
import com.devaulty.backend.application.port.in.security.CheckMasterPasswordSetupUseCase;
import com.devaulty.backend.application.port.out.persistence.CredentialRepositoryPort;
import com.devaulty.backend.application.port.out.persistence.ProjectRepositoryPort;
import com.devaulty.backend.application.port.out.security.CryptoPort;
import com.devaulty.backend.application.port.out.security.MasterKeySessionPort;
import com.devaulty.backend.application.port.out.security.dto.CryptoResultDto;
import com.devaulty.backend.domain.model.Credential;
import org.springframework.transaction.annotation.Transactional;

import javax.crypto.SecretKey;
import java.nio.ByteBuffer;
import java.nio.CharBuffer;
import java.nio.charset.StandardCharsets;
import java.time.LocalDateTime;
import java.util.Arrays;

public class UpdateCredentialImpl implements UpdateCredentialUseCase {

    private final CredentialRepositoryPort credentialRepositoryPort;
    private final ProjectRepositoryPort projectRepositoryPort;
    private final CryptoPort cryptoPort;
    private final MasterKeySessionPort masterKeySessionPort;
    private final CheckMasterPasswordSetupUseCase checkMasterPasswordSetupUseCase;

    public UpdateCredentialImpl(CredentialRepositoryPort credentialRepositoryPort, ProjectRepositoryPort projectRepositoryPort, CryptoPort cryptoPort, MasterKeySessionPort masterKeySessionPort, CheckMasterPasswordSetupUseCase checkMasterPasswordSetupUseCase) {
        this.credentialRepositoryPort = credentialRepositoryPort;
        this.projectRepositoryPort = projectRepositoryPort;
        this.cryptoPort = cryptoPort;
        this.masterKeySessionPort = masterKeySessionPort;
        this.checkMasterPasswordSetupUseCase = checkMasterPasswordSetupUseCase;
    }

    /**
     * Updates an existing credential, re-encrypting the new payload if provided.
     *
     * <p><b>⚠ MEMORY OWNERSHIP:</b> the returned {@link DecryptedCredential#decryptedPayload()}
     * is NOT zeroed here — ownership passes to the caller. Whoever consumes it MUST
     * wipe it after use:
     * <pre>{@code Arrays.fill(decryptedCredential.decryptedPayload(), (byte) 0);}</pre>
     * See {@code docs/security/memory-hygiene.md} for the full rule.
     *
     * @param command update details containing the target IDs, metadata, and optional new plaintext
     * payload ({@code char[]}); this method zeroes the input payload if present —
     * caller does not need to.
     * @return the updated credential with decrypted payload; caller owns and must zero
     * {@code decryptedPayload} after use.
     */

    @Override
    @Transactional
    public DecryptedCredential execute(UpdateCredentialCommand command) {
        SecretKey key = masterKeySessionPort.getKey();

        if (key == null) throw new VaultLockedException();
        if (checkMasterPasswordSetupUseCase.isSetupRequired()) throw new MasterPasswordNotConfiguredException();
        if (!projectRepositoryPort.existsById(command.projectId()))
            throw new ResourceNotFoundException("Project", command.projectId());

        Credential credential = credentialRepositoryPort.findById(command.id())
                .filter(c -> command.projectId().equals(c.getProjectId()))
                .orElseThrow(() -> new ResourceNotFoundException("Credential", command.id()));

        byte[] payloadBytes = null;
        byte[] decryptedBytes = null;

        try {
            if (command.title() != null) credential.setTitle(command.title());
            if (command.secretType() != null) credential.setSecretType(command.secretType());
            if (command.notes() != null) credential.setNotes(command.notes());
            if (command.relatedUrl() != null) credential.setRelatedUrl(command.relatedUrl());
            if (command.payload() != null) {

                ByteBuffer buffer = StandardCharsets.UTF_8.encode(CharBuffer.wrap(command.payload()));
                payloadBytes = new byte[buffer.remaining()];
                buffer.get(payloadBytes);

                CryptoResultDto cryptoResult = cryptoPort.encrypt(payloadBytes, key);

                credential.setPayloadEncrypted(cryptoResult.cipherText());
                credential.setEncryptionIv(cryptoResult.iv());
                credential.setEncryptionAuthTag(cryptoResult.authTag());
            }

            credential.setUpdatedAt(LocalDateTime.now());
            credentialRepositoryPort.save(credential);

            decryptedBytes = cryptoPort.decrypt(
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
                    decryptedBytes,
                    credential.getNotes(),
                    credential.getRelatedUrl(),
                    credential.getCreatedAt(),
                    credential.getUpdatedAt()
            );

        } finally {
            // decryptedBytes is NOT cleared here — its ownership was transferred
            // to the returned DecryptedCredential, and it will be cleared by
            // whoever consumes it (CredentialWebMapper.jsonToMap).
            if(payloadBytes != null) Arrays.fill(payloadBytes, (byte) 0);
            if(command.payload() != null) Arrays.fill(command.payload(), '\0');
        }
    }
}
