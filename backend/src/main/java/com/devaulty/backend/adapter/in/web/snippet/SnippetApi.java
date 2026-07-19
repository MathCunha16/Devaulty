package com.devaulty.backend.adapter.in.web.snippet;

import com.devaulty.backend.adapter.in.web.exception.ApiErrorResponse;
import com.devaulty.backend.adapter.in.web.snippet.dto.CreateSnippetRequest;
import com.devaulty.backend.adapter.in.web.snippet.dto.SnippetSummaryResponse;
import com.devaulty.backend.adapter.in.web.snippet.dto.SnippetViewResponse;
import com.devaulty.backend.adapter.in.web.snippet.dto.UpdateSnippetRequest;
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

@Tag(name = "Snippets", description = "Endpoints for managing snippets within a project")
public interface SnippetApi {

    @Operation(
        summary = "Create a new snippet",
        description = "Creates a new snippet associated with the specified project ID."
    )
    @ApiResponses(value = {
        @ApiResponse(
            responseCode = "201",
            description = "Snippet created successfully",
            content = @Content(
                mediaType = "application/json",
                schema = @Schema(implementation = SnippetViewResponse.class)
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
    ResponseEntity<SnippetViewResponse> create(
            @PathVariable UUID projectId,
            @RequestBody @Valid CreateSnippetRequest request
    );

    @Operation(
        summary = "Get all snippets by project",
        description = "Retrieves a paginated list of all snippets belonging to the specified project ID."
    )
    @ApiResponses(value = {
        @ApiResponse(
            responseCode = "200",
            description = "Successfully retrieved list of snippets"
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
    ResponseEntity<Page<SnippetSummaryResponse>> getAllByProject(
            @PathVariable UUID projectId,
            @RequestParam(defaultValue = "0") int page,
            @RequestParam(defaultValue = "10") int size
    );

    @Operation(
        summary = "Get snippet by ID",
        description = "Retrieves a single snippet by its unique identifier and its project ID."
    )
    @ApiResponses(value = {
        @ApiResponse(
            responseCode = "200",
            description = "Successfully retrieved the snippet",
            content = @Content(
                mediaType = "application/json",
                schema = @Schema(implementation = SnippetViewResponse.class)
            )
        ),
        @ApiResponse(
            responseCode = "404",
            description = "Project or Snippet not found",
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
    @GetMapping("/{snippetId}")
    ResponseEntity<SnippetViewResponse> getById(
            @PathVariable UUID projectId,
            @PathVariable UUID snippetId
    );

    @Operation(
        summary = "Update a snippet",
        description = "Updates an existing snippet by its unique identifier and its project ID."
    )
    @ApiResponses(value = {
        @ApiResponse(
            responseCode = "200",
            description = "Snippet updated successfully",
            content = @Content(
                mediaType = "application/json",
                schema = @Schema(implementation = SnippetViewResponse.class)
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
            description = "Project or Snippet not found",
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
    @PatchMapping("/{snippetId}")
    ResponseEntity<SnippetViewResponse> update(
            @PathVariable UUID projectId,
            @PathVariable UUID snippetId,
            @RequestBody @Valid UpdateSnippetRequest request
    );

    @Operation(
        summary = "Delete a snippet",
        description = "Deletes an existing snippet by its unique identifier and its project ID."
    )
    @ApiResponses(value = {
        @ApiResponse(
            responseCode = "204",
            description = "Snippet deleted successfully"
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
    @DeleteMapping("/{snippetId}")
    ResponseEntity<Void> delete(
            @PathVariable UUID projectId,
            @PathVariable UUID snippetId
    );
}
