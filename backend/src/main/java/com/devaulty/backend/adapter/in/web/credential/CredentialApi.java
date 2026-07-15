package com.devaulty.backend.adapter.in.web.credential;

import com.devaulty.backend.adapter.in.web.exception.ApiErrorResponse;
import com.devaulty.backend.adapter.in.web.credential.dto.CreateCredentialRequest;
import com.devaulty.backend.adapter.in.web.credential.dto.CredentialSummaryResponse;
import com.devaulty.backend.adapter.in.web.credential.dto.CredentialViewResponse;
import com.devaulty.backend.adapter.in.web.credential.dto.UpdateCredentialRequest;
import com.devaulty.backend.domain.model.enums.CredentialSecretType;
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

@Tag(name = "Credentials", description = "Endpoints for managing secure credentials inside projects")
public interface CredentialApi {

    @Schema(name = "CreateCredentialRequest", description = "Payload for creating a new credential")
    record CreateCredentialRequestDoc(
        @Schema(description = "Title of the credential", minLength = 2, maxLength = 255, example = "My Database Password")
        String title,

        @Schema(description = "Type of the credential secret")
        CredentialSecretType secretType,

        @Schema(type = "string", description = "Username for LOGIN type", example = "admin")
        String username,

        @Schema(type = "string", description = "Password for LOGIN type", example = "supersecret")
        String password,

        @Schema(type = "string", description = "API Key for API_KEY type", example = "api_key_12345")
        String apiKey,

        @Schema(type = "string", description = "Raw text content for RAW_TEXT type", example = "my secret private key content")
        String rawTextContent,

        @Schema(description = "Optional additional notes", example = "Used for staging env")
        String notes,

        @Schema(description = "Optional related website URL", example = "https://database.internal")
        String relatedUrl
    ) {}

    @Schema(name = "UpdateCredentialRequest", description = "Payload for updating a credential")
    record UpdateCredentialRequestDoc(
        @Schema(description = "Updated title of the credential", minLength = 2, maxLength = 255, example = "My Database Password")
        String title,

        @Schema(description = "Updated type of the credential secret")
        CredentialSecretType secretType,

        @Schema(type = "string", description = "Updated username for LOGIN type", example = "admin")
        String username,

        @Schema(type = "string", description = "Updated password for LOGIN type", example = "supersecret")
        String password,

        @Schema(type = "string", description = "Updated API Key for API_KEY type", example = "api_key_12345")
        String apiKey,

        @Schema(type = "string", description = "Updated raw text content for RAW_TEXT type", example = "my secret private key content")
        String rawTextContent,

        @Schema(description = "Updated notes", example = "Used for staging env")
        String notes,

        @Schema(description = "Updated related website URL", example = "https://database.internal")
        String relatedUrl
    ) {}

