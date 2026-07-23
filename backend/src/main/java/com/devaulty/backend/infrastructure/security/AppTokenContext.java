package com.devaulty.backend.infrastructure.security;

import java.util.UUID;

public class AppTokenContext {

    public static final String HEADER_NAME = "X-Devaulty-Internal-Token";
    // Unique token generated in memory
    public static final String PROCESS_TOKEN = UUID.randomUUID().toString();

    private AppTokenContext() {
        /* This utility class should not be instantiated */
    }
}
