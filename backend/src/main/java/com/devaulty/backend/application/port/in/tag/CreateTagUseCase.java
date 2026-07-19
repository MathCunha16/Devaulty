package com.devaulty.backend.application.port.in.tag;

import com.devaulty.backend.domain.model.Tag;

public interface CreateTagUseCase {
    Tag execute(CreateTagCommand command);
}
