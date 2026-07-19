package com.devaulty.backend.application.impl.tag.item;

import com.devaulty.backend.application.exception.ResourceNotFoundException;
import com.devaulty.backend.application.port.in.tag.item.GetTagsForItemsUseCase;
import com.devaulty.backend.application.port.out.persistence.ItemTagRepositoryPort;
import com.devaulty.backend.application.port.out.persistence.ProjectRepositoryPort;
import com.devaulty.backend.domain.model.Tag;
import org.springframework.transaction.annotation.Transactional;

import java.util.List;
import java.util.Map;
import java.util.UUID;

public class GetTagsForItemsImpl implements GetTagsForItemsUseCase {

    private final ItemTagRepositoryPort itemTagRepositoryPort;
    private final ProjectRepositoryPort projectRepositoryPort;

    public GetTagsForItemsImpl(ItemTagRepositoryPort itemTagRepositoryPort, ProjectRepositoryPort projectRepositoryPort) {
        this.itemTagRepositoryPort = itemTagRepositoryPort;
        this.projectRepositoryPort = projectRepositoryPort;
    }

    @Override
    @Transactional(readOnly = true)
    public Map<UUID, List<Tag>> execute(String itemType, UUID projectId, List<UUID> itemIds) {

        if(!projectRepositoryPort.existsById(projectId)) throw new ResourceNotFoundException("Project", projectId);

        return itemTagRepositoryPort.findTagsForItems(itemType, itemIds);
    }
}
