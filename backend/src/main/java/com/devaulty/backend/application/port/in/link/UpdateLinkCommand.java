package com.devaulty.backend.application.port.in.link;

import java.util.UUID;

public record UpdateLinkCommand(
        UUID id,
        UUID projectId,
        String title,
        String url,
        String description
) {
}
