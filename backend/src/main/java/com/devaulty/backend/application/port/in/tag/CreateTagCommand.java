package com.devaulty.backend.application.port.in.tag;

import java.util.UUID;

public record CreateTagCommand(
        UUID projectId,
        String name,
        String color
) {
}
