package com.devaulty.backend.application.port.in.tag;

import com.devaulty.backend.domain.model.Tag;

import java.util.List;
import java.util.UUID;

public interface SearchTagByNameUseCase {
    List<Tag> execute(UUID projectId, String name);
}
