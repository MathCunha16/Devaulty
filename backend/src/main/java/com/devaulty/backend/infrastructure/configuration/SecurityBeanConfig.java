package com.devaulty.backend.infrastructure.configuration;

import com.devaulty.backend.application.port.out.security.MasterKeySessionPort;
import com.devaulty.backend.application.impl.security.CheckMasterPasswordSetupImpl;
import com.devaulty.backend.application.impl.security.LockVaultImpl;
import com.devaulty.backend.application.impl.security.SetupMasterPasswordImpl;
import com.devaulty.backend.application.impl.security.UnlockVaultImpl;
import com.devaulty.backend.application.port.in.security.CheckMasterPasswordSetupUseCase;
import com.devaulty.backend.application.port.in.security.LockVaultUseCase;
import com.devaulty.backend.application.port.in.security.SetupMasterPasswordUseCase;
import com.devaulty.backend.application.port.in.security.UnlockVaultUseCase;
import com.devaulty.backend.application.port.out.persistence.AppSettingRepositoryPort;
import com.devaulty.backend.application.port.out.security.KeyDerivationPort;
import org.springframework.context.annotation.Bean;
import org.springframework.context.annotation.Configuration;

@Configuration
public class SecurityBeanConfig {

    @Bean
    public CheckMasterPasswordSetupUseCase checkMasterPasswordSetupUseCase(AppSettingRepositoryPort appSettingRepositoryPort){
        return new CheckMasterPasswordSetupImpl(appSettingRepositoryPort);
    }

    @Bean
    public LockVaultUseCase lockVaultUseCase(MasterKeySessionPort masterKeySessionPort){
        return new LockVaultImpl(masterKeySessionPort);
    }

    @Bean
    public SetupMasterPasswordUseCase setupMasterPasswordUseCase(
            AppSettingRepositoryPort appSettingRepositoryPort,
            MasterKeySessionPort masterKeySessionPort,
            KeyDerivationPort keyDerivationPort
    ){
        return new SetupMasterPasswordImpl(
                appSettingRepositoryPort,
                masterKeySessionPort,
                keyDerivationPort
        );
    }

    @Bean
    public UnlockVaultUseCase unlockVaultUseCase(
            AppSettingRepositoryPort appSettingRepositoryPort,
            MasterKeySessionPort sessionHolder,
            KeyDerivationPort keyDerivationPort
    ){
        return new UnlockVaultImpl(
                appSettingRepositoryPort,
                sessionHolder,
                keyDerivationPort
        );
    }
}
