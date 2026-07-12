package com.devaulty.backend.adapter.in.web.problem.dto;

import com.devaulty.backend.domain.model.enums.ProblemSeverity;
import com.devaulty.backend.domain.model.enums.ProblemStatus;

import java.time.LocalDateTime;
import java.util.UUID;

public record ProblemSummaryResponse(
        UUID id,
        UUID projectId,
        String title,
        ProblemStatus status,
        ProblemSeverity severity,
        LocalDateTime createdAt,
        LocalDateTime updatedAt
) {
}
