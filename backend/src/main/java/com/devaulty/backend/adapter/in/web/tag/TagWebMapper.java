package com.devaulty.backend.adapter.in.web.tag;

import com.devaulty.backend.adapter.in.web.tag.dto.CreateTagRequest;
import com.devaulty.backend.adapter.in.web.tag.dto.TagViewResponse;
import com.devaulty.backend.adapter.in.web.tag.dto.UpdateTagRequest;
import com.devaulty.backend.application.port.in.tag.CreateTagCommand;
import com.devaulty.backend.application.port.in.tag.UpdateTagCommand;
import com.devaulty.backend.domain.model.Tag;
import org.mapstruct.Mapper;
import org.mapstruct.Mapping;

import java.util.UUID;

@Mapper(componentModel = "spring")
public interface TagWebMapper {

    TagViewResponse toViewResponse(Tag tag);

    @Mapping(target = "projectId", source = "projectId")
    CreateTagCommand toCreateTagCommand(CreateTagRequest request, UUID projectId);

    @Mapping(target = "projectId", source = "projectId")
    @Mapping(target = "id", source = "tagId")
    UpdateTagCommand toUpdateTagCommand(UpdateTagRequest request, UUID projectId, UUID tagId);
}
