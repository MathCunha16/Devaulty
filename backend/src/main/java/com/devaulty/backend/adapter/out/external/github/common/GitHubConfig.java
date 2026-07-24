package com.devaulty.backend.adapter.out.external.github.common;

import org.springframework.context.annotation.Bean;
import org.springframework.context.annotation.Configuration;
import org.springframework.web.reactive.function.client.WebClient;

@Configuration
public class GitHubConfig {

    @Bean
    public WebClient githubWebClient() {
        return WebClient.builder()
                .baseUrl("https://api.github.com")
                .defaultHeader("User-Agent", "Devaulty-Desktop-App")
                .defaultHeader("Accept", "application/vnd.github+json")
                .defaultHeader("X-GitHub-Api-Version", "2026-03-10")
                .build();
    }
}
