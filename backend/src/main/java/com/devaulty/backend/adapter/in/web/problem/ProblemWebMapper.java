package com.devaulty.backend.adapter.in.web.problem;

import com.devaulty.backend.adapter.in.web.problem.dto.*;
import com.devaulty.backend.application.port.in.problem.CreateProblemCommand;
import com.devaulty.backend.application.port.in.problem.UpdateProblemCommand;
import com.devaulty.backend.application.port.in.problem.UpdateProblemStatusCommand;
import com.devaulty.backend.domain.model.Problem;
import org.mapstruct.Mapper;
import org.mapstruct.Mapping;

import java.util.UUID;

@Mapper(componentModel = "spring")
public interface ProblemWebMapper {

    ProblemViewResponse toViewResponse(Problem problem);

    ProblemSummaryResponse toSummaryResponse(Problem problem);

    @Mapping(target = "projectId", source = "projectId")
    CreateProblemCommand toCreateProblemCommand(CreateProblemRequest request, UUID projectId);

    @Mapping(target = "projectId", source = "projectId")
    @Mapping(target = "id", source = "problemId")
    UpdateProblemCommand toUpdateProblemCommand(UpdateProblemRequest request, UUID projectId, UUID problemId);

    @Mapping(target = "projectId", source = "projectId")
    @Mapping(target = "id", source = "problemId")
    UpdateProblemStatusCommand toUpdateProblemStatusCommand(UpdateProblemStatusRequest request, UUID projectId, UUID problemId);
}
