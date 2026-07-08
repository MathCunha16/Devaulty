package com.devaulty.backend.adapter.in.web.exception;

import com.devaulty.backend.application.exception.ResourceNotFoundException;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;
import org.springframework.http.HttpStatus;
import org.springframework.http.ResponseEntity;
import org.springframework.web.bind.annotation.ExceptionHandler;
import org.springframework.web.bind.annotation.RestControllerAdvice;

import java.time.LocalDateTime;

@RestControllerAdvice
public class GlobalExceptionHandler {

    private static final Logger logger = LoggerFactory.getLogger(GlobalExceptionHandler.class);

    @ExceptionHandler(ResourceNotFoundException.class)
    public ResponseEntity<ApiErrorResponse> handleResourceNotFoundException(ResourceNotFoundException exception){
        return buildResponse(HttpStatus.NOT_FOUND, exception.getMessage());
    }

    @ExceptionHandler(Exception.class)
    public ResponseEntity<ApiErrorResponse> handleGeral(Exception exception){
        logger.error("Unexpected error: ", exception);
        return buildResponse(HttpStatus.INTERNAL_SERVER_ERROR, "Internal server error, Please contact support.");
    }

    private ResponseEntity<ApiErrorResponse> buildResponse(HttpStatus status, String message){
        ApiErrorResponse response = new ApiErrorResponse(
                status.value(),
                message,
                LocalDateTime.now(),
                null // null because generic exceptions don't have a validation list field
        );

        return ResponseEntity.status(status).body(response);
    }
}
