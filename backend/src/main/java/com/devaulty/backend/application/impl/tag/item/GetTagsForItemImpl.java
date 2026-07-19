package com.devaulty.backend.application.impl.tag.item;

import com.devaulty.backend.application.exception.ResourceNotFoundException;
import com.devaulty.backend.application.port.in.tag.item.GetTagsForItemUseCase;
import com.devaulty.backend.application.port.out.persistence.ItemTagRepositoryPort;
import com.devaulty.backend.application.port.out.persistence.ProjectRepositoryPort;
import com.devaulty.backend.application.port.out.persistence.ProjectScopedRepositoryPort;
import com.devaulty.backend.domain.model.Tag;
import org.springframework.transaction.annotation.Transactional;

import java.util.List;
import java.util.UUID;

public class GetTagsForItemImpl implements GetTagsForItemUseCase {

    private final ItemTagRepositoryPort itemTagRepositoryPort;
    private final ProjectRepositoryPort projectRepositoryPort;
    private final List<ProjectScopedRepositoryPort> projectScopedRepositories;

    public GetTagsForItemImpl(ItemTagRepositoryPort itemTagRepositoryPort,
                              ProjectRepositoryPort projectRepositoryPort,
                              List<ProjectScopedRepositoryPort> projectScopedRepositories) {
        this.itemTagRepositoryPort = itemTagRepositoryPort;
        this.projectRepositoryPort = projectRepositoryPort;
        this.projectScopedRepositories = projectScopedRepositories;
    }

    @Override
    @Transactional(readOnly = true)
    public List<Tag> execute(String itemType, UUID projectId, UUID itemId) {

        if (!projectRepositoryPort.existsById(projectId)) {
            throw new ResourceNotFoundException("Project", projectId);
        }

        String canonicalType = validateItemOwnership(itemType, projectId, itemId);

        return itemTagRepositoryPort.findTagsForItem(canonicalType, projectId, itemId);
    }

    private String validateItemOwnership(String itemType, UUID projectId, UUID itemId) {
        ProjectScopedRepositoryPort repo = projectScopedRepositories.stream()
                .filter(r -> r.getSupportedType().equalsIgnoreCase(itemType))
                .findFirst()
                .orElseThrow(() -> new IllegalArgumentException("Unsupported item type: " + itemType));

        boolean exists = repo.existsByIdAndProjectId(itemId, projectId);

        if (!exists) {
            String resourceName = itemType.substring(0, 1).toUpperCase() + itemType.substring(1).toLowerCase();
            throw new ResourceNotFoundException(resourceName, itemId);
        }

        return repo.getSupportedType();
    }
}
