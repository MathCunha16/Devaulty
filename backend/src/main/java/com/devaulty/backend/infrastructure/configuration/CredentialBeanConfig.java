package com.devaulty.backend.infrastructure.configuration;

import com.devaulty.backend.application.impl.credential.CreateCredentialImpl;
import com.devaulty.backend.application.port.in.credential.CreateCredentialUseCase;
import com.devaulty.backend.application.port.in.security.CheckMasterPasswordSetupUseCase;
import com.devaulty.backend.application.port.out.persistence.CredentialRepositoryPort;
import com.devaulty.backend.application.port.out.persistence.ProjectRepositoryPort;
import com.devaulty.backend.application.port.out.security.CryptoPort;
import com.devaulty.backend.application.port.out.security.MasterKeySessionPort;
import org.springframework.context.annotation.Bean;
import org.springframework.context.annotation.Configuration;

@Configuration
public class CredentialBeanConfig {

    @Bean
    public CreateCredentialUseCase createCredentialUseCase(
            CredentialRepositoryPort credentialRepositoryPort,
            ProjectRepositoryPort projectRepositoryPort,
            CryptoPort cryptoPort,
            MasterKeySessionPort masterKeySessionPort,
            CheckMasterPasswordSetupUseCase checkMasterPasswordSetupUseCase
    ){
        return new CreateCredentialImpl(
                credentialRepositoryPort,
                projectRepositoryPort,
                cryptoPort,
                masterKeySessionPort,
                checkMasterPasswordSetupUseCase
        );
    }
}
