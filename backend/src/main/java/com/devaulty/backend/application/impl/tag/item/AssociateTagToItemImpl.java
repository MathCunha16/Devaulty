package com.devaulty.backend.application.impl.tag.item;

import com.devaulty.backend.application.exception.ResourceNotFoundException;
import com.devaulty.backend.application.port.in.tag.item.AssociateTagToItemUseCase;
import com.devaulty.backend.application.port.out.persistence.ItemTagRepositoryPort;
import com.devaulty.backend.application.port.out.persistence.ProjectRepositoryPort;
import com.devaulty.backend.application.port.out.persistence.ProjectScopedRepositoryPort;
import com.devaulty.backend.application.port.out.persistence.TagRepositoryPort;
import org.springframework.transaction.annotation.Transactional;

import java.util.List;
import java.util.UUID;

public class AssociateTagToItemImpl implements AssociateTagToItemUseCase {

    private final ItemTagRepositoryPort itemTagRepositoryPort;
    private final TagRepositoryPort tagRepositoryPort;
    private final ProjectRepositoryPort projectRepositoryPort;
    private final List<ProjectScopedRepositoryPort> projectScopedRepositories;

    public AssociateTagToItemImpl(ItemTagRepositoryPort itemTagRepositoryPort,
                                  TagRepositoryPort tagRepositoryPort,
                                  ProjectRepositoryPort projectRepositoryPort,
                                  List<ProjectScopedRepositoryPort> projectScopedRepositories) {
        this.itemTagRepositoryPort = itemTagRepositoryPort;
        this.tagRepositoryPort = tagRepositoryPort;
        this.projectRepositoryPort = projectRepositoryPort;
        this.projectScopedRepositories = projectScopedRepositories;
    }

    @Override
    @Transactional
    public void execute(UUID projectId, String itemType, UUID itemId, UUID tagId) {

        if (!projectRepositoryPort.existsById(projectId)) {
            throw new ResourceNotFoundException("Project", projectId);
        }
        if (!tagRepositoryPort.existsByIdAndProjectId(tagId, projectId)) {
            throw new ResourceNotFoundException("Tag", tagId);
        }

        validateItemOwnership(itemType, projectId, itemId);

        itemTagRepositoryPort.associateTagToItem(tagId, itemType, itemId);
    }

    private void validateItemOwnership(String itemType, UUID projectId, UUID itemId) {
        boolean exists = projectScopedRepositories.stream()
                .filter(repo -> repo.getSupportedType().equalsIgnoreCase(itemType))
                .findFirst()
                .map(repo -> repo.existsByIdAndProjectId(itemId, projectId))
                .orElseThrow(() -> new IllegalArgumentException("Unsupported item type: " + itemType));

        if (!exists) {
            String resourceName = itemType.substring(0, 1).toUpperCase() + itemType.substring(1).toLowerCase();
            throw new ResourceNotFoundException(resourceName, itemId);
        }
    }
}
