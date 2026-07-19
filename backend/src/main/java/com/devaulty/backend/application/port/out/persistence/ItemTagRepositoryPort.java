package com.devaulty.backend.application.port.out.persistence;

import com.devaulty.backend.domain.model.Tag;

import java.util.List;
import java.util.Map;
import java.util.UUID;

public interface ItemTagRepositoryPort {

    void associateTagToItem(UUID tagId, String itemType, UUID itemId);

    void disassembleTagFromItem(UUID projectId, UUID tagId, String itemType, UUID itemId);

    void removeAllTagsFromItem(String itemType, UUID itemId);

    List<Tag> findTagsForItem(String itemType, UUID projectId, UUID itemId);

    Map<UUID, List<Tag>> findTagsForItems(String itemType, UUID projectId, List<UUID> itemIds);
}
