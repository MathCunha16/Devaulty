package com.devaulty.backend.application.port.in.problem;

import com.devaulty.backend.domain.model.enums.ProblemSeverity;

import java.util.UUID;

public record UpdateProblemCommand(
        UUID id,
        UUID projectId,
        String title,
        String errorDescription,
        String solution,
        // Problem Status has a dedicated update logic
        ProblemSeverity severity
) {
}
