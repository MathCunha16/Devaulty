package com.devaulty.backend.adapter.out.persistence.link;

import com.devaulty.backend.domain.model.Link;
import org.mapstruct.Mapper;
import org.mapstruct.Mapping;

@Mapper(componentModel = "spring")
public interface LinkMapper {

    @Mapping(target = "project", ignore = true)
    LinkEntity toEntity(Link link);

    @Mapping(target = "projectId", source = "project.id")
    Link toDomain(LinkEntity linkEntity);
}
