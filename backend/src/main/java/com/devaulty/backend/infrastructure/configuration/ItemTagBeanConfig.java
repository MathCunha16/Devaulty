package com.devaulty.backend.infrastructure.configuration;

import com.devaulty.backend.application.impl.tag.item.GetTagsForItemImpl;
import com.devaulty.backend.application.impl.tag.item.GetTagsForItemsImpl;
import com.devaulty.backend.application.port.in.tag.item.GetTagsForItemUseCase;
import com.devaulty.backend.application.port.in.tag.item.GetTagsForItemsUseCase;
import com.devaulty.backend.application.port.out.persistence.ItemTagRepositoryPort;
import com.devaulty.backend.application.port.out.persistence.ProjectRepositoryPort;
import org.springframework.context.annotation.Bean;
import org.springframework.context.annotation.Configuration;

@Configuration
public class ItemTagBeanConfig {

    @Bean
    public GetTagsForItemUseCase getTagsForItemUseCase(
            ItemTagRepositoryPort itemTagRepositoryPort,
            ProjectRepositoryPort projectRepositoryPort
    ){
        return new GetTagsForItemImpl(
                itemTagRepositoryPort,
                projectRepositoryPort
        );
    }

    @Bean
    public GetTagsForItemsUseCase getTagsForItemsUseCase(
            ItemTagRepositoryPort itemTagRepositoryPort,
            ProjectRepositoryPort projectRepositoryPort
    ){
        return new GetTagsForItemsImpl(
                itemTagRepositoryPort,
                projectRepositoryPort
        );
    }
}
