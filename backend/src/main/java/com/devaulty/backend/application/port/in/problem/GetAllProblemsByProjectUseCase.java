package com.devaulty.backend.application.port.in.problem;

import com.devaulty.backend.domain.model.Problem;
import org.springframework.data.domain.Page;

import java.util.UUID;

public interface GetAllProblemsByProjectUseCase {
    Page<Problem> execute(UUID projectId, int page, int size);
}
