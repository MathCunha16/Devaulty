package com.devaulty.backend.adapter.out.external.github.common;

import io.netty.channel.ChannelOption;
import io.netty.handler.timeout.ReadTimeoutHandler;
import io.netty.handler.timeout.WriteTimeoutHandler;
import org.springframework.context.annotation.Bean;
import org.springframework.context.annotation.Configuration;
import org.springframework.http.client.reactive.ReactorClientHttpConnector;
import org.springframework.web.reactive.function.client.WebClient;
import reactor.netty.http.client.HttpClient;

import java.time.Duration;
import java.util.concurrent.TimeUnit;

@Configuration
public class GitHubConfig {

    private static final int CONNECT_TIMEOUT_MS = 10000;
    private static final int API_READ_WRITE_TIMEOUT_SECONDS = 15;

    /** Generous timeout for downloading large binary assets (e.g. .deb, .rpm, .msi). */
    private static final int DOWNLOAD_READ_WRITE_TIMEOUT_SECONDS = 600; // 10 minutes

    /**
     * WebClient for GitHub REST API calls (JSON responses, short timeout).
     */
    @Bean(name = "githubWebClient")
    public WebClient githubWebClient() {
        HttpClient httpClient = HttpClient.create()
                .option(ChannelOption.CONNECT_TIMEOUT_MILLIS, CONNECT_TIMEOUT_MS)
                .responseTimeout(Duration.ofSeconds(API_READ_WRITE_TIMEOUT_SECONDS))
                .doOnConnected(conn -> conn
                        .addHandlerLast(new ReadTimeoutHandler(API_READ_WRITE_TIMEOUT_SECONDS, TimeUnit.SECONDS))
                        .addHandlerLast(new WriteTimeoutHandler(API_READ_WRITE_TIMEOUT_SECONDS, TimeUnit.SECONDS)))
                .followRedirect(true);

        return WebClient.builder()
                .clientConnector(new ReactorClientHttpConnector(httpClient))
                .baseUrl("https://api.github.com")
                .defaultHeader("User-Agent", "Devaulty-Desktop-App")
                .defaultHeader("Accept", "application/vnd.github+json")
                .defaultHeader("X-GitHub-Api-Version", "2026-03-10")
                .build();
    }

    /**
     * WebClient for downloading binary release assets from GitHub CDN.
     * No base URL, no GitHub API Accept header, and a much longer timeout
     * to prevent large installer files from being truncated mid-download.
     */
    @Bean(name = "downloadWebClient")
    public WebClient downloadWebClient() {
        HttpClient httpClient = HttpClient.create()
                .option(ChannelOption.CONNECT_TIMEOUT_MILLIS, CONNECT_TIMEOUT_MS)
                .responseTimeout(Duration.ofSeconds(DOWNLOAD_READ_WRITE_TIMEOUT_SECONDS))
                .doOnConnected(conn -> conn
                        .addHandlerLast(new ReadTimeoutHandler(DOWNLOAD_READ_WRITE_TIMEOUT_SECONDS, TimeUnit.SECONDS))
                        .addHandlerLast(new WriteTimeoutHandler(DOWNLOAD_READ_WRITE_TIMEOUT_SECONDS, TimeUnit.SECONDS)))
                .followRedirect(true);

        return WebClient.builder()
                .clientConnector(new ReactorClientHttpConnector(httpClient))
                .defaultHeader("User-Agent", "Devaulty-Desktop-App")
                .build();
    }
}