    @Operation(
        summary = "Create a new credential",
        description = "Creates and encrypts a new credential under the specified project ID. The vault must be unlocked."
    )
    @ApiResponses(value = {
        @ApiResponse(
            responseCode = "201",
            description = "Credential created successfully",
            content = @Content(
                mediaType = "application/json",
                schema = @Schema(implementation = CredentialViewResponse.class)
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
            responseCode = "403",
            description = "Master password not configured",
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
            responseCode = "423",
            description = "Vault is locked (master key session is inactive)",
            content = @Content(
                mediaType = "application/json",
                schema = @Schema(implementation = ApiErrorResponse.class)
            )
        ),
        @ApiResponse(
            responseCode = "500",
            description = "Internal server error or encryption failure",
            content = @Content(
                mediaType = "application/json",
                schema = @Schema(implementation = ApiErrorResponse.class)
            )
        )
    })
    @PostMapping
    ResponseEntity<CredentialViewResponse> create(
            @PathVariable UUID projectId,
            @io.swagger.v3.oas.annotations.parameters.RequestBody(
                description = "Credential creation details",
                required = true,
                content = @Content(
                    mediaType = "application/json",
                    schema = @Schema(implementation = CreateCredentialRequestDoc.class)
                )
            )
            @RequestBody @Valid CreateCredentialRequest request
    );

    @Operation(
        summary = "Get all credentials by project",
        description = "Retrieves a paginated list of credential summaries belonging to the specified project. The vault must be unlocked."
    )
    @ApiResponses(value = {
        @ApiResponse(
            responseCode = "200",
            description = "Successfully retrieved list of credential summaries"
        ),
        @ApiResponse(
            responseCode = "403",
            description = "Master password not configured",
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
            responseCode = "423",
            description = "Vault is locked (master key session is inactive)",
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
    ResponseEntity<Page<CredentialSummaryResponse>> getAllByProject(
            @PathVariable UUID projectId,
            @RequestParam(defaultValue = "0") int page,
            @RequestParam(defaultValue = "10") int size
    );

    @Operation(
        summary = "Get credential by ID",
        description = "Retrieves and decrypts a single credential by its ID within a specific project. The vault must be unlocked."
    )
    @ApiResponses(value = {
        @ApiResponse(
            responseCode = "200",
            description = "Successfully retrieved and decrypted the credential",
            content = @Content(
                mediaType = "application/json",
                schema = @Schema(implementation = CredentialViewResponse.class)
            )
        ),
        @ApiResponse(
            responseCode = "403",
            description = "Master password not configured",
            content = @Content(
                mediaType = "application/json",
                schema = @Schema(implementation = ApiErrorResponse.class)
            )
        ),
        @ApiResponse(
            responseCode = "404",
            description = "Project or Credential not found (or credential does not belong to the project)",
            content = @Content(
                mediaType = "application/json",
                schema = @Schema(implementation = ApiErrorResponse.class)
            )
        ),
        @ApiResponse(
            responseCode = "423",
            description = "Vault is locked (master key session is inactive)",
            content = @Content(
                mediaType = "application/json",
                schema = @Schema(implementation = ApiErrorResponse.class)
            )
        ),
        @ApiResponse(
            responseCode = "500",
            description = "Internal server error or decryption failure",
            content = @Content(
                mediaType = "application/json",
                schema = @Schema(implementation = ApiErrorResponse.class)
            )
        )
    })
    @GetMapping("/{credentialId}")
    ResponseEntity<CredentialViewResponse> getById(
            @PathVariable UUID projectId,
            @PathVariable UUID credentialId
    );

    @Operation(
        summary = "Update a credential",
        description = "Updates metadata and/or the secret payload of an existing credential. Re-encrypts the new secret payload if provided. The vault must be unlocked."
    )
    @ApiResponses(value = {
        @ApiResponse(
            responseCode = "200",
            description = "Credential updated successfully",
            content = @Content(
                mediaType = "application/json",
                schema = @Schema(implementation = CredentialViewResponse.class)
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
            responseCode = "403",
            description = "Master password not configured",
            content = @Content(
                mediaType = "application/json",
                schema = @Schema(implementation = ApiErrorResponse.class)
            )
        ),
        @ApiResponse(
            responseCode = "404",
            description = "Project or Credential not found (or credential does not belong to the project)",
            content = @Content(
                mediaType = "application/json",
                schema = @Schema(implementation = ApiErrorResponse.class)
            )
        ),
        @ApiResponse(
            responseCode = "423",
            description = "Vault is locked (master key session is inactive)",
            content = @Content(
                mediaType = "application/json",
                schema = @Schema(implementation = ApiErrorResponse.class)
            )
        ),
        @ApiResponse(
            responseCode = "500",
            description = "Internal server error or encryption failure",
            content = @Content(
                mediaType = "application/json",
                schema = @Schema(implementation = ApiErrorResponse.class)
            )
        )
    })
    @PatchMapping("/{credentialId}")
    ResponseEntity<CredentialViewResponse> update(
            @PathVariable UUID projectId,
            @PathVariable UUID credentialId,
            @io.swagger.v3.oas.annotations.parameters.RequestBody(
                description = "Credential update details",
                required = true,
                content = @Content(
                    mediaType = "application/json",
                    schema = @Schema(implementation = UpdateCredentialRequestDoc.class)
                )
            )
            @RequestBody @Valid UpdateCredentialRequest request
    );

    @Operation(
        summary = "Delete a credential",
        description = "Deletes a credential by its ID within a specific project. The vault must be unlocked."
    )
    @ApiResponses(value = {
        @ApiResponse(
            responseCode = "204",
            description = "Credential deleted successfully"
        ),
        @ApiResponse(
            responseCode = "403",
            description = "Master password not configured",
            content = @Content(
                mediaType = "application/json",
                schema = @Schema(implementation = ApiErrorResponse.class)
            )
        ),
        @ApiResponse(
            responseCode = "404",
            description = "Project or Credential not found (or credential does not belong to the project)",
            content = @Content(
                mediaType = "application/json",
                schema = @Schema(implementation = ApiErrorResponse.class)
            )
        ),
        @ApiResponse(
            responseCode = "423",
            description = "Vault is locked (master key session is inactive)",
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
    @DeleteMapping("/{credentialId}")
    ResponseEntity<Void> delete(
            @PathVariable UUID projectId,
            @PathVariable UUID credentialId
    );
}
