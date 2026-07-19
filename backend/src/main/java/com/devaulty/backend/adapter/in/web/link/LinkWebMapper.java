package com.devaulty.backend.adapter.in.web.link;

import com.devaulty.backend.adapter.in.web.link.dto.CreateLinkRequest;
import com.devaulty.backend.adapter.in.web.link.dto.LinkViewResponse;
import com.devaulty.backend.adapter.in.web.link.dto.UpdateLinkRequest;
import com.devaulty.backend.adapter.in.web.tag.TagWebMapper;
import com.devaulty.backend.application.port.in.link.CreateLinkCommand;
import com.devaulty.backend.application.port.in.link.UpdateLinkCommand;
import com.devaulty.backend.domain.model.Link;
import com.devaulty.backend.domain.model.Tag;
import org.mapstruct.Mapper;
import org.mapstruct.Mapping;

import java.util.List;
import java.util.UUID;

@Mapper(componentModel = "spring", uses = TagWebMapper.class)
public interface LinkWebMapper {

    @Mapping(target = "tags", source = "tags")
    LinkViewResponse toViewResponse(Link link, List<Tag> tags);

    @Mapping(target = "projectId", source = "projectId")
    CreateLinkCommand toCreateLinkCommand(CreateLinkRequest request, UUID projectId);

    @Mapping(target = "projectId", source = "projectId")
    @Mapping(target = "id", source = "linkId")
    UpdateLinkCommand toUpdateLinkCommand(UpdateLinkRequest request, UUID projectId, UUID linkId);
}
