package com.devaulty.backend.application.port.in.problem;

import com.devaulty.backend.domain.model.enums.ProblemSeverity;
import com.devaulty.backend.domain.model.enums.ProblemStatus;

import java.util.UUID;

public record CreateProblemCommand(
        UUID projectId,
        String title,
        String errorDescription,
        String solution,
        ProblemStatus status,
        ProblemSeverity severity
) {
}
