package com.devaulty.backend.application.port.in.tag.item;

import java.util.UUID;

public interface RemoveTagFromItemUseCase {
    void execute(UUID projectId, String itemType ,UUID itemId, UUID tagId);
}
