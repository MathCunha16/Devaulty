package com.devaulty.backend.adapter.in.web.tag.item;

import com.devaulty.backend.adapter.in.web.exception.ApiErrorResponse;
import io.swagger.v3.oas.annotations.Operation;
import io.swagger.v3.oas.annotations.media.Content;
import io.swagger.v3.oas.annotations.media.Schema;
import io.swagger.v3.oas.annotations.responses.ApiResponse;
import io.swagger.v3.oas.annotations.responses.ApiResponses;
import io.swagger.v3.oas.annotations.tags.Tag;
import org.springframework.http.ResponseEntity;
import org.springframework.web.bind.annotation.*;

import java.util.UUID;

@Tag(name = "Item Tags", description = "Endpoints for associating and disassociating tags with project items")
public interface ItemTagApi {

    @Operation(
        summary = "Associate a tag with an item",
        description = "Associates the specified tag ID with a project item (snippet, credential, link, note, problem)."
    )
    @ApiResponses(value = {
        @ApiResponse(
            responseCode = "204",
            description = "Association completed (returns 204 No Content)",
            content = @Content
        ),
        @ApiResponse(
            responseCode = "400",
            description = "Invalid request or unsupported item type",
            content = @Content(
                mediaType = "application/json",
                schema = @Schema(implementation = ApiErrorResponse.class)
            )
        ),
        @ApiResponse(
            responseCode = "404",
            description = "Project, Item, or Tag not found",
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
    @PutMapping("/{tagId}")
    ResponseEntity<Void> associate(
            @PathVariable UUID projectId,
            @PathVariable String itemType,
            @PathVariable UUID itemId,
            @PathVariable UUID tagId
    );

    @Operation(
        summary = "Disassociate a tag from an item",
        description = "Removes the association between the specified tag ID and a project item."
    )
    @ApiResponses(value = {
        @ApiResponse(
            responseCode = "204",
            description = "Association removed successfully",
            content = @Content
        ),
        @ApiResponse(
            responseCode = "400",
            description = "Invalid request or unsupported item type",
            content = @Content(
                mediaType = "application/json",
                schema = @Schema(implementation = ApiErrorResponse.class)
            )
        ),
        @ApiResponse(
            responseCode = "404",
            description = "Project, Item, or Tag association not found",
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
    ResponseEntity<Void> disassociate(
            @PathVariable UUID projectId,
            @PathVariable String itemType,
            @PathVariable UUID itemId,
            @PathVariable UUID tagId
    );
}
