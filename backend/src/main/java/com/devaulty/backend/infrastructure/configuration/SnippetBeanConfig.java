package com.devaulty.backend.infrastructure.configuration;

import com.devaulty.backend.application.impl.snippet.CreateSnippetImpl;
import com.devaulty.backend.application.impl.snippet.GetAllSnippetsByProjectImpl;
import com.devaulty.backend.application.impl.snippet.GetSnippetByIdImpl;
import com.devaulty.backend.application.port.in.snippet.CreateSnippetUseCase;
import com.devaulty.backend.application.port.in.snippet.GetAllSnippetsByProjectUseCase;
import com.devaulty.backend.application.port.in.snippet.GetSnippetByIdUseCase;
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
}
