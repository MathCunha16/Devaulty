package com.devaulty.backend.infrastructure.lifecycle;

import com.devaulty.backend.infrastructure.security.AppTokenContext;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;
import org.springframework.boot.web.server.context.WebServerInitializedEvent;
import org.springframework.context.ApplicationListener;
import org.springframework.stereotype.Component;

@Component
public class RuntimeSessionLogger implements ApplicationListener<WebServerInitializedEvent> {

    private static final Logger log = LoggerFactory.getLogger(RuntimeSessionLogger.class);

    public static final String SESSION_PREFIX = "[DEVAULTY_SESSION]";

    @Override
    public void onApplicationEvent(WebServerInitializedEvent event) {
        int port = event.getWebServer().getPort();
        String token = AppTokenContext.PROCESS_TOKEN;

        String sessionPayload = String.format("%s PORT=%d TOKEN=%s", SESSION_PREFIX, port, token);

        System.out.println(sessionPayload);

        log.info("Devaulty Backend initialized on dynamic port: {}", port);
    }
}
