package com.devaulty.backend.adapter.in.web.problem.dto;

import com.devaulty.backend.domain.model.enums.ProblemSeverity;
import com.devaulty.backend.domain.model.enums.ProblemStatus;
import jakarta.validation.constraints.NotBlank;
import jakarta.validation.constraints.NotNull;
import jakarta.validation.constraints.Size;

public record CreateProblemRequest(
        @NotBlank(message = "Title must not be blank")
        @Size(min = 2, max = 255, message = "Title must be between 2 and 255 characters")
        String title,
        String errorDescription,
        String solution,
        @NotNull(message = "Status must not be null")
        ProblemStatus status,
        @NotNull(message = "Severity must not be null")
        ProblemSeverity severity
) {
}
