package com.devaulty.backend.infrastructure.configuration;

import com.devaulty.backend.application.impl.project.*;
import com.devaulty.backend.application.port.in.project.*;
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

    @Bean
    public ArchiveProjectUseCase archiveProjectUseCase(ProjectRepositoryPort projectRepository) {
        return new ArchiveProjectImpl(projectRepository);
    }

    @Bean
    public UnarchiveProjectUseCase unarchiveProjectUseCase(ProjectRepositoryPort projectRepository) {
        return new UnarchiveProjectImpl(projectRepository);
    }

    @Bean
    public DeleteProjectUseCase deleteProjectUseCase(ProjectRepositoryPort projectRepository) {
        return new DeleteProjectImpl(projectRepository);
    }

    @Bean
    public UpdateProjectUseCase updateProjectUseCase(ProjectRepositoryPort projectRepository) {
        return new UpdateProjectImpl(projectRepository);
    }
}
