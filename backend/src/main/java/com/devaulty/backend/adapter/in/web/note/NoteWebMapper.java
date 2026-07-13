package com.devaulty.backend.adapter.in.web.note;

import com.devaulty.backend.adapter.in.web.note.dto.CreateNoteRequest;
import com.devaulty.backend.adapter.in.web.note.dto.NoteSummaryResponse;
import com.devaulty.backend.adapter.in.web.note.dto.NoteViewResponse;
import com.devaulty.backend.adapter.in.web.note.dto.UpdateNoteRequest;
import com.devaulty.backend.application.port.in.note.CreateNoteCommand;
import com.devaulty.backend.application.port.in.note.UpdateNoteCommand;
import com.devaulty.backend.domain.model.Note;
import org.mapstruct.Mapper;
import org.mapstruct.Mapping;

import java.util.UUID;

@Mapper(componentModel = "spring")
public interface NoteWebMapper {

    NoteViewResponse toViewResponse(Note note);

    NoteSummaryResponse toSummaryResponse(Note note);

    @Mapping( target = "projectId", source = "projectId")
    CreateNoteCommand toCreateNoteCommand(CreateNoteRequest request, UUID projectId);

    @Mapping( target = "projectId", source = "projectId")
    @Mapping( target = "id", source = "noteId")
    UpdateNoteCommand toUpdateNoteCommand(UpdateNoteRequest request, UUID projectId, UUID noteId);
}
