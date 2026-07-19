package com.devaulty.backend.adapter.in.web.note;

import com.devaulty.backend.adapter.in.web.note.dto.CreateNoteRequest;
import com.devaulty.backend.adapter.in.web.note.dto.NoteSummaryResponse;
import com.devaulty.backend.adapter.in.web.note.dto.NoteViewResponse;
import com.devaulty.backend.adapter.in.web.note.dto.UpdateNoteRequest;
import com.devaulty.backend.adapter.in.web.util.UriLocationBuilderHelper;
import com.devaulty.backend.application.port.in.note.*;
import com.devaulty.backend.domain.model.Note;
import com.devaulty.backend.application.port.in.tag.item.GetTagsForItemUseCase;
import com.devaulty.backend.application.port.in.tag.item.GetTagsForItemsUseCase;
import com.devaulty.backend.domain.model.Tag;
import jakarta.validation.Valid;
import org.springframework.data.domain.Page;
import org.springframework.http.ResponseEntity;
import org.springframework.web.bind.annotation.*;

import java.net.URI;
import java.util.List;
import java.util.Map;
import java.util.UUID;

@RestController
@RequestMapping("/ap1/v1/projects/{projectId}/notes")
public class NoteController implements NoteApi{

    private final CreateNoteUseCase createNoteUseCase;
    private final GetAllNotesByProjectUseCase getAllNotesByProjectUseCase;
    private final GetNoteByIdUseCase getNoteByIdUseCase;
    private final ArchiveNoteUseCase archiveNoteUseCase;
    private final UpdateNoteUseCase updateNoteUseCase;
    private final UnarchiveNoteUseCase unarchiveNoteUseCase;
    private final DeleteNoteUseCase deleteNoteUseCase;
    private final GetTagsForItemUseCase getTagsForItemUseCase;
    private final GetTagsForItemsUseCase getTagsForItemsUseCase;
    private final NoteWebMapper mapper;
    private final UriLocationBuilderHelper uriLocationBuilderHelper;

    private static final String ITEM_TYPE = "note";

    public NoteController(CreateNoteUseCase createNoteUseCase, GetAllNotesByProjectUseCase getAllNotesByProjectUseCase, GetNoteByIdUseCase getNoteByIdUseCase, ArchiveNoteUseCase archiveNoteUseCase, UpdateNoteUseCase updateNoteUseCase, UnarchiveNoteUseCase unarchiveNoteUseCase, DeleteNoteUseCase deleteNoteUseCase, GetTagsForItemUseCase getTagsForItemUseCase, GetTagsForItemsUseCase getTagsForItemsUseCase, NoteWebMapper mapper, UriLocationBuilderHelper uriLocationBuilderHelper) {
        this.createNoteUseCase = createNoteUseCase;
        this.getAllNotesByProjectUseCase = getAllNotesByProjectUseCase;
        this.getNoteByIdUseCase = getNoteByIdUseCase;
        this.archiveNoteUseCase = archiveNoteUseCase;
        this.updateNoteUseCase = updateNoteUseCase;
        this.unarchiveNoteUseCase = unarchiveNoteUseCase;
        this.deleteNoteUseCase = deleteNoteUseCase;
        this.getTagsForItemUseCase = getTagsForItemUseCase;
        this.getTagsForItemsUseCase = getTagsForItemsUseCase;
        this.mapper = mapper;
        this.uriLocationBuilderHelper = uriLocationBuilderHelper;
    }

    @Override
    @PostMapping
    public ResponseEntity<NoteViewResponse> create(
            @PathVariable UUID projectId,
            @RequestBody @Valid CreateNoteRequest request
            )
    {
        CreateNoteCommand command = mapper.toCreateNoteCommand(request, projectId);
        Note note = createNoteUseCase.execute(command);
        URI location = uriLocationBuilderHelper.buildLocationUri(note.getId());
        return ResponseEntity.created(location).body(mapper.toViewResponse(note, List.of()));
    }

    @Override
    @GetMapping
    public ResponseEntity<Page<NoteSummaryResponse>> getAllByProject(
            @PathVariable UUID projectId,
            @RequestParam(defaultValue = "0") int page,
            @RequestParam(defaultValue = "10") int size
    ){
        Page<Note> notes = getAllNotesByProjectUseCase.execute(projectId, page, size);
        List<UUID> ids = notes.getContent().stream().map(Note::getId).toList();
        Map<UUID, List<Tag>> tagsByItem = getTagsForItemsUseCase.execute(ITEM_TYPE, projectId, ids);
        Page<NoteSummaryResponse> response = notes.map(s ->
                mapper.toSummaryResponse(s, tagsByItem.getOrDefault(s.getId(), List.of())));
        return ResponseEntity.ok(response);
    }

    @Override
    @GetMapping("/{noteID}")
    public ResponseEntity<NoteViewResponse> getById(
            @PathVariable UUID projectId,
            @PathVariable UUID noteID
    )
    {
        Note note = getNoteByIdUseCase.execute(projectId, noteID);
        List<Tag> tags = getTagsForItemUseCase.execute(ITEM_TYPE, projectId, noteID);
        return ResponseEntity.ok(mapper.toViewResponse(note, tags));
    }

    @Override
    @PatchMapping("/{noteId}")
    public ResponseEntity<NoteViewResponse> update(
            @PathVariable UUID projectId,
            @PathVariable UUID noteId,
            @RequestBody @Valid UpdateNoteRequest request
    ){
        UpdateNoteCommand command = mapper.toUpdateNoteCommand(request, projectId, noteId);
        Note note = updateNoteUseCase.execute(command);
        List<Tag> tags = getTagsForItemUseCase.execute(ITEM_TYPE, projectId, noteId);
        return ResponseEntity.ok(mapper.toViewResponse(note, tags));
    }

    @Override
    @PatchMapping("/{noteId}/archive")
    public ResponseEntity<Void> archive(
            @PathVariable UUID projectId,
            @PathVariable UUID noteId
    ){
        archiveNoteUseCase.execute(projectId, noteId);
        return ResponseEntity.noContent().build();
    }

    @Override
    @PatchMapping("/{noteId}/unarchive")
    public ResponseEntity<Void> unarchive(
            @PathVariable UUID projectId,
            @PathVariable UUID noteId
    ){
        unarchiveNoteUseCase.execute(projectId, noteId);
        return ResponseEntity.noContent().build();
    }

    @Override
    @DeleteMapping("/{noteId}")
    public ResponseEntity<Void> delete(
            @PathVariable UUID projectId,
            @PathVariable UUID noteId
    ){
        deleteNoteUseCase.execute(projectId, noteId);
        return ResponseEntity.noContent().build();
    }
}
