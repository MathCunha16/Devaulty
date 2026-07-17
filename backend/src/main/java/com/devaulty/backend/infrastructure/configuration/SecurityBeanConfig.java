package com.devaulty.backend.infrastructure.configuration;

import com.devaulty.backend.application.impl.security.*;
import com.devaulty.backend.application.port.in.security.*;
import com.devaulty.backend.application.port.out.security.MasterKeySessionPort;
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

    @Bean
    public GetSessionStatusUseCase getSessionStatusUseCase(
            MasterKeySessionPort sessionHolder
    ){
        return new GetSessionStatusImpl(
                sessionHolder
        );
    }
}
