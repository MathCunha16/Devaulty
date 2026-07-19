package com.devaulty.backend.adapter.in.web.tag.dto;

import java.util.UUID;

public record TagSummaryResponse(
        UUID id,
        String name,
        String color
) {
}
