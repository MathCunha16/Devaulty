package com.devaulty.backend.adapter.out.persistence.project;

import com.devaulty.backend.domain.model.Project;
import org.mapstruct.Mapper;

@Mapper(componentModel = "spring")
public interface ProjectMapper {
    ProjectEntity toEntity(Project project);

    Project toDomain(ProjectEntity projectEntity);
}
