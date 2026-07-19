package com.devaulty.backend.application.exception;

public class ResourceAlreadyExistsException extends DevaultyException{
    public ResourceAlreadyExistsException(String resourceName, Object identifier) {
        super(String.format("%s already exists with identifier %s", resourceName, identifier));
    }
}
