package com.devaulty.backend.adapter.out.external.github;

import com.devaulty.backend.adapter.out.external.github.dto.GitHubReleaseResponse;
import org.springframework.beans.factory.annotation.Qualifier;
import org.springframework.http.HttpStatus;
import org.springframework.stereotype.Component;
import org.springframework.web.client.RestClient;
import org.springframework.web.client.RestClientResponseException;

import java.io.InputStream;
import java.net.URI;

@Component
public class GitHubClient {

    private final RestClient githubRestClient;
    private final RestClient downloadRestClient;

    private static final String LATEST_RELEASES_URL = "/repos/MathCunha16/Devaulty/releases/latest";

    public GitHubClient(@Qualifier("githubRestClient") RestClient githubRestClient,
                        @Qualifier("downloadRestClient") RestClient downloadRestClient) {
        this.githubRestClient = githubRestClient;
        this.downloadRestClient = downloadRestClient;
    }

    public GitHubReleaseResponse getLatestRelease() {
        try {
            return githubRestClient.get()
                    .uri(LATEST_RELEASES_URL)
                    .retrieve()
                    .body(GitHubReleaseResponse.class);
        } catch (RestClientResponseException ex) {
            if (ex.getStatusCode() == HttpStatus.NOT_FOUND) {
                return null;
            }
            throw ex;
        }
    }

    /**
     * Downloads a binary release asset from an absolute URL (e.g. GitHub CDN).
     * Returns the response body as a live InputStream (streaming) — the caller
     * MUST consume it inside a try-with-resources block, otherwise the
     * underlying HTTP connection will leak.
     */
    public InputStream downloadAsset(String downloadUrl) {
        return downloadRestClient.get()
                .uri(URI.create(downloadUrl))
                .exchange((request, response) -> {
                    if (response.getStatusCode().isError()) {
                        throw new RestClientResponseException(
                                "Failed to download asset: HTTP " + response.getStatusCode(),
                                response.getStatusCode().value(),
                                response.getStatusText(),
                                response.getHeaders(),
                                null,
                                null
                        );
                    }
                    return response.getBody();
                }, false);
    }
}