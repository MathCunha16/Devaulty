package com.devaulty.backend.adapter.in.web.exception;

import java.time.LocalDateTime;
import java.util.List;

public record ApiErrorResponse(
        int status,
        String message,
        LocalDateTime timestamp,
        List<ValidationErrors> errors
) {
    // Internal record to detail which field has an error
    public record ValidationErrors(String field, String message) {}
}
