package com.devaulty.backend.application.impl.security;

import com.devaulty.backend.application.port.out.security.MasterKeySessionPort;
import com.devaulty.backend.application.exception.MasterPasswordAlreadyConfiguredException;
import com.devaulty.backend.application.port.in.security.SetupMasterPasswordUseCase;
import com.devaulty.backend.application.port.out.persistence.AppSettingRepositoryPort;
import com.devaulty.backend.application.port.out.security.KeyDerivationPort;
import com.devaulty.backend.domain.model.AppSetting;
import org.springframework.security.crypto.argon2.Argon2PasswordEncoder;
import org.springframework.transaction.annotation.Transactional;
import org.springframework.transaction.support.TransactionSynchronization;
import org.springframework.transaction.support.TransactionSynchronizationManager;

import javax.crypto.SecretKey;
import java.nio.CharBuffer;
import java.security.SecureRandom;
import java.util.Arrays;
import java.util.Base64;

public class SetupMasterPasswordImpl implements SetupMasterPasswordUseCase {

    private final AppSettingRepositoryPort appSettingRepositoryPort;
    private final MasterKeySessionPort sessionHolder;
    private final Argon2PasswordEncoder argon2PasswordEncoder;
    private final KeyDerivationPort keyDerivationPort;

    private static final String MASTER_PASSWORD_HASH_KEY = "master_password_hash";
    private static final String MASTER_PASSWORD_SALT_KEY = "master_password_salt";

    public SetupMasterPasswordImpl(AppSettingRepositoryPort appSettingRepositoryPort, MasterKeySessionPort sessionHolder, KeyDerivationPort keyDerivationPort) {
        this.appSettingRepositoryPort = appSettingRepositoryPort;
        this.sessionHolder = sessionHolder;
        this.keyDerivationPort = keyDerivationPort;
        // Instantiates the default Spring Security encoder for the login verification hash
        this.argon2PasswordEncoder = new Argon2PasswordEncoder(16, 32, 2, 65536, 3);
    }

    @Override
    @Transactional
    public void execute(char[] password) {
        try {
            // Prevents overwriting if already configured
            if(appSettingRepositoryPort.existsByKey(MASTER_PASSWORD_HASH_KEY)) throw new MasterPasswordAlreadyConfiguredException();

            // 1. Generates a strong, random 16-byte global application salt
            byte[] saltBytes = new byte[16];
            new SecureRandom().nextBytes(saltBytes);
            String saltBase64 = Base64.getEncoder().encodeToString(saltBytes);

            // 2. Generates the quick-match authentication hash for future validation
            String hash = argon2PasswordEncoder.encode(CharBuffer.wrap(password));

            // 3. Persists both elements inside the app_settings table via SQLite
            appSettingRepositoryPort.save(new AppSetting(MASTER_PASSWORD_SALT_KEY, saltBase64));
            appSettingRepositoryPort.save(new AppSetting(MASTER_PASSWORD_HASH_KEY, hash));

            // 4. Derives the AES-GCM key once and populates the memory session holder immediately
            SecretKey secretKey = keyDerivationPort.deriveKey(password, saltBytes);
            Arrays.fill(saltBytes, (byte) 0);

            // 5. Publishes the key only after the transaction commits successfully if active, otherwise publishes immediately
            if (TransactionSynchronizationManager.isSynchronizationActive()) {
                TransactionSynchronizationManager.registerSynchronization(new TransactionSynchronization() {
                    @Override
                    public void afterCommit() {
                        sessionHolder.setKey(secretKey);
                    }
                });
            } else {
                sessionHolder.setKey(secretKey);
            }

        } finally {
            // CRITICAL PROTECTION: Overwrites the input char array immediately after the use case completes
            Arrays.fill(password, '\0');
        }
    }

}
