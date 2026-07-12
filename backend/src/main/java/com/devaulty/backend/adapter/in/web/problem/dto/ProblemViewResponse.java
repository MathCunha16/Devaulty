package com.devaulty.backend.adapter.in.web.problem.dto;

import com.devaulty.backend.domain.model.enums.ProblemSeverity;
import com.devaulty.backend.domain.model.enums.ProblemStatus;

import java.time.LocalDateTime;
import java.util.UUID;

public record ProblemViewResponse(
        UUID id,
        UUID projectId,
        String title,
        String errorDescription,
        String solution,
        ProblemStatus status,
        ProblemSeverity severity,
        LocalDateTime createdAt,
        LocalDateTime updatedAt
) {
}
