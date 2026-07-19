package com.devaulty.backend.application.impl.tag.item;

import com.devaulty.backend.application.exception.ResourceNotFoundException;
import com.devaulty.backend.application.port.in.tag.item.RemoveTagFromItemUseCase;
import com.devaulty.backend.application.port.out.persistence.ItemTagRepositoryPort;
import com.devaulty.backend.application.port.out.persistence.ProjectRepositoryPort;
import org.springframework.transaction.annotation.Transactional;

import java.util.UUID;

public class RemoveTagFromItemImpl implements RemoveTagFromItemUseCase {

    private final ItemTagRepositoryPort itemTagRepositoryPort;
    private final ProjectRepositoryPort projectRepositoryPort;

    public RemoveTagFromItemImpl(ItemTagRepositoryPort itemTagRepositoryPort, ProjectRepositoryPort projectRepositoryPort) {
        this.itemTagRepositoryPort = itemTagRepositoryPort;
        this.projectRepositoryPort = projectRepositoryPort;
    }

    @Override
    @Transactional
    public void execute(UUID projectId, String itemType, UUID itemId, UUID tagId) {
        if(!projectRepositoryPort.existsById(projectId)) throw new ResourceNotFoundException("Project", projectId);

        itemTagRepositoryPort.disassembleTagFromItem(tagId, itemType, itemId);
    }
}
