package com.devaulty.backend.adapter.in.web.snippet;

import com.devaulty.backend.adapter.in.web.snippet.dto.CreateSnippetRequest;
import com.devaulty.backend.adapter.in.web.snippet.dto.SnippetViewResponse;
import com.devaulty.backend.adapter.in.web.snippet.dto.UpdateSnippetRequest;
import com.devaulty.backend.application.port.in.snippet.CreateSnippetCommand;
import com.devaulty.backend.application.port.in.snippet.UpdateSnippetCommand;
import com.devaulty.backend.domain.model.Snippet;
import org.mapstruct.Mapper;
import org.mapstruct.Mapping;

import java.util.UUID;

@Mapper(componentModel = "spring")
public interface SnippetWebMapper {

    SnippetViewResponse toViewResponse(Snippet snippet);

    @Mapping(target = "projectId", source = "projectId")
    CreateSnippetCommand toCreateSnippetCommand(CreateSnippetRequest request, UUID projectId);

    @Mapping(target = "projectId", source = "projectId")
    @Mapping(target = "id", source = "snippetId")
    UpdateSnippetCommand toUpdateSnippetCommand(UpdateSnippetRequest request, UUID projectId, UUID snippetId);
}
