package com.devaulty.backend.adapter.in.web.problem;

import com.devaulty.backend.adapter.in.web.problem.dto.*;
import com.devaulty.backend.adapter.in.web.tag.TagWebMapper;
import com.devaulty.backend.application.port.in.problem.CreateProblemCommand;
import com.devaulty.backend.application.port.in.problem.UpdateProblemCommand;
import com.devaulty.backend.application.port.in.problem.UpdateProblemStatusCommand;
import com.devaulty.backend.domain.model.Problem;
import com.devaulty.backend.domain.model.Tag;
import org.mapstruct.Mapper;
import org.mapstruct.Mapping;

import java.util.List;
import java.util.UUID;

@Mapper(componentModel = "spring", uses = TagWebMapper.class)
public interface ProblemWebMapper {

    @Mapping(target = "tags", source = "tags")
    ProblemViewResponse toViewResponse(Problem problem, List<Tag> tags);

    @Mapping(target = "tags", source = "tags")
    ProblemSummaryResponse toSummaryResponse(Problem problem, List<Tag> tags);

    @Mapping(target = "projectId", source = "projectId")
    CreateProblemCommand toCreateProblemCommand(CreateProblemRequest request, UUID projectId);

    @Mapping(target = "projectId", source = "projectId")
    @Mapping(target = "id", source = "problemId")
    UpdateProblemCommand toUpdateProblemCommand(UpdateProblemRequest request, UUID projectId, UUID problemId);

    @Mapping(target = "projectId", source = "projectId")
    @Mapping(target = "id", source = "problemId")
    UpdateProblemStatusCommand toUpdateProblemStatusCommand(UpdateProblemStatusRequest request, UUID projectId, UUID problemId);
}
