package com.devaulty.backend.application.port.in.problem;

import com.devaulty.backend.domain.model.enums.ProblemStatus;

import java.util.UUID;

public record UpdateProblemStatusCommand(
        UUID id,
        UUID projectId,
        ProblemStatus status
) {
}
