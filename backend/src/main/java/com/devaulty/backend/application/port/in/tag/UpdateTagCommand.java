package com.devaulty.backend.application.port.in.tag;

import java.util.UUID;

public record UpdateTagCommand(
        UUID id,
        UUID projectId,
        String name,
        String color
) {
}
