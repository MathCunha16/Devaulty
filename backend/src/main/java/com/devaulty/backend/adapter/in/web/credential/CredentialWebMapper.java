package com.devaulty.backend.adapter.in.web.credential;

import com.devaulty.backend.adapter.in.web.credential.dto.CreateCredentialRequest;
import com.devaulty.backend.adapter.in.web.credential.dto.CredentialSummaryResponse;
import com.devaulty.backend.adapter.in.web.credential.dto.CredentialViewResponse;
import com.devaulty.backend.adapter.in.web.credential.dto.UpdateCredentialRequest;
import com.devaulty.backend.application.exception.JsonProcessingException;
import com.devaulty.backend.application.port.in.credential.CreateCredentialCommand;
import com.devaulty.backend.application.port.in.credential.CredentialSummary;
import com.devaulty.backend.application.port.in.credential.DecryptedCredential;
import com.devaulty.backend.application.port.in.credential.UpdateCredentialCommand;
import com.devaulty.backend.domain.model.enums.CredentialSecretType;
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

    public abstract CredentialSummaryResponse toSummaryResponse(CredentialSummary credential);

    public CreateCredentialCommand toCreateCredentialCommand(CreateCredentialRequest request, UUID projectId) {
        char[] serializedPayload = null;
        try {
            if (request.secretType() != null) {
                serializedPayload = serializeSecretStructure(
                        request.secretType(),
                        request.username(),
                        request.password(),
                        request.apiKey(),
                        request.rawTextContent()
                );
            }

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
            wipeSensitiveFields(request.username(), request.password(), request.apiKey(), request.rawTextContent());
        }
    }

    public UpdateCredentialCommand toUpdateCredentialCommand(UpdateCredentialRequest request, UUID projectId, UUID id) {
        char[] serializedPayload = null;
        try {
            if (request.secretType() != null) {
                serializedPayload = serializeSecretStructure(
                        request.secretType(),
                        request.username(),
                        request.password(),
                        request.apiKey(),
                        request.rawTextContent()
                );
            }

            return new UpdateCredentialCommand(
                    id,
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
            wipeSensitiveFields(request.username(), request.password(), request.apiKey(), request.rawTextContent());
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

    // --- REUSABLE PRIVATE METHODS ---

    private char[] serializeSecretStructure(CredentialSecretType type, char[] username, char[] password, char[] apiKey, char[] rawText) {
        Object secretStructure = switch (type) {
            case LOGIN -> Map.of(
                    "username", username != null ? String.valueOf(username) : "",
                    "password", password != null ? String.valueOf(password) : ""
            );
            case API_KEY -> Map.of(
                    "apiKey", apiKey != null ? String.valueOf(apiKey) : ""
            );
            case RAW_TEXT -> Map.of(
                    "rawText", rawText != null ? String.valueOf(rawText) : ""
            );
        };
        return jsonMapper.writeValueAsString(secretStructure).toCharArray();
    }

    private void wipeSensitiveFields(char[] username, char[] password, char[] apiKey, char[] rawText) {
        if (username != null) Arrays.fill(username, '\0');
        if (password != null) Arrays.fill(password, '\0');
        if (apiKey != null) Arrays.fill(apiKey, '\0');
        if (rawText != null) Arrays.fill(rawText, '\0');
    }

}