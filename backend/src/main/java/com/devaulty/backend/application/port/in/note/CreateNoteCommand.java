package com.devaulty.backend.application.port.in.note;

import java.util.UUID;

public record CreateNoteCommand(
        UUID projectId,
        String title,
        String content
) {
}
