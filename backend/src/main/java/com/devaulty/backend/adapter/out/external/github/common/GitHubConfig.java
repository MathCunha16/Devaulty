package com.devaulty.backend.adapter.out.external.github.common;

import org.springframework.context.annotation.Bean;
import org.springframework.context.annotation.Configuration;
import org.springframework.http.client.JdkClientHttpRequestFactory;
import org.springframework.web.client.RestClient;

import java.net.http.HttpClient;
import java.time.Duration;

@Configuration
public class GitHubConfig {

    private static final Duration CONNECT_TIMEOUT = Duration.ofSeconds(10);
    private static final Duration API_READ_TIMEOUT = Duration.ofSeconds(15);
    /**
     * Generous timeout for downloading large binary assets (e.g. .deb, .rpm, .msi).
     */
    private static final Duration DOWNLOAD_READ_TIMEOUT = Duration.ofMinutes(10);
    private static final String USER_AGENT_HEADER = "Devaulty-Desktop-App";
    private static final String GITHUB_API_BASE_URL = "https://api.github.com";
    private static final String GITHUB_ACCEPT_HEADER = "application/vnd.github+json";
    private static final String GITHUB_API_VERSION_HEADER = "2026-03-10";


    /**
     * Default RestClient for GitHub API requests.
     */
    @Bean
    public HttpClient sharedHttpClient() {
        return HttpClient.newBuilder()
                .connectTimeout(CONNECT_TIMEOUT)
                .followRedirects(HttpClient.Redirect.NORMAL)
                .build();
    }

    @Bean(name = "githubRestClient")
    public RestClient githubRestClient(HttpClient sharedHttpClient) {
        JdkClientHttpRequestFactory requestFactory = new JdkClientHttpRequestFactory(sharedHttpClient);
        requestFactory.setReadTimeout(API_READ_TIMEOUT);

        return RestClient.builder()
                .requestFactory(requestFactory)
                .baseUrl(GITHUB_API_BASE_URL)
                .defaultHeader("User-Agent", USER_AGENT_HEADER)
                .defaultHeader("Accept", GITHUB_ACCEPT_HEADER)
                .defaultHeader("X-GitHub-Api-Version", GITHUB_API_VERSION_HEADER)
                .build();
    }


    /**
     * RestClient for downloading large binary assets (e.g. .deb, .rpm, .msi).
     */
    @Bean(name = "downloadRestClient")
    public RestClient downloadRestClient(HttpClient sharedHttpClient) {
        JdkClientHttpRequestFactory requestFactory = new JdkClientHttpRequestFactory(sharedHttpClient);
        requestFactory.setReadTimeout(DOWNLOAD_READ_TIMEOUT);

        return RestClient.builder()
                .requestFactory(requestFactory)
                .defaultHeader("User-Agent", USER_AGENT_HEADER)
                .build();
    }

}
