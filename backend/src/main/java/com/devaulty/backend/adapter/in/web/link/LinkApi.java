package com.devaulty.backend.adapter.in.web.link;

import com.devaulty.backend.adapter.in.web.exception.ApiErrorResponse;
import com.devaulty.backend.adapter.in.web.link.dto.CreateLinkRequest;
import com.devaulty.backend.adapter.in.web.link.dto.LinkViewResponse;
import com.devaulty.backend.adapter.in.web.link.dto.UpdateLinkRequest;
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

@Tag(name = "Links", description = "Endpoints for managing links within a project")
public interface LinkApi {

    @Operation(
        summary = "Create a new link",
        description = "Creates a new link associated with the specified project ID."
    )
    @ApiResponses(value = {
        @ApiResponse(
            responseCode = "201",
            description = "Link created successfully",
            content = @Content(
                mediaType = "application/json",
                schema = @Schema(implementation = LinkViewResponse.class)
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
    ResponseEntity<LinkViewResponse> create(
            @PathVariable UUID projectId,
            @RequestBody @Valid CreateLinkRequest request
    );

    @Operation(
        summary = "Get all links by project",
        description = "Retrieves a paginated list of all links belonging to the specified project ID."
    )
    @ApiResponses(value = {
        @ApiResponse(
            responseCode = "200",
            description = "Successfully retrieved list of links"
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
    ResponseEntity<Page<LinkViewResponse>> getAllByProject(
            @PathVariable UUID projectId,
            @RequestParam(defaultValue = "0") int page,
            @RequestParam(defaultValue = "10") int size
    );

    @Operation(
        summary = "Get link by ID",
        description = "Retrieves a single link by its unique identifier and its project ID."
    )
    @ApiResponses(value = {
        @ApiResponse(
            responseCode = "200",
            description = "Successfully retrieved the link",
            content = @Content(
                mediaType = "application/json",
                schema = @Schema(implementation = LinkViewResponse.class)
            )
        ),
        @ApiResponse(
            responseCode = "404",
            description = "Project or Link not found",
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
    @GetMapping("/{linkId}")
    ResponseEntity<LinkViewResponse> getById(
            @PathVariable UUID projectId,
            @PathVariable UUID linkId
    );

    @Operation(
        summary = "Update a link",
        description = "Updates an existing link by its unique identifier and its project ID."
    )
    @ApiResponses(value = {
        @ApiResponse(
            responseCode = "200",
            description = "Link updated successfully",
            content = @Content(
                mediaType = "application/json",
                schema = @Schema(implementation = LinkViewResponse.class)
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
            description = "Project or Link not found",
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
    @PatchMapping("/{linkId}")
    ResponseEntity<LinkViewResponse> update(
            @PathVariable UUID projectId,
            @PathVariable UUID linkId,
            @RequestBody @Valid UpdateLinkRequest request
    );

    @Operation(
        summary = "Delete a link",
        description = "Deletes an existing link by its unique identifier and its project ID."
    )
    @ApiResponses(value = {
        @ApiResponse(
            responseCode = "204",
            description = "Link deleted successfully"
        ),
        @ApiResponse(
            responseCode = "404",
            description = "Project or Link not found",
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
    @DeleteMapping("/{linkId}")
    ResponseEntity<Void> delete(
            @PathVariable UUID projectId,
            @PathVariable UUID linkId
    );
}
