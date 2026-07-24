package com.devaulty.backend.adapter.out.external.github;

import com.devaulty.backend.adapter.out.external.github.dto.GitHubReleaseResponse;
import org.springframework.stereotype.Component;
import org.springframework.web.reactive.function.client.WebClient;
import org.springframework.web.reactive.function.client.WebClientException;
import reactor.core.publisher.Mono;

@Component
public class GitHubReleaseClient {

    private final WebClient githubWebClient;

    private static final String LATEST_RELEASES_URL = "/repos/MathCunha16/Devaulty/releases/latest";

    public GitHubReleaseClient(WebClient githubWebClient) {
        this.githubWebClient = githubWebClient;
    }

    public GitHubReleaseResponse getLatestRelease() {
        return githubWebClient.get()
                .uri(LATEST_RELEASES_URL)
                .retrieve()
                .bodyToMono(GitHubReleaseResponse.class)
                .onErrorResume(WebClientException.class, ex -> Mono.empty())
                .block();
    }
}
