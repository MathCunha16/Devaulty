package com.devaulty.backend.adapter.in.web.credential;

import com.devaulty.backend.adapter.in.web.credential.dto.CreateCredentialRequest;
import com.devaulty.backend.adapter.in.web.credential.dto.CredentialViewResponse;
import com.devaulty.backend.adapter.in.web.util.UriLocationBuilderHelper;
import com.devaulty.backend.application.port.in.credential.CreateCredentialCommand;
import com.devaulty.backend.application.port.in.credential.CreateCredentialUseCase;
import com.devaulty.backend.application.port.in.credential.DecryptedCredential;
import jakarta.validation.Valid;
import org.springframework.http.ResponseEntity;
import org.springframework.web.bind.annotation.*;

import java.net.URI;
import java.util.UUID;

@RestController
@RequestMapping("/ap1/v1/projects/{projectId}/credentials")
public class CredentialController implements CredentialApi{

    private final CreateCredentialUseCase createCredentialUseCase;
    private final CredentialWebMapper mapper;
    private final UriLocationBuilderHelper uriLocationBuilderHelper;

    public CredentialController(CreateCredentialUseCase createCredentialUseCase, CredentialWebMapper mapper, UriLocationBuilderHelper uriLocationBuilderHelper) {
        this.createCredentialUseCase = createCredentialUseCase;
        this.mapper = mapper;
        this.uriLocationBuilderHelper = uriLocationBuilderHelper;
    }

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

}
