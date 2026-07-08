package com.devaulty.backend.infrastructure.configuration;

import com.devaulty.backend.application.impl.project.CreateProjectImpl;
import com.devaulty.backend.application.impl.project.GetAllProjectsImpl;
import com.devaulty.backend.application.impl.project.GetProjectByIdImpl;
import com.devaulty.backend.application.port.in.project.CreateProjectUseCase;
import com.devaulty.backend.application.port.in.project.GetAllProjectsUseCase;
import com.devaulty.backend.application.port.in.project.GetProjectByIdUseCase;
import com.devaulty.backend.application.port.out.persistence.ProjectRepositoryPort;
import org.springframework.context.annotation.Bean;
import org.springframework.context.annotation.Configuration;

@Configuration
public class ProjectBeanConfig {

    @Bean
    public CreateProjectUseCase createProjectUseCase(ProjectRepositoryPort projectRepository) {
        return new CreateProjectImpl(projectRepository);
    }

    @Bean
    public GetAllProjectsUseCase getAllProjectsUseCase(ProjectRepositoryPort projectRepository) {
        return new GetAllProjectsImpl(projectRepository);
    }

    @Bean
    public GetProjectByIdUseCase getProjectByIdUseCase(ProjectRepositoryPort projectRepository) {
        return new GetProjectByIdImpl(projectRepository);
    }
}
