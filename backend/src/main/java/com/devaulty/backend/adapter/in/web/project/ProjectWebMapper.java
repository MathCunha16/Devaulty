package com.devaulty.backend.adapter.in.web.project;

import com.devaulty.backend.adapter.in.web.project.dto.CreateProjectRequest;
import com.devaulty.backend.adapter.in.web.project.dto.ProjectViewResponse;
import com.devaulty.backend.adapter.in.web.project.dto.UpdateProjectRequest;
import com.devaulty.backend.application.port.in.project.CreateProjectCommand;
import com.devaulty.backend.domain.model.Project;
import org.mapstruct.Mapper;

@Mapper(componentModel = "spring")
public interface ProjectWebMapper {

    ProjectViewResponse toViewResponse(Project project);

    CreateProjectCommand toCreateProjectCommand(CreateProjectRequest request);

    CreateProjectCommand toUpdateProjectCommand(UpdateProjectRequest request);
}
