package com.devaulty.backend.application.port.in.security;

public interface UnlockVaultUseCase {
    boolean execute(char[] password);
}
