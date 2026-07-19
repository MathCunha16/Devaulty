package com.devaulty.backend.application.impl.tag.item;

import com.devaulty.backend.application.exception.ResourceNotFoundException;
import com.devaulty.backend.application.port.in.tag.item.AssociateTagToItemUseCase;
import com.devaulty.backend.application.port.out.persistence.ItemTagRepositoryPort;
import com.devaulty.backend.application.port.out.persistence.ProjectRepositoryPort;
import com.devaulty.backend.application.port.out.persistence.TagRepositoryPort;
import org.springframework.transaction.annotation.Transactional;

import java.util.UUID;

public class AssociateTagToItemImpl implements AssociateTagToItemUseCase {

    private final ItemTagRepositoryPort itemTagRepositoryPort;
    private final TagRepositoryPort tagRepositoryPort;
    private final ProjectRepositoryPort projectRepositoryPort;

    public AssociateTagToItemImpl(ItemTagRepositoryPort itemTagRepositoryPort, TagRepositoryPort tagRepositoryPort, ProjectRepositoryPort projectRepositoryPort) {
        this.itemTagRepositoryPort = itemTagRepositoryPort;
        this.tagRepositoryPort = tagRepositoryPort;
        this.projectRepositoryPort = projectRepositoryPort;
    }

    @Override
    @Transactional
    public void execute(UUID projectId, String itemType, UUID itemId, UUID tagId) {

        if(!projectRepositoryPort.existsById(projectId)) throw new ResourceNotFoundException("Project", projectId);
        if(!tagRepositoryPort.existsByIdAndProjectId(tagId, projectId)) throw new ResourceNotFoundException("Tag", tagId);

        itemTagRepositoryPort.associateTagToItem(tagId, itemType, itemId);
    }
}
