package com.devaulty.backend.adapter.out.persistence.credential;

import com.devaulty.backend.domain.model.Credential;
import org.mapstruct.Mapper;
import org.mapstruct.Mapping;

@Mapper(componentModel = "spring")
public interface CredentialMapper {

    @Mapping( target = "project", ignore = true)
    CredentialEntity toEntity(Credential credential);

    @Mapping( target = "projectId", source = "project.id")
    Credential toDomain(CredentialEntity credentialEntity);

}
