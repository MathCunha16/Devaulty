package com.devaulty.backend.infrastructure.configuration;

import com.devaulty.backend.application.impl.link.*;
import com.devaulty.backend.application.port.in.link.*;
import com.devaulty.backend.application.port.out.persistence.ItemTagRepositoryPort;
import com.devaulty.backend.application.port.out.persistence.LinkRepositoryPort;
import com.devaulty.backend.application.port.out.persistence.ProjectRepositoryPort;
import org.springframework.context.annotation.Bean;
import org.springframework.context.annotation.Configuration;

@Configuration
public class LinkBeanConfig {

    @Bean
    public CreateLinkUseCase createLinkUseCase (
            LinkRepositoryPort linkRepositoryPort,
            ProjectRepositoryPort projectRepositoryPort
    ){
        return new CreateLinkImpl(
                linkRepositoryPort,
                projectRepositoryPort
        );
    }

    @Bean
    public GetAllLinksByProjectUseCase getAllLinksByProjectUseCase(
            LinkRepositoryPort linkRepositoryPort,
            ProjectRepositoryPort projectRepositoryPort
    ){
        return new GetAllLinksByProjectImpl(
                linkRepositoryPort,
                projectRepositoryPort
        );
    }

    @Bean
    public GetLinkByIdUseCase getLinkByIdUseCase(
            LinkRepositoryPort linkRepositoryPort,
            ProjectRepositoryPort projectRepositoryPort
    ){
        return new GetLinkByIdImpl(
                linkRepositoryPort,
                projectRepositoryPort
        );
    }

    @Bean
    public DeleteLinkUseCase deleteLinkUseCase(
            LinkRepositoryPort linkRepositoryPort,
            ProjectRepositoryPort projectRepositoryPort,
            ItemTagRepositoryPort itemTagRepositoryPort
    ){
        return new DeleteLinkImpl(
                linkRepositoryPort,
                projectRepositoryPort,
                itemTagRepositoryPort
        );
    }

    @Bean
    public UpdateLinkUseCase updateLinkUseCase(
            LinkRepositoryPort linkRepositoryPort,
            ProjectRepositoryPort projectRepositoryPort
    ){
        return new UpdateLinkImpl(
                linkRepositoryPort,
                projectRepositoryPort
        );
    }
}
