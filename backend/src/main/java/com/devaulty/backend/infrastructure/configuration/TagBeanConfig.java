package com.devaulty.backend.infrastructure.configuration;

import com.devaulty.backend.application.impl.tag.*;
import com.devaulty.backend.application.port.in.tag.*;
import com.devaulty.backend.application.port.out.persistence.ProjectRepositoryPort;
import com.devaulty.backend.application.port.out.persistence.TagRepositoryPort;
import org.springframework.context.annotation.Bean;
import org.springframework.context.annotation.Configuration;

@Configuration
public class TagBeanConfig {

    @Bean
    public CreateTagUseCase createTagUseCase(
            TagRepositoryPort tagRepositoryPort,
            ProjectRepositoryPort projectRepositoryPort) {
        return new CreateTagImpl(
                tagRepositoryPort,
                projectRepositoryPort
        );
    }

    @Bean
    public GetTagByIdUseCase getTagByIdUseCase(
            TagRepositoryPort tagRepositoryPort,
            ProjectRepositoryPort projectRepositoryPort
    ) {
        return new GetTagByIdImpl(
                tagRepositoryPort,
                projectRepositoryPort
        );
    }

    @Bean
    public SearchTagByNameUseCase searchTagByNameUseCase(
            TagRepositoryPort tagRepositoryPort,
            ProjectRepositoryPort projectRepositoryPort
    ) {
        return new SearchTagByNameImpl(
                tagRepositoryPort,
                projectRepositoryPort
        );
    }

    @Bean
    public GetAllTagsByProjectUseCase getAllTagsByProjectUseCase(
            TagRepositoryPort tagRepositoryPort,
            ProjectRepositoryPort projectRepositoryPort
    ){
        return new GetAllTagsByProjectImpl(
                tagRepositoryPort,
                projectRepositoryPort
        );
    }

    @Bean
    public UpdateTagUseCase updateTagUseCase(
            TagRepositoryPort tagRepositoryPort,
            ProjectRepositoryPort projectRepositoryPort
    ){
        return new UpdateTagImpl(
                tagRepositoryPort,
                projectRepositoryPort
        );
    }

    @Bean
    public DeleteTagUseCase deleteTagUseCase(
            TagRepositoryPort tagRepositoryPort,
            ProjectRepositoryPort projectRepositoryPort
    ){
        return new DeleteTagImpl(
                tagRepositoryPort,
                projectRepositoryPort
        );
    }
}
