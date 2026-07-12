package com.devaulty.backend.adapter.in.web.problem;

import com.devaulty.backend.adapter.in.web.exception.ApiErrorResponse;
import com.devaulty.backend.adapter.in.web.problem.dto.*;
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

@Tag(name = "Problems", description = "Endpoints for managing problems/errors encountered within a project")
public interface ProblemApi {

    @Operation(
        summary = "Create a new problem log",
        description = "Logs a new problem associated with the specified project ID."
    )
    @ApiResponses(value = {
        @ApiResponse(
            responseCode = "201",
            description = "Problem logged successfully",
            content = @Content(
                mediaType = "application/json",
                schema = @Schema(implementation = ProblemViewResponse.class)
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
    ResponseEntity<ProblemViewResponse> create(
            @PathVariable UUID projectId,
            @RequestBody @Valid CreateProblemRequest request
    );

    @Operation(
        summary = "Get all problems by project",
        description = "Retrieves a paginated list of all problem summaries belonging to the specified project ID."
    )
    @ApiResponses(value = {
        @ApiResponse(
            responseCode = "200",
            description = "Successfully retrieved list of problem summaries"
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
    ResponseEntity<Page<ProblemSummaryResponse>> getAllByProject(
            @PathVariable UUID projectId,
            @RequestParam(defaultValue = "0") int page,
            @RequestParam(defaultValue = "10") int size
    );

    @Operation(
        summary = "Get problem details by ID",
        description = "Retrieves details of a single problem by its unique identifier and its project ID."
    )
    @ApiResponses(value = {
        @ApiResponse(
            responseCode = "200",
            description = "Successfully retrieved the problem details",
            content = @Content(
                mediaType = "application/json",
                schema = @Schema(implementation = ProblemViewResponse.class)
            )
        ),
        @ApiResponse(
            responseCode = "404",
            description = "Project or Problem not found",
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
    @GetMapping("/{problemId}")
    ResponseEntity<ProblemViewResponse> getById(
            @PathVariable UUID projectId,
            @PathVariable UUID problemId
    );

    @Operation(
        summary = "Update a problem",
        description = "Updates details of an existing problem by its unique identifier and its project ID."
    )
    @ApiResponses(value = {
        @ApiResponse(
            responseCode = "200",
            description = "Problem updated successfully",
            content = @Content(
                mediaType = "application/json",
                schema = @Schema(implementation = ProblemViewResponse.class)
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
            description = "Project or Problem not found",
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
    @PatchMapping("/{problemId}")
    ResponseEntity<ProblemViewResponse> update(
            @PathVariable UUID projectId,
            @PathVariable UUID problemId,
            @RequestBody @Valid UpdateProblemRequest request
    );

    @Operation(
        summary = "Update problem status",
        description = "Updates status of an existing problem by its unique identifier and its project ID."
    )
    @ApiResponses(value = {
        @ApiResponse(
            responseCode = "200",
            description = "Problem status updated successfully",
            content = @Content(
                mediaType = "application/json",
                schema = @Schema(implementation = ProblemViewResponse.class)
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
            description = "Project or Problem not found",
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
    @PatchMapping("/{problemId}/status")
    ResponseEntity<ProblemViewResponse> updateStatus(
            @PathVariable UUID projectId,
            @PathVariable UUID problemId,
            @RequestBody @Valid UpdateProblemStatusRequest request
    );

    @Operation(
        summary = "Delete a problem log",
        description = "Deletes an existing problem log by its unique identifier and its project ID."
    )
    @ApiResponses(value = {
        @ApiResponse(
            responseCode = "204",
            description = "Problem deleted successfully"
        ),
        @ApiResponse(
            responseCode = "404",
            description = "Project or Problem not found",
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
    @DeleteMapping("/{problemId}")
    ResponseEntity<Void> delete(
            @PathVariable UUID projectId,
            @PathVariable UUID problemId
    );
}
