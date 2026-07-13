package com.devaulty.backend.adapter.in.web.note;

import com.devaulty.backend.adapter.in.web.exception.ApiErrorResponse;
import com.devaulty.backend.adapter.in.web.note.dto.CreateNoteRequest;
import com.devaulty.backend.adapter.in.web.note.dto.NoteSummaryResponse;
import com.devaulty.backend.adapter.in.web.note.dto.NoteViewResponse;
import com.devaulty.backend.adapter.in.web.note.dto.UpdateNoteRequest;
import io.swagger.v3.oas.annotations.Operation;
import io.swagger.v3.oas.annotations.media.Content;
import io.swagger.v3.oas.annotations.media.Schema;
import io.swagger.v3.oas.annotations.responses.ApiResponse;
import io.swagger.v3.oas.annotations.responses.ApiResponses;
import io.swagger.v3.oas.annotations.tags.Tag;
import jakarta.validation.Valid;
import org.springframework.data.domain.Page;
import org.springframework.http.ResponseEntity;
import org.springframework.web.bind.annotation.*;

import java.util.UUID;

@Tag(name = "Notes", description = "Endpoints for managing project notes")
public interface NoteApi {

    @Operation(
        summary = "Create a new note",
        description = "Creates a new note associated with the specified project ID."
    )
    @ApiResponses(value = {
        @ApiResponse(
            responseCode = "201",
            description = "Note created successfully",
            content = @Content(
                mediaType = "application/json",
                schema = @Schema(implementation = NoteViewResponse.class)
            )
        ),
        @ApiResponse(
            responseCode = "400",
            description = "Invalid request body or validation constraints failed",
            content = @Content(
                mediaType = "application/json",
                schema = @Schema(implementation = ApiErrorResponse.class)
            )
        ),
        @ApiResponse(
            responseCode = "404",
            description = "Project not found",
            content = @Content(
                mediaType = "application/json",
                schema = @Schema(implementation = ApiErrorResponse.class)
            )
        ),
        @ApiResponse(
            responseCode = "500",
            description = "Internal server error occurred",
            content = @Content(
                mediaType = "application/json",
                schema = @Schema(implementation = ApiErrorResponse.class)
            )
        )
    })
    @PostMapping
    ResponseEntity<NoteViewResponse> create(
            @PathVariable UUID projectId,
            @RequestBody @Valid CreateNoteRequest request
    );

    @Operation(
        summary = "Get all notes by project",
        description = "Retrieves a paginated list of all note summaries belonging to the specified project ID."
    )
    @ApiResponses(value = {
        @ApiResponse(
            responseCode = "200",
            description = "Successfully retrieved list of note summaries"
        ),
        @ApiResponse(
            responseCode = "404",
            description = "Project not found",
            content = @Content(
                mediaType = "application/json",
                schema = @Schema(implementation = ApiErrorResponse.class)
            )
        ),
        @ApiResponse(
            responseCode = "500",
            description = "Internal server error occurred",
            content = @Content(
                mediaType = "application/json",
                schema = @Schema(implementation = ApiErrorResponse.class)
            )
        )
    })
    @GetMapping
    ResponseEntity<Page<NoteSummaryResponse>> getAllByProject(
            @PathVariable UUID projectId,
            @RequestParam(defaultValue = "0") int page,
            @RequestParam(defaultValue = "10") int size
    );

    @Operation(
        summary = "Get note by ID",
        description = "Retrieves a single note by its unique identifier and its project ID."
    )
    @ApiResponses(value = {
        @ApiResponse(
            responseCode = "200",
            description = "Successfully retrieved the note",
            content = @Content(
                mediaType = "application/json",
                schema = @Schema(implementation = NoteViewResponse.class)
            )
        ),
        @ApiResponse(
            responseCode = "404",
            description = "Project or Note not found",
            content = @Content(
                mediaType = "application/json",
                schema = @Schema(implementation = ApiErrorResponse.class)
            )
        ),
        @ApiResponse(
            responseCode = "500",
            description = "Internal server error occurred",
            content = @Content(
                mediaType = "application/json",
                schema = @Schema(implementation = ApiErrorResponse.class)
            )
        )
    })
    @GetMapping("/{noteID}")
    ResponseEntity<NoteViewResponse> getById(
            @PathVariable UUID projectId,
            @PathVariable UUID noteID
    );

