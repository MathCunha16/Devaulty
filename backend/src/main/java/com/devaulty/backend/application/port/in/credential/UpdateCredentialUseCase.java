package com.devaulty.backend.application.port.in.credential;

public interface UpdateCredentialUseCase {
    DecryptedCredential execute(UpdateCredentialCommand command);
}
