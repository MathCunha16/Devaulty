package com.devaulty.backend.application.exception;

public class VaultLockedException extends DevaultyException{
    public VaultLockedException() {
        super("Vault is locked, please unlock it first");
    }
}
