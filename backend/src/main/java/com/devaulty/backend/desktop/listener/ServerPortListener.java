package com.devaulty.backend.desktop.listener;

import org.springframework.boot.web.server.context.WebServerInitializedEvent;
import org.springframework.context.ApplicationListener;
import org.springframework.stereotype.Component;

import java.util.concurrent.CompletableFuture;

@Component
public class ServerPortListener implements ApplicationListener<WebServerInitializedEvent> {

    public static final CompletableFuture<Integer> serverPortFuture = new CompletableFuture<>();

    @Override
    public void onApplicationEvent(WebServerInitializedEvent event) {
        int port = event.getWebServer().getPort();
        serverPortFuture.complete(port);
    }
}
