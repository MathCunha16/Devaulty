package com.devaulty.backend.infrastructure.configuration;

import com.devaulty.backend.application.impl.credential.*;
import com.devaulty.backend.application.port.in.credential.*;
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

    @Bean
    public GetAllCredentialsByProjectUseCase getCredentialsByProjectUseCase(
            CredentialRepositoryPort credentialRepositoryPort,
            ProjectRepositoryPort projectRepositoryPort,
            MasterKeySessionPort masterKeySessionPort,
            CheckMasterPasswordSetupUseCase checkMasterPasswordSetupUseCase
    ){
        return new GetAllCredentialsByProjectImpl(
                credentialRepositoryPort,
                projectRepositoryPort,
                masterKeySessionPort,
                checkMasterPasswordSetupUseCase
        );
    }

    @Bean
    public GetCredentialByIdUseCase getCredentialByIdUseCase(
            CredentialRepositoryPort credentialRepositoryPort,
            ProjectRepositoryPort projectRepositoryPort,
            CryptoPort cryptoPort,
            MasterKeySessionPort masterKeySessionPort,
            CheckMasterPasswordSetupUseCase checkMasterPasswordSetupUseCase
    ){
        return new GetCredentialByIdImpl(
                credentialRepositoryPort,
                projectRepositoryPort,
                cryptoPort,
                masterKeySessionPort,
                checkMasterPasswordSetupUseCase
        );
    }

    @Bean
    public UpdateCredentialUseCase updateCredentialUseCase(
            CredentialRepositoryPort credentialRepositoryPort,
            ProjectRepositoryPort projectRepositoryPort,
            CryptoPort cryptoPort,
            MasterKeySessionPort masterKeySessionPort,
            CheckMasterPasswordSetupUseCase checkMasterPasswordSetupUseCase
    ){
        return new UpdateCredentialImpl(
                credentialRepositoryPort,
                projectRepositoryPort,
                cryptoPort,
                masterKeySessionPort,
                checkMasterPasswordSetupUseCase
        );
    }

    @Bean
    public DeleteCredentialUseCase deleteCredentialUseCase(
            CredentialRepositoryPort credentialRepositoryPort,
            ProjectRepositoryPort projectRepositoryPort,
            MasterKeySessionPort masterKeySessionPort,
            CheckMasterPasswordSetupUseCase checkMasterPasswordSetupUseCase
    ){
        return new DeleteCredentialImpl(
                credentialRepositoryPort,
                projectRepositoryPort,
                masterKeySessionPort,
                checkMasterPasswordSetupUseCase
        );
    }
}