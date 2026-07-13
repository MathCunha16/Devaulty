package com.devaulty.backend.application.port.in.note;

import java.util.UUID;

public record UpdateNoteCommand(
        UUID id,
        UUID projectId,
        String title,
        String content
) {
}
