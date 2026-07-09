package com.devaulty.backend.application.exception;

public abstract class DevaultyException extends RuntimeException{
    protected DevaultyException(String message) {
        super(message);
    }

    protected DevaultyException(String message, Throwable cause) {
        super(message, cause);
    }
}
