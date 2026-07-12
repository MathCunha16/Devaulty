package com.devaulty.backend.adapter.out.persistence.problem;

import com.devaulty.backend.domain.model.Problem;
import org.mapstruct.Mapper;
import org.mapstruct.Mapping;

@Mapper(componentModel = "spring")
public interface ProblemMapper {

    @Mapping(target = "project", ignore = true)
    ProblemEntity toEntity(Problem problem);

    @Mapping(target = "projectId", source = "project.id")
    Problem toDomain(ProblemEntity problemEntity);
}
