package com.devaulty.backend.application.port.in.credential;

public interface CreateCredentialUseCase {
    DecryptedCredential execute(CreateCredentialCommand command);
}
