package com.devaulty.backend.application.port.in.credential;

import com.devaulty.backend.domain.model.Credential;

public interface CreateCredentialUseCase {
    Credential execute(CreateCredentialCommand command);
}
