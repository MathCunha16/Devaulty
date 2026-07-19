package com.devaulty.backend.adapter.in.web.tag;

import com.devaulty.backend.adapter.in.web.exception.ApiErrorResponse;
import com.devaulty.backend.adapter.in.web.tag.dto.CreateTagRequest;
import com.devaulty.backend.adapter.in.web.tag.dto.TagViewResponse;
import com.devaulty.backend.adapter.in.web.tag.dto.UpdateTagRequest;
import io.swagger.v3.oas.annotations.Operation;
import io.swagger.v3.oas.annotations.media.Content;
import io.swagger.v3.oas.annotations.media.Schema;
import io.swagger.v3.oas.annotations.responses.ApiResponse;
import io.swagger.v3.oas.annotations.responses.ApiResponses;
import io.swagger.v3.oas.annotations.tags.Tag;
import jakarta.validation.Valid;
import org.springframework.http.ResponseEntity;
import org.springframework.web.bind.annotation.*;

import java.util.List;
import java.util.UUID;

@Tag(name = "Tags", description = "Endpoints for managing tags within a project")
public interface TagApi {

    @Operation(
        summary = "Create a new tag",
        description = "Creates a new tag associated with the specified project ID."
    )
    @ApiResponses(value = {
        @ApiResponse(
            responseCode = "201",
            description = "Tag created successfully",
            content = @Content(
                mediaType = "application/json",
                schema = @Schema(implementation = TagViewResponse.class)
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
    ResponseEntity<TagViewResponse> create(
            @PathVariable UUID projectId,
            @RequestBody @Valid CreateTagRequest request
    );

    @Operation(
        summary = "Get all tags by project",
        description = "Retrieves a list of all tags belonging to the specified project ID."
    )
    @ApiResponses(value = {
        @ApiResponse(
            responseCode = "200",
            description = "Successfully retrieved tags list",
            content = @Content(
                mediaType = "application/json",
                schema = @Schema(implementation = TagViewResponse.class)
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
    @GetMapping
    ResponseEntity<List<TagViewResponse>> getAll(
            @PathVariable UUID projectId
    );

    @Operation(
        summary = "Get tag by ID",
        description = "Retrieves a single tag by its ID and project ID."
    )
    @ApiResponses(value = {
        @ApiResponse(
            responseCode = "200",
            description = "Tag retrieved successfully",
            content = @Content(
                mediaType = "application/json",
                schema = @Schema(implementation = TagViewResponse.class)
            )
        ),
        @ApiResponse(
            responseCode = "404",
            description = "Project or Tag not found",
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
    @GetMapping("/{tagId}")
    ResponseEntity<TagViewResponse> getById(
            @PathVariable UUID projectId,
            @PathVariable UUID tagId
    );

    @Operation(
        summary = "Search tag by name",
        description = "Searches for tags belonging to the project that contain the given name string."
    )
    @ApiResponses(value = {
        @ApiResponse(
            responseCode = "200",
            description = "Successfully performed search",
            content = @Content(
                mediaType = "application/json",
                schema = @Schema(implementation = TagViewResponse.class)
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
    @GetMapping("/search")
    ResponseEntity<List<TagViewResponse>> searchByName(
            @PathVariable UUID projectId,
            @RequestParam String name
    );

    @Operation(
        summary = "Update tag",
        description = "Updates the fields of an existing tag within a project."
    )
    @ApiResponses(value = {
        @ApiResponse(
            responseCode = "200",
            description = "Tag updated successfully",
            content = @Content(
                mediaType = "application/json",
                schema = @Schema(implementation = TagViewResponse.class)
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
            description = "Project or Tag not found",
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
    @PatchMapping("/{tagId}")
    ResponseEntity<TagViewResponse> update(
            @PathVariable UUID projectId,
            @PathVariable UUID tagId,
            @RequestBody @Valid UpdateTagRequest request
    );

    @Operation(
        summary = "Delete tag",
        description = "Deletes a tag by its ID and project ID."
    )
    @ApiResponses(value = {
        @ApiResponse(
            responseCode = "204",
            description = "Tag deleted successfully"
        ),
        @ApiResponse(
            responseCode = "404",
            description = "Project or Tag not found",
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
    @DeleteMapping("/{tagId}")
    ResponseEntity<Void> delete(
            @PathVariable UUID projectId,
            @PathVariable UUID tagId
    );
}
