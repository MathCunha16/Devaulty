package com.devaulty.backend.adapter.in.web.problem.dto;

import io.swagger.v3.oas.annotations.media.Schema;
import com.devaulty.backend.adapter.in.web.tag.dto.TagSummaryResponse;
import com.devaulty.backend.domain.model.enums.ProblemSeverity;
import com.devaulty.backend.domain.model.enums.ProblemStatus;

import java.time.LocalDateTime;
import java.util.List;
import java.util.UUID;

public record ProblemViewResponse(
        UUID id,
        UUID projectId,
        String title,
        String errorDescription,
        String solution,
        ProblemStatus status,
        ProblemSeverity severity,
        @Schema(description = "List of tags associated with this problem")
        List<TagSummaryResponse> tags,
        LocalDateTime createdAt,
        LocalDateTime updatedAt
) {
}
