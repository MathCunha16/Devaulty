package com.devaulty.backend.adapter.out.persistence.snippet;

import com.devaulty.backend.domain.model.Snippet;
import org.mapstruct.Mapper;
import org.mapstruct.Mapping;

@Mapper(componentModel = "spring")
public interface SnippetMapper {

    @Mapping(target = "project", ignore = true)
    SnippetEntity toEntity(Snippet snippet);

    @Mapping(target = "projectId", source = "project.id")
    Snippet toDomain(SnippetEntity snippetEntity);
}
