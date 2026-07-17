package com.devaulty.backend.adapter.in.web.exception;

import com.devaulty.backend.application.exception.*;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;
import org.springframework.http.HttpStatus;
import org.springframework.http.ResponseEntity;
import org.springframework.web.bind.MethodArgumentNotValidException;
import org.springframework.web.bind.annotation.ExceptionHandler;
import org.springframework.web.bind.annotation.RestControllerAdvice;

import java.time.LocalDateTime;
import java.util.List;

import com.devaulty.backend.adapter.in.web.security.dto.MasterPasswordRequest;
import com.devaulty.backend.adapter.in.web.credential.dto.CreateCredentialRequest;
import com.devaulty.backend.adapter.in.web.credential.dto.UpdateCredentialRequest;

@RestControllerAdvice
public class GlobalExceptionHandler {

    private static final Logger logger = LoggerFactory.getLogger(GlobalExceptionHandler.class);

    @ExceptionHandler(ResourceNotFoundException.class)
    public ResponseEntity<ApiErrorResponse> handleResourceNotFoundException(ResourceNotFoundException exception) {
        return buildResponse(HttpStatus.NOT_FOUND, exception.getMessage());
    }

    @ExceptionHandler(BusinessRuleException.class)
    public ResponseEntity<ApiErrorResponse> handleBusinessRuleException(BusinessRuleException exception) {
        return buildResponse(HttpStatus.BAD_REQUEST, exception.getMessage());
    }

    @ExceptionHandler(CryptoException.class)
    public ResponseEntity<ApiErrorResponse> handleCryptoException(CryptoException exception) {
        return buildResponse(HttpStatus.INTERNAL_SERVER_ERROR, exception.getMessage());
    }

    @ExceptionHandler(InvalidMasterPasswordException.class)
    public ResponseEntity<ApiErrorResponse> handleInvalidMasterPasswordException(InvalidMasterPasswordException exception) {
        return buildResponse(HttpStatus.UNAUTHORIZED, exception.getMessage());
    }

    @ExceptionHandler(JsonProcessingException.class)
    public ResponseEntity<ApiErrorResponse> handleJsonProcessingException(JsonProcessingException exception) {
        return buildResponse(HttpStatus.INTERNAL_SERVER_ERROR, exception.getMessage());
    }

    @ExceptionHandler(MasterPasswordAlreadyConfiguredException.class)
    public ResponseEntity<ApiErrorResponse> handleMasterPasswordAlreadyConfiguredException(MasterPasswordAlreadyConfiguredException exception) {
        return buildResponse(HttpStatus.CONFLICT, exception.getMessage());
    }

    @ExceptionHandler(MasterPasswordNotConfiguredException.class)
    public ResponseEntity<ApiErrorResponse> handleMasterPasswordNotConfiguredException(MasterPasswordNotConfiguredException exception) {
        return buildResponse(HttpStatus.FORBIDDEN, exception.getMessage());
    }

    @ExceptionHandler(VaultLockedException.class)
    public ResponseEntity<ApiErrorResponse> handleVaultLockedException(VaultLockedException exception) {
        return buildResponse(HttpStatus.LOCKED, exception.getMessage());
    }

    @ExceptionHandler(MethodArgumentNotValidException.class)
    public ResponseEntity<ApiErrorResponse> handleValidation(MethodArgumentNotValidException exception) {
        wipeRejectedRequestSecrets(exception.getBindingResult().getTarget());

        List<ApiErrorResponse.ValidationErrors> errors = exception.getBindingResult().getFieldErrors().stream()
                .map(err -> new ApiErrorResponse.ValidationErrors(err.getField(), err.getDefaultMessage()))
                .toList();

        ApiErrorResponse response = new ApiErrorResponse(
                HttpStatus.BAD_REQUEST.value(),
                "Validation error",
                LocalDateTime.now(),
                errors
        );

        return ResponseEntity.status(HttpStatus.BAD_REQUEST).body(response);
    }

    @ExceptionHandler(Exception.class)
    public ResponseEntity<ApiErrorResponse> handleGeral(Exception exception) {
        logger.error("Unexpected error: ", exception);
        return buildResponse(HttpStatus.INTERNAL_SERVER_ERROR, "Internal server error, Please contact support.");
    }

    private ResponseEntity<ApiErrorResponse> buildResponse(HttpStatus status, String message) {
        ApiErrorResponse response = new ApiErrorResponse(
                status.value(),
                message,
                LocalDateTime.now(),
                null // null because generic exceptions don't have a validation list field
        );

        return ResponseEntity.status(status).body(response);
    }

    private void wipeRejectedRequestSecrets(Object target) {
        if (target instanceof MasterPasswordRequest(char[] masterPassword)) {
            wipe(masterPassword);
        } else if (target instanceof CreateCredentialRequest(
                var title, var type, char[] username, char[] password, char[] apiKey, char[] rawTextContent, var notes,
                var url
        )) {
            wipe(username, password, apiKey, rawTextContent);
        } else if (target instanceof UpdateCredentialRequest(
                var title, var type, char[] username, char[] password, char[] apiKey, char[] rawTextContent, var notes,
                var url
        )) {
            wipe(username, password, apiKey, rawTextContent);
        }
    }

    private void wipe(char[]... arrays) {
        for (char[] array : arrays) {
            if (array != null) {
                java.util.Arrays.fill(array, '\0');
            }
        }
    }
}
