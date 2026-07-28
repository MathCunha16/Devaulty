package com.devaulty.backend.adapter.in.web.common;

import jakarta.annotation.PreDestroy;
import org.springframework.stereotype.Component;

import java.util.concurrent.ExecutorService;
import java.util.concurrent.Executors;
import java.util.concurrent.TimeUnit;

/**
 * Runs long-lived, blocking background tasks (e.g. file downloads) outside
 * of the servlet request thread pool, so HTTP worker threads are never held
 * hostage by slow I/O operations.
 */
@Component
public class BackgroundTaskRunner {

    private final ExecutorService executor = Executors.newSingleThreadExecutor();

    public void run(Runnable task) {
        executor.execute(task);
    }

    @PreDestroy
    public void shutdown() {
        executor.shutdown();
        try {
            if (!executor.awaitTermination(30, TimeUnit.SECONDS)) {
                executor.shutdownNow();
            }
        } catch (InterruptedException e) {
            executor.shutdownNow();
            Thread.currentThread().interrupt();
        }
    }
}