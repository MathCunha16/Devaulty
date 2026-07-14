package com.devaulty.backend.application.exception;

public class MasterPasswordNotConfiguredException extends DevaultyException{
    public MasterPasswordNotConfiguredException() {
        super("Master password not configured");
    }
}
