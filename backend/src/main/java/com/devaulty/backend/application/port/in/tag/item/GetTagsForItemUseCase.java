package com.devaulty.backend.application.port.in.tag.item;

import com.devaulty.backend.domain.model.Tag;

import java.util.List;
import java.util.UUID;

public interface GetTagsForItemUseCase {
    List<Tag> execute(String itemType, UUID projectId, UUID itemId);
}
