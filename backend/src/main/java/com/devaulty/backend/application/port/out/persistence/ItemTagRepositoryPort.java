package com.devaulty.backend.application.port.out.persistence;

import java.util.UUID;

public interface ItemTagRepositoryPort {

    void associateTagToItem(UUID tagId, String itemType, UUID itemId);

    void disassembleTagFromItem(UUID tagId, String itemType, UUID itemId);

    void removeAllTagsFromItem(String itemType, UUID itemId);
}
