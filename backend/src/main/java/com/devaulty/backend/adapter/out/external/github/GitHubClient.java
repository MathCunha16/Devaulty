package com.devaulty.backend.adapter.out.external.github;

import com.devaulty.backend.adapter.out.external.github.dto.GitHubReleaseResponse;
import org.springframework.core.io.buffer.DataBuffer;
import org.springframework.http.HttpStatus;
import org.springframework.stereotype.Component;
import org.springframework.web.reactive.function.client.WebClient;
import org.springframework.web.reactive.function.client.WebClientResponseException;
import reactor.core.publisher.Flux;
import reactor.core.publisher.Mono;

import java.net.URI;

@Component
public class GitHubClient {

    private final WebClient githubWebClient;

    private static final String LATEST_RELEASES_URL = "/repos/MathCunha16/Devaulty/releases/latest";

    public GitHubClient(WebClient githubWebClient) {
        this.githubWebClient = githubWebClient;
    }

    public GitHubReleaseResponse getLatestRelease() {
        return githubWebClient.get()
                .uri(LATEST_RELEASES_URL)
                .retrieve()
                .bodyToMono(GitHubReleaseResponse.class)
                .onErrorResume(
                        WebClientResponseException.class,
                        ex -> ex.getStatusCode() == HttpStatus.NOT_FOUND ? Mono.empty() : Mono.error(ex)
                )
                .block();
    }

    public Flux<DataBuffer> downloadAsset(String downloadUrl) {
        return githubWebClient.get()
                .uri(URI.create(downloadUrl)) // Clean URL
                .retrieve()
                .bodyToFlux(DataBuffer.class);
    }
}
