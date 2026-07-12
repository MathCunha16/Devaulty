package com.devaulty.backend.adapter.in.web.problem.dto;

import com.devaulty.backend.domain.model.enums.ProblemStatus;
import jakarta.validation.constraints.NotNull;

public record UpdateProblemStatusRequest(
        @NotNull(message = "Status must not be null")
        ProblemStatus status
) {
}
