package com.devaulty.backend.adapter.out.persistence.tag;

import com.devaulty.backend.domain.model.Tag;
import org.mapstruct.Mapper;
import org.mapstruct.Mapping;

@Mapper(componentModel = "spring")
public interface TagMapper {

    @Mapping(target = "project", ignore = true)
    TagEntity toEntity(Tag tag);

    @Mapping(target = "projectId", source = "project.id")
    Tag toDomain(TagEntity tagEntity);
}
