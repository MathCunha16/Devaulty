package com.devaulty.backend.application.impl.credential;

import com.devaulty.backend.application.exception.MasterPasswordNotConfiguredException;
import com.devaulty.backend.application.exception.ResourceNotFoundException;
import com.devaulty.backend.application.exception.VaultLockedException;
import com.devaulty.backend.application.port.in.credential.CreateCredentialCommand;
import com.devaulty.backend.application.port.in.credential.CreateCredentialUseCase;
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

    @Override
    @Transactional
    public Credential execute(CreateCredentialCommand command) {

        // 1. Basic validations
        if (!projectRepositoryPort.existsById(command.projectId())) throw new ResourceNotFoundException("Project", command.projectId());
        if (checkMasterPasswordSetupUseCase.isSetupRequired()) throw new MasterPasswordNotConfiguredException();

        byte[] payloadBytes = null;

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

            return credentialRepositoryPort.save(credential);

        } finally {
            // 6. Clears the payload from memory
            if(payloadBytes != null) {
                Arrays.fill(payloadBytes, (byte) 0);
            }
            Arrays.fill(command.payload(), '\0');
        }

    }
}
