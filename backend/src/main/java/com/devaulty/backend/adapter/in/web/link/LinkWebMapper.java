package com.devaulty.backend.adapter.in.web.link;

import com.devaulty.backend.adapter.in.web.link.dto.CreateLinkRequest;
import com.devaulty.backend.adapter.in.web.link.dto.LinkViewResponse;
import com.devaulty.backend.adapter.in.web.link.dto.UpdateLinkRequest;
import com.devaulty.backend.application.port.in.link.CreateLinkCommand;
import com.devaulty.backend.application.port.in.link.UpdateLinkCommand;
import com.devaulty.backend.domain.model.Link;
import org.mapstruct.Mapper;
import org.mapstruct.Mapping;

import java.util.UUID;

@Mapper(componentModel = "spring")
public interface LinkWebMapper {

    LinkViewResponse toViewResponse(Link link);

    @Mapping(target = "projectId", source = "projectId")
    CreateLinkCommand toCreateLinkCommand(CreateLinkRequest request, UUID projectId);

    @Mapping(target = "projectId", source = "projectId")
    @Mapping(target = "id", source = "linkId")
    UpdateLinkCommand toUpdateLinkCommand(UpdateLinkRequest request, UUID projectId, UUID linkId);
}
