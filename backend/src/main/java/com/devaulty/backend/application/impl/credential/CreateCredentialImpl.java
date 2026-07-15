package com.devaulty.backend.application.impl.credential;

import com.devaulty.backend.application.exception.MasterPasswordNotConfiguredException;
import com.devaulty.backend.application.exception.ResourceNotFoundException;
import com.devaulty.backend.application.exception.VaultLockedException;
import com.devaulty.backend.application.port.in.credential.CreateCredentialCommand;
import com.devaulty.backend.application.port.in.credential.CreateCredentialUseCase;
import com.devaulty.backend.application.port.in.credential.DecryptedCredential;
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
import java.util.UUID;

public class CreateCredentialImpl implements CreateCredentialUseCase {

    private final CredentialRepositoryPort credentialRepositoryPort;
    private final ProjectRepositoryPort projectRepositoryPort;
    private final CryptoPort cryptoPort;
    private final MasterKeySessionPort masterKeySessionPort;
    private final CheckMasterPasswordSetupUseCase checkMasterPasswordSetupUseCase;

    public CreateCredentialImpl(CredentialRepositoryPort credentialRepositoryPort, ProjectRepositoryPort projectRepositoryPort, CryptoPort cryptoPort, MasterKeySessionPort masterKeySessionPort, CheckMasterPasswordSetupUseCase checkMasterPasswordSetupUseCase) {
        this.credentialRepositoryPort = credentialRepositoryPort;
        this.projectRepositoryPort = projectRepositoryPort;
        this.cryptoPort = cryptoPort;
        this.masterKeySessionPort = masterKeySessionPort;
        this.checkMasterPasswordSetupUseCase = checkMasterPasswordSetupUseCase;
    }

    /**
     * Creates a new credential, encrypting its payload with the active master key.
     *
     * <p><b>⚠ MEMORY OWNERSHIP:</b> the returned {@link DecryptedCredential#decryptedPayload()}
     * is NOT zeroed here — ownership passes to the caller. Whoever consumes it MUST
     * wipe it after use:
     * <pre>{@code Arrays.fill(decryptedCredential.decryptedPayload(), (byte) 0);}</pre>
     * See {@code docs/security/memory-hygiene.md} for the full rule.
     *
     * @param command plaintext payload ({@code char[]}); this method zeroes it —
     *                caller does not need to.
     * @return credential with decrypted payload; caller owns and must zero
     *         {@code decryptedPayload} after use.
     */

    @Override
    @Transactional
    public DecryptedCredential execute(CreateCredentialCommand command) {

        // 1. Basic validations
        if (!projectRepositoryPort.existsById(command.projectId())) throw new ResourceNotFoundException("Project", command.projectId());
        if (checkMasterPasswordSetupUseCase.isSetupRequired()) throw new MasterPasswordNotConfiguredException();

        byte[] payloadBytes = null;
        byte[] decryptedBytes = null;

        try {
            // 2. Gets the key of the active on RAM
            SecretKey key = masterKeySessionPort.getKey();
            if (key == null) throw new VaultLockedException();

            // 3. Converts the payload into a byte array
            ByteBuffer buffer = StandardCharsets.UTF_8.encode(CharBuffer.wrap(command.payload()));
            payloadBytes = new byte[buffer.remaining()];
            buffer.get(payloadBytes);

            // 4. Encrypts the payload using the key and the Bouncy Castle
            CryptoResultDto cryptoResult = cryptoPort.encrypt(payloadBytes, key);

            // 5. Persists the encrypted payload into the database
            Credential credential = new Credential();
            credential.setId(UUID.randomUUID());
            credential.setProjectId(command.projectId());
            credential.setTitle(command.title());
            credential.setSecretType(command.secretType());
            credential.setPayloadEncrypted(cryptoResult.cipherText());
            credential.setEncryptionIv(cryptoResult.iv());
            credential.setEncryptionAuthTag(cryptoResult.authTag());
            credential.setNotes(command.notes());
            credential.setRelatedUrl(command.relatedUrl());
            credential.setCreatedAt(LocalDateTime.now());

            Credential savedCredential = credentialRepositoryPort.save(credential);

            // 6. Decrypts the payload
            decryptedBytes = cryptoPort.decrypt(
                    savedCredential.getPayloadEncrypted(),
                    savedCredential.getEncryptionIv(),
                    savedCredential.getEncryptionAuthTag(),
                    key)
            ;

            // 7. Returns the decrypted credentials
            return new DecryptedCredential(
                    savedCredential.getId(),
                    savedCredential.getProjectId(),
                    savedCredential.getTitle(),
                    savedCredential.getSecretType(),
                    decryptedBytes,
                    savedCredential.getNotes(),
                    savedCredential.getRelatedUrl(),
                    savedCredential.getCreatedAt(),
                    savedCredential.getUpdatedAt()
            );

        } finally {
            // 8. Clears the payload from memory
            // decryptedBytes is NOT cleared here — its ownership was transferred
            // to the returned DecryptedCredential, and it will be cleared by
            // whoever consumes it (CredentialWebMapper.jsonToMap).
            if(payloadBytes != null) Arrays.fill(payloadBytes, (byte) 0);
            Arrays.fill(command.payload(), '\0');
        }
    }
}