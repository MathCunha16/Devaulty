package com.devaulty.backend.application.port.in.link;

import com.devaulty.backend.domain.model.Link;

import java.util.UUID;

public interface GetLinkByIdUseCase {
    Link execute(UUID projectId ,UUID id);
}
