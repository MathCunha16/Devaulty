package com.devaulty.backend.infrastructure.configuration;

import com.devaulty.backend.application.impl.problem.*;
import com.devaulty.backend.application.port.in.problem.*;
import com.devaulty.backend.application.port.out.persistence.ProblemRepositoryPort;
import com.devaulty.backend.application.port.out.persistence.ProjectRepositoryPort;
import org.springframework.context.annotation.Bean;
import org.springframework.context.annotation.Configuration;

@Configuration
public class ProblemBeanConfig {

    @Bean
    public CreateProblemUseCase createProblemUseCase(
            ProblemRepositoryPort problemRepositoryPort,
            ProjectRepositoryPort projectRepositoryPort
    ){
        return new CreateProblemImpl(
                problemRepositoryPort,
                projectRepositoryPort
        );
    }

    @Bean
    public GetAllProblemsByProjectUseCase getAllProblemsByProjectUseCase(
            ProblemRepositoryPort problemRepositoryPort,
            ProjectRepositoryPort projectRepositoryPort
    ){
        return new GetAllProblemsByProjectImpl(
                problemRepositoryPort,
                projectRepositoryPort
        );
    }

    @Bean
    public GetProblemByIdUseCase getProblemByIdUseCase(
            ProblemRepositoryPort problemRepositoryPort,
            ProjectRepositoryPort projectRepositoryPort
    ){
        return new GetProblemByIdImpl(
                problemRepositoryPort,
                projectRepositoryPort
        );
    }

    @Bean
    public DeleteProblemUseCase deleteProblemUseCase(
            ProblemRepositoryPort problemRepositoryPort,
            ProjectRepositoryPort projectRepositoryPort
    ){
        return new DeleteProblemImpl(
                problemRepositoryPort,
                projectRepositoryPort
        );
    }

    @Bean
    public UpdateProblemUseCase updateProblemUseCase(
            ProblemRepositoryPort problemRepositoryPort,
            ProjectRepositoryPort projectRepositoryPort
    ){
        return new UpdateProblemImpl(
                problemRepositoryPort,
                projectRepositoryPort
        );
    }

    @Bean
    public UpdateProblemStatusUseCase updateProblemStatusUseCase(
            ProblemRepositoryPort problemRepositoryPort,
            ProjectRepositoryPort projectRepositoryPort
    ){
        return new UpdateProblemStatusImpl(
                problemRepositoryPort,
                projectRepositoryPort
        );
    }
}
