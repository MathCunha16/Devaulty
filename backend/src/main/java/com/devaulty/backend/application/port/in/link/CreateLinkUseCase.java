package com.devaulty.backend.application.port.in.link;

import com.devaulty.backend.domain.model.Link;

public interface CreateLinkUseCase {
    Link execute(CreateLinkCommand command);
}
