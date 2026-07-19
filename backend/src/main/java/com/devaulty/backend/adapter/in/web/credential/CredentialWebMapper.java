package com.devaulty.backend.adapter.in.web.credential;

import com.devaulty.backend.adapter.in.web.credential.dto.CreateCredentialRequest;
import com.devaulty.backend.adapter.in.web.credential.dto.CredentialSummaryResponse;
import com.devaulty.backend.adapter.in.web.credential.dto.CredentialViewResponse;
import com.devaulty.backend.adapter.in.web.credential.dto.UpdateCredentialRequest;
import com.devaulty.backend.adapter.in.web.tag.TagWebMapper;
import com.devaulty.backend.application.exception.BusinessRuleException;
import com.devaulty.backend.application.exception.JsonProcessingException;
import com.devaulty.backend.application.port.in.credential.CreateCredentialCommand;
import com.devaulty.backend.application.port.in.credential.CredentialSummary;
import com.devaulty.backend.application.port.in.credential.DecryptedCredential;
import com.devaulty.backend.application.port.in.credential.UpdateCredentialCommand;
import com.devaulty.backend.domain.model.Tag;
import com.devaulty.backend.domain.model.enums.CredentialSecretType;
import org.mapstruct.Mapper;
import org.mapstruct.Mapping;
import org.mapstruct.Named;
import org.springframework.beans.factory.annotation.Autowired;
import tools.jackson.core.type.TypeReference;
import tools.jackson.databind.json.JsonMapper;

import java.util.Arrays;
import java.util.List;
import java.util.Map;
import java.util.UUID;

@Mapper(componentModel = "spring", uses = TagWebMapper.class)
public abstract class CredentialWebMapper {

    @Autowired
    protected JsonMapper jsonMapper;

    @Mapping(source = "credential.decryptedPayload", target = "decryptedPayload", qualifiedByName = "jsonToMap")
    @Mapping(source = "tags", target = "tags")
    public abstract CredentialViewResponse toViewResponse(DecryptedCredential credential, List<Tag> tags);

    @Mapping(source = "tags", target = "tags")
    public abstract CredentialSummaryResponse toSummaryResponse(CredentialSummary credential, List<Tag> tags);

    public CreateCredentialCommand toCreateCredentialCommand(CreateCredentialRequest request, UUID projectId) {
        char[] serializedPayload = null;
        try {
            if (request.secretType() == null) {
                throw new BusinessRuleException("secretType is required");
            }

            boolean isValid = switch (request.secretType()) {
                case LOGIN -> request.password() != null && request.password().length > 0;
                case API_KEY -> request.apiKey() != null && request.apiKey().length > 0;
                case RAW_TEXT -> request.rawTextContent() != null && request.rawTextContent().length > 0;
            };

            if (!isValid) {
                throw new BusinessRuleException("Required secret fields are missing for the selected secret type: " + request.secretType());
            }

            serializedPayload = serializeSecretStructure(
                    request.secretType(),
                    request.username(),
                    request.password(),
                    request.apiKey(),
                    request.rawTextContent()
            );

            return new CreateCredentialCommand(
                    projectId,
                    request.title(),
                    request.secretType(),
                    serializedPayload,
                    request.notes(),
                    request.relatedUrl()
            );
        } catch (BusinessRuleException e) {
            throw e;
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
                boolean isValid = switch (request.secretType()) {
                    case LOGIN -> request.password() != null && request.password().length > 0;
                    case API_KEY -> request.apiKey() != null && request.apiKey().length > 0;
                    case RAW_TEXT -> request.rawTextContent() != null && request.rawTextContent().length > 0;
                };

                if (!isValid) {
                    throw new BusinessRuleException("Required secret fields are missing for the selected secret type: " + request.secretType());
                }

                serializedPayload = serializeSecretStructure(
                        request.secretType(),
                        request.username(),
                        request.password(),
                        request.apiKey(),
                        request.rawTextContent()
                );
            } else {
                boolean hasAnySecretField = (request.username() != null && request.username().length > 0)
                        || (request.password() != null && request.password().length > 0)
                        || (request.apiKey() != null && request.apiKey().length > 0)
                        || (request.rawTextContent() != null && request.rawTextContent().length > 0);
                if (hasAnySecretField) {
                    throw new BusinessRuleException("secretType is required when secret fields are provided");
                }
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
        } catch (BusinessRuleException e) {
            throw e;
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