package com.devaulty.backend.application.port.in.security;

public interface SetupMasterPasswordUseCase {
    void execute(char[] password);
}
