package com.devaulty.backend.application.port.in.tag;

import com.devaulty.backend.domain.model.Tag;

import java.util.UUID;

public interface GetTagByIdUseCase {
    Tag execute(UUID projectId ,UUID id);
}
