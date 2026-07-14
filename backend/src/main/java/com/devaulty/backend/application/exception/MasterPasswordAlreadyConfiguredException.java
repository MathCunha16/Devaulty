package com.devaulty.backend.application.exception;

public class MasterPasswordAlreadyConfiguredException extends DevaultyException{
    public MasterPasswordAlreadyConfiguredException() {
        super("Master password already configured");
    }
}
