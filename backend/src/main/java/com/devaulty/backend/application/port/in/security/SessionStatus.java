package com.devaulty.backend.application.port.in.security;

public record SessionStatus(
        boolean active,
        long secondsLeft
) {
}
