package com.devaulty.backend.adapter.in.web.credential;

import com.devaulty.backend.adapter.in.web.credential.dto.CreateCredentialRequest;
import com.devaulty.backend.adapter.in.web.credential.dto.CredentialSummaryResponse;
import com.devaulty.backend.adapter.in.web.credential.dto.CredentialViewResponse;
import com.devaulty.backend.application.exception.JsonProcessingException;
import com.devaulty.backend.application.port.in.credential.CreateCredentialCommand;
import com.devaulty.backend.application.port.in.credential.DecryptedCredential;
import com.devaulty.backend.domain.model.Credential;
import org.mapstruct.Mapper;
import org.mapstruct.Mapping;
import org.mapstruct.Named;
import org.springframework.beans.factory.annotation.Autowired;
import tools.jackson.core.type.TypeReference;
import tools.jackson.databind.json.JsonMapper;

import java.util.Arrays;
import java.util.Map;
import java.util.UUID;

@Mapper(componentModel = "spring")
public abstract class CredentialWebMapper {

    @Autowired
    protected JsonMapper jsonMapper;

    @Mapping(source = "decryptedPayload", target = "decryptedPayload", qualifiedByName = "jsonToMap")
    public abstract CredentialViewResponse toViewResponse(DecryptedCredential credential);

    public abstract CredentialSummaryResponse toSummaryResponse(Credential credential);

    public CreateCredentialCommand toCreateCredentialCommand(CreateCredentialRequest request, UUID projectId) {

        Object secretStructure = switch (request.secretType()) {
            case LOGIN -> Map.of(
                    "username", request.username() != null ? String.valueOf(request.username()) : "",
                    "password", request.password() != null ? String.valueOf(request.password()) : ""
            );

            case API_KEY -> Map.of(
                    "apiKey", request.apiKey() != null ? String.valueOf(request.apiKey()) : ""
            );

            case RAW_TEXT -> Map.of(
                    "rawText", request.rawTextContent() != null ? String.valueOf(request.rawTextContent()) : ""
            );
        };

        char[] serializedPayload = null;
        try {
            String json = jsonMapper.writeValueAsString(secretStructure);
            serializedPayload = json.toCharArray();

            return new CreateCredentialCommand(
                    projectId,
                    request.title(),
                    request.secretType(),
                    serializedPayload,
                    request.notes(),
                    request.relatedUrl()
            );

        } catch (Exception e) {
            throw new JsonProcessingException("Error trying to structure credential secrets.", e);
        } finally {
            // Clears the raw secret fields from memory
            if (request.username() != null) Arrays.fill(request.username(), '\0');
            if (request.password() != null) Arrays.fill(request.password(), '\0');
            if (request.apiKey() != null) Arrays.fill(request.apiKey(), '\0');
            if (request.rawTextContent() != null) Arrays.fill(request.rawTextContent(), '\0');
        }

    }

    @Named("jsonToMap")
    protected Map<String, String> jsonToMap(byte[] decryptedBytes) {
        if (decryptedBytes == null || decryptedBytes.length == 0) return Map.of();

        try {
            return jsonMapper.readValue(decryptedBytes, new TypeReference<>() {
            });
        } catch (Exception e) {
            throw new JsonProcessingException("Failed to parse decrypted payload.", e);
        } finally {
            // Clears the decryptedBytes from memory
            Arrays.fill(decryptedBytes, (byte) 0);
        }
    }
}
