package com.devaulty.backend.application.exception;

public class ResourceNotFoundException extends DevaultyException {
    public ResourceNotFoundException(String resourceName, Object identifier) {
        super(String.format("%s not found with identifier %s", resourceName, identifier));
    }
}
