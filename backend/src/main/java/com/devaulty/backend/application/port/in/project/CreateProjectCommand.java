package com.devaulty.backend.application.port.in.project;

public record CreateProjectCommand(
        String name,
        String description,
        String icon,
        String color
) {
}
