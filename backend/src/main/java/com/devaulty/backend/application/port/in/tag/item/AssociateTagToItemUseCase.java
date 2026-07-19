package com.devaulty.backend.application.port.in.tag.item;

import java.util.UUID;

public interface AssociateTagToItemUseCase {
    void execute(UUID projectId, String itemType ,UUID itemId, UUID tagId);
}
