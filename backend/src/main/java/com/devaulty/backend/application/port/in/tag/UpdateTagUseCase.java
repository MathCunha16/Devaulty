package com.devaulty.backend.application.port.in.tag;

import com.devaulty.backend.domain.model.Tag;

public interface UpdateTagUseCase {
    Tag execute(UpdateTagCommand command);
}
