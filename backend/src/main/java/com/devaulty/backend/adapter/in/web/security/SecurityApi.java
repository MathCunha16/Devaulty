package com.devaulty.backend.adapter.in.web.security;

import com.devaulty.backend.adapter.in.web.exception.ApiErrorResponse;
import com.devaulty.backend.adapter.in.web.security.dto.MasterPasswordRequest;
import io.swagger.v3.oas.annotations.Operation;
import io.swagger.v3.oas.annotations.media.Content;
import io.swagger.v3.oas.annotations.media.Schema;
import io.swagger.v3.oas.annotations.responses.ApiResponse;
import io.swagger.v3.oas.annotations.responses.ApiResponses;
import io.swagger.v3.oas.annotations.tags.Tag;
import jakarta.validation.Valid;
import org.springframework.http.ResponseEntity;
import org.springframework.web.bind.annotation.*;

@Tag(name = "Security", description = "Endpoints for master password configuration, vault locking, and unlocking")
public interface SecurityApi {

    @Schema(name = "MasterPasswordRequest", description = "Payload containing the master password")
    record MasterPasswordRequestDoc(
        @Schema(type = "string", minLength = 8, maxLength = 255, description = "The master password", example = "mySuperSecurePassword123")
        String masterPassword
    ) {}

    @Operation(
        summary = "Configure master password",
        description = "Initializes the master password for the vault if it has not been configured yet."
    )
    @ApiResponses(value = {
        @ApiResponse(
            responseCode = "204",
            description = "Master password configured successfully"
        ),
        @ApiResponse(
            responseCode = "400",
            description = "Validation failed (password too short or null)",
            content = @Content(
                mediaType = "application/json",
                schema = @Schema(implementation = ApiErrorResponse.class)
            )
        ),
        @ApiResponse(
            responseCode = "409",
            description = "Master password is already configured",
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
    @PostMapping("/master-password")
    ResponseEntity<Void> setupMasterPassword(
            @io.swagger.v3.oas.annotations.parameters.RequestBody(
                description = "Master password details",
                required = true,
                content = @Content(
                    mediaType = "application/json",
                    schema = @Schema(implementation = MasterPasswordRequestDoc.class)
                )
            )
            @RequestBody @Valid MasterPasswordRequest request
    );

    @Operation(
        summary = "Check master password setup status",
        description = "Checks if the master password has not been configured yet (returns true if setup is required, false if already configured)."
    )
    @ApiResponses(value = {
        @ApiResponse(
            responseCode = "200",
            description = "Status check completed successfully",
            content = @Content(
                mediaType = "application/json",
                schema = @Schema(type = "boolean")
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
    @GetMapping("/master-password/required-status")
    ResponseEntity<Boolean> checkMasterPasswordSetup();

    @Operation(
        summary = "Unlock the vault",
        description = "Unlocks the vault with the master password, activating the cryptographic session in RAM."
    )
    @ApiResponses(value = {
        @ApiResponse(
            responseCode = "200",
            description = "Vault unlocked successfully",
            content = @Content(
                mediaType = "application/json",
                schema = @Schema(type = "boolean", example = "true")
            )
        ),
        @ApiResponse(
            responseCode = "400",
            description = "Validation failed (password too short or null)",
            content = @Content(
                mediaType = "application/json",
                schema = @Schema(implementation = ApiErrorResponse.class)
            )
        ),
        @ApiResponse(
            responseCode = "401",
            description = "Invalid master password",
            content = @Content(
                mediaType = "application/json",
                schema = @Schema(implementation = ApiErrorResponse.class)
            )
        ),
        @ApiResponse(
            responseCode = "403",
            description = "Master password has not been configured yet",
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
    @PostMapping("/vault/unlock")
    ResponseEntity<Boolean> unlockVault(
            @io.swagger.v3.oas.annotations.parameters.RequestBody(
                description = "Master password details",
                required = true,
                content = @Content(
                    mediaType = "application/json",
                    schema = @Schema(implementation = MasterPasswordRequestDoc.class)
                )
            )
            @RequestBody @Valid MasterPasswordRequest request
    );

    @Operation(
        summary = "Lock the vault",
        description = "Locks the vault, clearing all cryptographic keys from memory."
    )
    @ApiResponses(value = {
        @ApiResponse(
            responseCode = "204",
            description = "Vault locked successfully (session cleared)"
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
    @PostMapping("/vault/lock")
    ResponseEntity<Void> lockVault();
}
