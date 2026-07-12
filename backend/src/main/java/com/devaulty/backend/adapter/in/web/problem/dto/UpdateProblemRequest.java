package com.devaulty.backend.adapter.in.web.problem.dto;

import com.devaulty.backend.domain.model.enums.ProblemSeverity;
import jakarta.validation.constraints.Pattern;
import jakarta.validation.constraints.Size;

public record UpdateProblemRequest(
        @Pattern(regexp = "(?s).*\\S.*", message = "Title must not be blank")
        @Size(min = 2, max = 255, message = "Title must be between 2 and 255 characters")
        String title,
        String errorDescription,
        String solution,
        ProblemSeverity severity
) {
}
