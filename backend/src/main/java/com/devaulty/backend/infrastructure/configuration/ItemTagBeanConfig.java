package com.devaulty.backend.infrastructure.configuration;

import com.devaulty.backend.application.impl.tag.item.AssociateTagToItemImpl;
import com.devaulty.backend.application.impl.tag.item.GetTagsForItemImpl;
import com.devaulty.backend.application.impl.tag.item.GetTagsForItemsImpl;
import com.devaulty.backend.application.impl.tag.item.RemoveTagFromItemImpl;
import com.devaulty.backend.application.port.in.tag.item.AssociateTagToItemUseCase;
import com.devaulty.backend.application.port.in.tag.item.GetTagsForItemUseCase;
import com.devaulty.backend.application.port.in.tag.item.GetTagsForItemsUseCase;
import com.devaulty.backend.application.port.in.tag.item.RemoveTagFromItemUseCase;
import com.devaulty.backend.application.port.out.persistence.ItemTagRepositoryPort;
import com.devaulty.backend.application.port.out.persistence.ProjectRepositoryPort;
import com.devaulty.backend.application.port.out.persistence.ProjectScopedRepositoryPort;
import com.devaulty.backend.application.port.out.persistence.TagRepositoryPort;
import org.springframework.context.annotation.Bean;
import org.springframework.context.annotation.Configuration;

import java.util.List;

@Configuration
public class ItemTagBeanConfig {

    @Bean
    public GetTagsForItemUseCase getTagsForItemUseCase(
            ItemTagRepositoryPort itemTagRepositoryPort,
            ProjectRepositoryPort projectRepositoryPort,
            List<ProjectScopedRepositoryPort> projectScopedRepositories
    ){
        return new GetTagsForItemImpl(
                itemTagRepositoryPort,
                projectRepositoryPort,
                projectScopedRepositories
        );
    }

    @Bean
    public GetTagsForItemsUseCase getTagsForItemsUseCase(
            ItemTagRepositoryPort itemTagRepositoryPort,
            ProjectRepositoryPort projectRepositoryPort,
            List<ProjectScopedRepositoryPort> projectScopedRepositories
    ){
        return new GetTagsForItemsImpl(
                itemTagRepositoryPort,
                projectRepositoryPort,
                projectScopedRepositories
        );
    }

    @Bean
    public AssociateTagToItemUseCase associateTagToItemUseCase(
            ItemTagRepositoryPort itemTagRepositoryPort,
            TagRepositoryPort tagRepositoryPort,
            ProjectRepositoryPort projectRepositoryPort,
            List<ProjectScopedRepositoryPort> projectScopedRepositories
    ){
        return new AssociateTagToItemImpl(
                itemTagRepositoryPort,
                tagRepositoryPort,
                projectRepositoryPort,
                projectScopedRepositories
        );
    }

    @Bean
    public RemoveTagFromItemUseCase removeTagFromItemUseCase(
            ItemTagRepositoryPort itemTagRepositoryPort,
            ProjectRepositoryPort projectRepositoryPort,
            TagRepositoryPort tagRepositoryPort,
            List<ProjectScopedRepositoryPort> projectScopedRepositories
    ){
        return new RemoveTagFromItemImpl(
                itemTagRepositoryPort,
                projectRepositoryPort,
                tagRepositoryPort,
                projectScopedRepositories
        );
    }
}
