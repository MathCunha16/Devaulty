package com.devaulty.backend.adapter.out.external.github;

import com.devaulty.backend.adapter.out.external.github.dto.GitHubReleaseResponse;
import org.springframework.beans.factory.annotation.Qualifier;
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

    /**
     * A separate WebClient used exclusively for downloading binary release assets.
     * It has no base URL, no GitHub-specific Accept header, and a much longer
     * read/write timeout to handle large file downloads without timing out.
     */
    private final WebClient downloadWebClient;

    private static final String LATEST_RELEASES_URL = "/repos/MathCunha16/Devaulty/releases/latest";

    public GitHubClient(@Qualifier("githubWebClient") WebClient githubWebClient,
                        @Qualifier("downloadWebClient") WebClient downloadWebClient) {
        this.githubWebClient = githubWebClient;
        this.downloadWebClient = downloadWebClient;
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

    /**
     * Downloads a binary release asset from an absolute URL (e.g. GitHub CDN).
     * Uses a dedicated WebClient with no GitHub API headers and a generous
     * timeout so that large installers (.deb, .rpm, .msi) are not truncated.
     */
    public Flux<DataBuffer> downloadAsset(String downloadUrl) {
        return downloadWebClient.get()
                .uri(URI.create(downloadUrl))
                .retrieve()
                .bodyToFlux(DataBuffer.class);
    }
}