    @Operation(
        summary = "Update a note",
        description = "Updates an existing note by its unique identifier and its project ID."
    )
    @ApiResponses(value = {
        @ApiResponse(
            responseCode = "200",
            description = "Note updated successfully",
            content = @Content(
                mediaType = "application/json",
                schema = @Schema(implementation = NoteViewResponse.class)
            )
        ),
        @ApiResponse(
            responseCode = "400",
            description = "Invalid request body or validation constraints failed",
            content = @Content(
                mediaType = "application/json",
                schema = @Schema(implementation = ApiErrorResponse.class)
            )
        ),
        @ApiResponse(
            responseCode = "404",
            description = "Project or Note not found",
            content = @Content(
                mediaType = "application/json",
                schema = @Schema(implementation = ApiErrorResponse.class)
            )
        ),
        @ApiResponse(
            responseCode = "500",
            description = "Internal server error occurred",
            content = @Content(
                mediaType = "application/json",
                schema = @Schema(implementation = ApiErrorResponse.class)
            )
        )
    })
    @PatchMapping("/{noteId}")
    ResponseEntity<NoteViewResponse> update(
            @PathVariable UUID projectId,
            @PathVariable UUID noteId,
            @RequestBody @Valid UpdateNoteRequest request
    );

    @Operation(
        summary = "Archive a note",
        description = "Archives an existing note by its unique identifier."
    )
    @ApiResponses(value = {
        @ApiResponse(
            responseCode = "204",
            description = "Note archived successfully"
        ),
        @ApiResponse(
            responseCode = "400",
            description = "Note is already archived",
            content = @Content(
                mediaType = "application/json",
                schema = @Schema(implementation = ApiErrorResponse.class)
            )
        ),
        @ApiResponse(
            responseCode = "404",
            description = "Project or Note not found",
            content = @Content(
                mediaType = "application/json",
                schema = @Schema(implementation = ApiErrorResponse.class)
            )
        ),
        @ApiResponse(
            responseCode = "500",
            description = "Internal server error occurred",
            content = @Content(
                mediaType = "application/json",
                schema = @Schema(implementation = ApiErrorResponse.class)
            )
        )
    })
    @PatchMapping("/{noteId}/archive")
    ResponseEntity<Void> archive(
            @PathVariable UUID projectId,
            @PathVariable UUID noteId
    );

    @Operation(
        summary = "Unarchive a note",
        description = "Unarchives an archived note by its unique identifier."
    )
    @ApiResponses(value = {
        @ApiResponse(
            responseCode = "204",
            description = "Note unarchived successfully"
        ),
        @ApiResponse(
            responseCode = "400",
            description = "Note is not archived",
            content = @Content(
                mediaType = "application/json",
                schema = @Schema(implementation = ApiErrorResponse.class)
            )
        ),
        @ApiResponse(
            responseCode = "404",
            description = "Project or Note not found",
            content = @Content(
                mediaType = "application/json",
                schema = @Schema(implementation = ApiErrorResponse.class)
            )
        ),
        @ApiResponse(
            responseCode = "500",
            description = "Internal server error occurred",
            content = @Content(
                mediaType = "application/json",
                schema = @Schema(implementation = ApiErrorResponse.class)
            )
        )
    })
    @PatchMapping("/{noteId}/unarchive")
    ResponseEntity<Void> unarchive(
            @PathVariable UUID projectId,
            @PathVariable UUID noteId
    );

    @Operation(
        summary = "Delete a note",
        description = "Deletes an existing note by its unique identifier."
    )
    @ApiResponses(value = {
        @ApiResponse(
            responseCode = "204",
            description = "Note deleted successfully"
        ),
        @ApiResponse(
            responseCode = "404",
            description = "Project or Note not found",
            content = @Content(
                mediaType = "application/json",
                schema = @Schema(implementation = ApiErrorResponse.class)
            )
        ),
        @ApiResponse(
            responseCode = "500",
            description = "Internal server error occurred",
            content = @Content(
                mediaType = "application/json",
                schema = @Schema(implementation = ApiErrorResponse.class)
            )
        )
    })
    @DeleteMapping("/{noteId}")
    ResponseEntity<Void> delete(
            @PathVariable UUID projectId,
            @PathVariable UUID noteId
    );
}
