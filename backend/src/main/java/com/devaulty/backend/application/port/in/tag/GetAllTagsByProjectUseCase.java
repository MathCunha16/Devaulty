package com.devaulty.backend.application.port.in.tag;

import com.devaulty.backend.domain.model.Tag;

import java.util.List;
import java.util.UUID;

public interface GetAllTagsByProjectUseCase {
    List<Tag> execute(UUID projectId);
}
