package com.devaulty.backend.application.impl.release;

import com.devaulty.backend.application.port.in.release.GetCurrentVersionUseCase;
import com.devaulty.backend.infrastructure.properties.DevaultyProperties;

public class GetCurrentVersionImpl implements GetCurrentVersionUseCase {

    private final DevaultyProperties devaultyProperties;

    public GetCurrentVersionImpl(DevaultyProperties devaultyProperties) {
        this.devaultyProperties = devaultyProperties;
    }

    @Override
    public String execute() {
        return devaultyProperties.getVersion();
    }
}
