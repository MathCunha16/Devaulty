package com.devaulty.backend.application.impl.tag.item;

import com.devaulty.backend.application.exception.ResourceNotFoundException;
import com.devaulty.backend.application.port.in.tag.item.GetTagsForItemsUseCase;
import com.devaulty.backend.application.port.out.persistence.ItemTagRepositoryPort;
import com.devaulty.backend.application.port.out.persistence.ProjectRepositoryPort;
import com.devaulty.backend.application.port.out.persistence.ProjectScopedRepositoryPort;
import com.devaulty.backend.domain.model.Tag;
import org.springframework.transaction.annotation.Transactional;

import java.util.List;
import java.util.Map;
import java.util.UUID;

public class GetTagsForItemsImpl implements GetTagsForItemsUseCase {

    private final ItemTagRepositoryPort itemTagRepositoryPort;
    private final ProjectRepositoryPort projectRepositoryPort;
    private final List<ProjectScopedRepositoryPort> projectScopedRepositories;

    public GetTagsForItemsImpl(ItemTagRepositoryPort itemTagRepositoryPort,
                               ProjectRepositoryPort projectRepositoryPort,
                               List<ProjectScopedRepositoryPort> projectScopedRepositories) {
        this.itemTagRepositoryPort = itemTagRepositoryPort;
        this.projectRepositoryPort = projectRepositoryPort;
        this.projectScopedRepositories = projectScopedRepositories;
    }

    @Override
    @Transactional(readOnly = true)
    public Map<UUID, List<Tag>> execute(String itemType, UUID projectId, List<UUID> itemIds) {

        if (!projectRepositoryPort.existsById(projectId)) {
            throw new ResourceNotFoundException("Project", projectId);
        }

        validateItemsOwnership(itemType, projectId, itemIds);

        return itemTagRepositoryPort.findTagsForItems(itemType, projectId, itemIds);
    }

    private void validateItemsOwnership(String itemType, UUID projectId, List<UUID> itemIds) {
        ProjectScopedRepositoryPort repo = projectScopedRepositories.stream()
                .filter(r -> r.getSupportedType().equalsIgnoreCase(itemType))
                .findFirst()
                .orElseThrow(() -> new IllegalArgumentException("Unsupported item type: " + itemType));

        List<UUID> existingIds = repo.findExistingIdsByProject(itemIds, projectId);

        if (existingIds.size() != itemIds.size()) {
            UUID missingId = itemIds.stream()
                    .filter(id -> !existingIds.contains(id))
                    .findFirst()
                    .orElse(null);

            String resourceName = itemType.substring(0, 1).toUpperCase() + itemType.substring(1).toLowerCase();
            throw new ResourceNotFoundException(resourceName, missingId);
        }
    }
}
