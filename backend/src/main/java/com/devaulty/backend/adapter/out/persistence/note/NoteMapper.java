package com.devaulty.backend.adapter.out.persistence.note;

import com.devaulty.backend.domain.model.Note;
import org.mapstruct.Mapper;
import org.mapstruct.Mapping;

@Mapper(componentModel = "spring")
public interface NoteMapper {

    @Mapping(target = "project", ignore = true)
    NoteEntity toEntity(Note note);

    @Mapping(target = "projectId", source = "project.id")
    Note toDomain(NoteEntity noteEntity);

}
