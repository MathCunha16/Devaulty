package com.devaulty.backend.adapter.in.web.credential;

import com.devaulty.backend.adapter.in.web.credential.dto.CreateCredentialRequest;
import com.devaulty.backend.adapter.in.web.credential.dto.CredentialSummaryResponse;
import com.devaulty.backend.adapter.in.web.credential.dto.CredentialViewResponse;
import com.devaulty.backend.adapter.in.web.credential.dto.UpdateCredentialRequest;
import com.devaulty.backend.adapter.in.web.util.UriLocationBuilderHelper;
import com.devaulty.backend.application.port.in.credential.*;
import jakarta.validation.Valid;
import org.springframework.data.domain.Page;
import org.springframework.http.ResponseEntity;
import org.springframework.web.bind.annotation.*;

import java.net.URI;
import java.util.UUID;

@RestController
@RequestMapping("/ap1/v1/projects/{projectId}/credentials")
public class CredentialController implements CredentialApi{

    private final CreateCredentialUseCase createCredentialUseCase;
    private final GetAllCredentialsByProjectUseCase getAllCredentialsByProjectUseCase;
    private final GetCredentialByIdUseCase getCredentialByIdUseCase;
    private final UpdateCredentialUseCase updateCredentialUseCase;
    private final DeleteCredentialUseCase deleteCredentialUseCase;
    private final CredentialWebMapper mapper;
    private final UriLocationBuilderHelper uriLocationBuilderHelper;

    public CredentialController(CreateCredentialUseCase createCredentialUseCase, GetAllCredentialsByProjectUseCase getAllCredentialsByProjectUseCase, GetCredentialByIdUseCase getCredentialByIdUseCase, UpdateCredentialUseCase updateCredentialUseCase, DeleteCredentialUseCase deleteCredentialUseCase, CredentialWebMapper mapper, UriLocationBuilderHelper uriLocationBuilderHelper) {
        this.createCredentialUseCase = createCredentialUseCase;
        this.getAllCredentialsByProjectUseCase = getAllCredentialsByProjectUseCase;
        this.getCredentialByIdUseCase = getCredentialByIdUseCase;
        this.updateCredentialUseCase = updateCredentialUseCase;
        this.deleteCredentialUseCase = deleteCredentialUseCase;
        this.mapper = mapper;
        this.uriLocationBuilderHelper = uriLocationBuilderHelper;
    }

    @Override
    @PostMapping
    public ResponseEntity<CredentialViewResponse> create(
            @PathVariable UUID projectId,
            @RequestBody @Valid CreateCredentialRequest request
    ){
        CreateCredentialCommand command = mapper.toCreateCredentialCommand(request, projectId);
        DecryptedCredential decryptedCredential = createCredentialUseCase.execute(command);
        URI location = uriLocationBuilderHelper.buildLocationUri(decryptedCredential.id());
        return ResponseEntity.created(location).body(mapper.toViewResponse(decryptedCredential));
    }

    @Override
    @GetMapping
    public ResponseEntity<Page<CredentialSummaryResponse>> getAllByProject(
            @PathVariable UUID projectId,
            @RequestParam(defaultValue = "0") int page,
            @RequestParam(defaultValue = "10") int size
    ){
        Page<CredentialSummary> summary = getAllCredentialsByProjectUseCase.execute(projectId, page, size);
        return ResponseEntity.ok(summary.map(mapper::toSummaryResponse));
    }

    @Override
    @GetMapping("/{credentialId}")
    public ResponseEntity<CredentialViewResponse> getById(
            @PathVariable UUID projectId,
            @PathVariable UUID credentialId
    ){
        DecryptedCredential credential = getCredentialByIdUseCase.execute(projectId, credentialId);
        return ResponseEntity.ok(mapper.toViewResponse(credential));
    }

    @Override
    @PatchMapping("/{credentialId}")
    public ResponseEntity<CredentialViewResponse> update(
            @PathVariable UUID projectId,
            @PathVariable UUID credentialId,
            @RequestBody @Valid UpdateCredentialRequest request
    ){
        UpdateCredentialCommand command = mapper.toUpdateCredentialCommand(request, projectId, credentialId);
        DecryptedCredential credential = updateCredentialUseCase.execute(command);
        return ResponseEntity.ok(mapper.toViewResponse(credential));
    }

    @Override
    @DeleteMapping("/{credentialId}")
    public ResponseEntity<Void> delete(
            @PathVariable UUID projectId,
            @PathVariable UUID credentialId
    ){
        deleteCredentialUseCase.execute(projectId, credentialId);
        return ResponseEntity.noContent().build();
    }

}
