package com.devaulty.backend.adapter.in.web.note;

import com.devaulty.backend.adapter.in.web.note.dto.CreateNoteRequest;
import com.devaulty.backend.adapter.in.web.note.dto.NoteSummaryResponse;
import com.devaulty.backend.adapter.in.web.note.dto.NoteViewResponse;
import com.devaulty.backend.adapter.in.web.note.dto.UpdateNoteRequest;
import com.devaulty.backend.adapter.in.web.tag.TagWebMapper;
import com.devaulty.backend.application.port.in.note.CreateNoteCommand;
import com.devaulty.backend.application.port.in.note.UpdateNoteCommand;
import com.devaulty.backend.domain.model.Note;
import com.devaulty.backend.domain.model.Tag;
import org.mapstruct.Mapper;
import org.mapstruct.Mapping;

import java.util.List;
import java.util.UUID;

@Mapper(componentModel = "spring", uses = TagWebMapper.class)
public interface NoteWebMapper {

    @Mapping(target = "tags", source = "tags")
    NoteViewResponse toViewResponse(Note note, List<Tag> tags);

    @Mapping(target = "tags", source = "tags")
    NoteSummaryResponse toSummaryResponse(Note note, List<Tag> tags);

    @Mapping( target = "projectId", source = "projectId")
    CreateNoteCommand toCreateNoteCommand(CreateNoteRequest request, UUID projectId);

    @Mapping( target = "projectId", source = "projectId")
    @Mapping( target = "id", source = "noteId")
    UpdateNoteCommand toUpdateNoteCommand(UpdateNoteRequest request, UUID projectId, UUID noteId);
}
