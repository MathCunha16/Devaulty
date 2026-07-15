package com.devaulty.backend.application.impl.security;

import com.devaulty.backend.application.exception.InvalidMasterPasswordException;
import com.devaulty.backend.application.port.out.security.MasterKeySessionPort;
import com.devaulty.backend.application.exception.MasterPasswordNotConfiguredException;
import com.devaulty.backend.application.port.in.security.UnlockVaultUseCase;
import com.devaulty.backend.application.port.out.persistence.AppSettingRepositoryPort;
import com.devaulty.backend.application.port.out.security.KeyDerivationPort;
import com.devaulty.backend.domain.model.AppSetting;
import org.springframework.security.crypto.argon2.Argon2PasswordEncoder;
import org.springframework.transaction.annotation.Transactional;

import javax.crypto.SecretKey;
import java.nio.CharBuffer;
import java.util.Arrays;
import java.util.Base64;

public class UnlockVaultImpl implements UnlockVaultUseCase {

    private final AppSettingRepositoryPort appSettingRepositoryPort;
    private final MasterKeySessionPort sessionHolder;
    private final Argon2PasswordEncoder argon2PasswordEncoder;
    private final KeyDerivationPort keyDerivationPort;

    private static final String MASTER_PASSWORD_HASH_KEY = "master_password_hash";
    private static final String MASTER_PASSWORD_SALT_KEY = "master_password_salt";

    public UnlockVaultImpl(AppSettingRepositoryPort appSettingRepositoryPort, MasterKeySessionPort sessionHolder, KeyDerivationPort keyDerivationPort) {
        this.appSettingRepositoryPort = appSettingRepositoryPort;
        this.sessionHolder = sessionHolder;
        this.keyDerivationPort = keyDerivationPort;
        this.argon2PasswordEncoder = new Argon2PasswordEncoder(16, 32, 2, 65536, 3);
    }

    @Override
    @Transactional(readOnly = true)
    public boolean execute(char[] password) {
        try {
            // 1. Get hash from DB
            String hash = appSettingRepositoryPort.findByKey(MASTER_PASSWORD_HASH_KEY)
                    .map(AppSetting::getValue)
                    .orElseThrow(MasterPasswordNotConfiguredException::new);

            // 2. Validates the password against the hash
            if (argon2PasswordEncoder.matches(CharBuffer.wrap(password), hash)) {

                // 3. Gets the salt from DB for key derivation
                String saltBase64 = appSettingRepositoryPort.findByKey(MASTER_PASSWORD_SALT_KEY)
                        .map(AppSetting::getValue)
                        .orElseThrow(MasterPasswordNotConfiguredException::new);

                byte[] saltBytes = Base64.getDecoder().decode(saltBase64);

                // 4. Derives the AES-GCM key once and populates the memory session holder immediately
                SecretKey secretKey = keyDerivationPort.deriveKey(password, saltBytes);
                sessionHolder.setKey(secretKey);
                return true;
            }

            throw  new InvalidMasterPasswordException("Invalid MasterPassword!");

        } finally {
            // CRITICAL PROTECTION: Overwrites the input char array immediately after the use case completes
            Arrays.fill(password, '\0');
        }

    }
}
