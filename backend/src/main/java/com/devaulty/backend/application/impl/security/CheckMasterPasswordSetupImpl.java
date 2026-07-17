package com.devaulty.backend.application.impl.security;

import com.devaulty.backend.application.port.in.security.CheckMasterPasswordSetupUseCase;
import com.devaulty.backend.application.port.out.persistence.AppSettingRepositoryPort;
import org.springframework.transaction.annotation.Transactional;

public class CheckMasterPasswordSetupImpl implements CheckMasterPasswordSetupUseCase {

    private final AppSettingRepositoryPort appSettingRepositoryPort;

    private static final String MASTER_PASSWORD_HASH_KEY = "master_password_hash";

    public CheckMasterPasswordSetupImpl(AppSettingRepositoryPort appSettingRepositoryPort) {
        this.appSettingRepositoryPort = appSettingRepositoryPort;
    }

    @Override
    @Transactional(readOnly = true)
    public boolean isSetupRequired() {
        return !appSettingRepositoryPort.existsByKey(MASTER_PASSWORD_HASH_KEY);
    }
}
