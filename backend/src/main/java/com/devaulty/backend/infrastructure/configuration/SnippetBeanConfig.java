package com.devaulty.backend.infrastructure.configuration;

import com.devaulty.backend.application.impl.snippet.*;
import com.devaulty.backend.application.port.in.snippet.*;
import com.devaulty.backend.application.port.out.persistence.ItemTagRepositoryPort;
import com.devaulty.backend.application.port.out.persistence.ProjectRepositoryPort;
import com.devaulty.backend.application.port.out.persistence.SnippetRepositoryPort;
import org.springframework.context.annotation.Bean;
import org.springframework.context.annotation.Configuration;

@Configuration
public class SnippetBeanConfig {

    @Bean
    public CreateSnippetUseCase createSnippetUseCase(
            SnippetRepositoryPort snippetRepositoryPort,
            ProjectRepositoryPort projectRepositoryPort
    ) {
        return new CreateSnippetImpl(
                snippetRepositoryPort,
                projectRepositoryPort
        );
    }

    @Bean
    public GetAllSnippetsByProjectUseCase getAllSnippetsByProjectUseCase(
            SnippetRepositoryPort snippetRepositoryPort,
            ProjectRepositoryPort projectRepositoryPort
    ){
        return new GetAllSnippetsByProjectImpl(
                snippetRepositoryPort,
                projectRepositoryPort
        );
    }

    @Bean
    public GetSnippetByIdUseCase getSnippetByIdUseCase(
            SnippetRepositoryPort snippetRepositoryPort,
            ProjectRepositoryPort projectRepositoryPort
    ){
        return new GetSnippetByIdImpl(
                snippetRepositoryPort,
                projectRepositoryPort
        );
    }

    @Bean
    public UpdateSnippetUseCase updateSnippetUseCase(
            SnippetRepositoryPort snippetRepositoryPort,
            ProjectRepositoryPort projectRepositoryPort
    ){
        return new UpdateSnippetImpl(
                snippetRepositoryPort,
                projectRepositoryPort
        );
    }

    @Bean
    public DeleteSnippetUseCase deleteSnippetUseCase(
            SnippetRepositoryPort snippetRepositoryPort,
            ProjectRepositoryPort projectRepositoryPort,
            ItemTagRepositoryPort itemTagRepositoryPort
    ){
        return new DeleteSnippetImpl(
                snippetRepositoryPort,
                projectRepositoryPort,
                itemTagRepositoryPort
        );
    }
}
