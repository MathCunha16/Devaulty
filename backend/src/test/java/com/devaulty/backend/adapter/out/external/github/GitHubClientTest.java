package com.devaulty.backend.adapter.out.external.github;

import com.devaulty.backend.adapter.out.external.github.dto.GitHubAssetResponse;
import com.devaulty.backend.adapter.out.external.github.dto.GitHubReleaseResponse;
import org.junit.jupiter.api.BeforeEach;
import org.junit.jupiter.api.DisplayName;
import org.junit.jupiter.api.Test;
import org.junit.jupiter.api.extension.ExtendWith;
import org.mockito.Mock;
import org.mockito.junit.jupiter.MockitoExtension;
import org.springframework.web.client.HttpClientErrorException;
import org.springframework.web.client.HttpServerErrorException;
import org.springframework.web.client.RestClient;
import org.springframework.web.client.RestClientResponseException;

import java.io.ByteArrayInputStream;
import java.io.InputStream;
import java.net.URI;
import java.nio.charset.StandardCharsets;
import java.time.Instant;
import java.util.List;

import static org.junit.jupiter.api.Assertions.*;
import static org.mockito.ArgumentMatchers.*;
import static org.mockito.Mockito.*;

@ExtendWith(MockitoExtension.class)
@SuppressWarnings({"unchecked", "rawtypes"})
class GitHubClientTest {

    @Mock
    private RestClient githubRestClient;

    @Mock
    private RestClient downloadRestClient;

    @Mock
    private RestClient.RequestHeadersUriSpec requestHeadersUriSpec;

    @Mock
    private RestClient.RequestHeadersSpec requestHeadersSpec;

    @Mock
    private RestClient.ResponseSpec responseSpec;

    @Mock
    private RestClient.RequestHeadersUriSpec downloadRequestHeadersUriSpec;

    @Mock
    private RestClient.RequestHeadersSpec downloadRequestHeadersSpec;

    private GitHubClient gitHubClient;

    @BeforeEach
    void setUp() {
        gitHubClient = new GitHubClient(githubRestClient, downloadRestClient);
    }

    @Test
    @DisplayName("getLatestRelease should return GitHubReleaseResponse when request succeeds")
    void getLatestRelease_shouldReturnReleaseResponse_whenGitHubReturnsRelease() {
        GitHubAssetResponse asset = new GitHubAssetResponse("devaulty_0.2.0.deb", "https://download.url", 50000000L, "application/octet-stream", "sha256:12345");
        GitHubReleaseResponse expectedResponse = new GitHubReleaseResponse(
                "v0.2.0",
                "Release v0.2.0",
                "Changelog notes",
                "https://github.com/tag/v0.2.0",
                Instant.now(),
                false,
                List.of(asset)
        );

        when(githubRestClient.get()).thenReturn(requestHeadersUriSpec);
        when(requestHeadersUriSpec.uri(anyString())).thenReturn(requestHeadersSpec);
        when(requestHeadersSpec.retrieve()).thenReturn(responseSpec);
        when(responseSpec.body(GitHubReleaseResponse.class)).thenReturn(expectedResponse);

        GitHubReleaseResponse actualResponse = gitHubClient.getLatestRelease();

        assertNotNull(actualResponse);
        assertEquals("v0.2.0", actualResponse.tagName());
        assertEquals("Release v0.2.0", actualResponse.name());
        assertEquals("Changelog notes", actualResponse.body());
        assertEquals(1, actualResponse.assets().size());

        verify(githubRestClient, times(1)).get();
    }

    @Test
    @DisplayName("getLatestRelease should return null when 404 Not Found RestClientResponseException occurs")
    void getLatestRelease_shouldReturnNull_when404NotFoundOccurs() {
        when(githubRestClient.get()).thenReturn(requestHeadersUriSpec);
        when(requestHeadersUriSpec.uri(anyString())).thenReturn(requestHeadersSpec);
        when(requestHeadersSpec.retrieve()).thenReturn(responseSpec);
        when(responseSpec.body(GitHubReleaseResponse.class))
                .thenThrow(HttpClientErrorException.create(
                        org.springframework.http.HttpStatus.NOT_FOUND,
                        "Not Found", null, null, null
                ));

        GitHubReleaseResponse actualResponse = gitHubClient.getLatestRelease();

        assertNull(actualResponse);
        verify(githubRestClient, times(1)).get();
    }

    @Test
    @DisplayName("getLatestRelease should rethrow exception when 500 Server Error RestClientResponseException occurs")
    void getLatestRelease_shouldThrowException_whenServerErrorOccurs() {
        when(githubRestClient.get()).thenReturn(requestHeadersUriSpec);
        when(requestHeadersUriSpec.uri(anyString())).thenReturn(requestHeadersSpec);
        when(requestHeadersSpec.retrieve()).thenReturn(responseSpec);
        when(responseSpec.body(GitHubReleaseResponse.class))
                .thenThrow(HttpServerErrorException.create(
                        org.springframework.http.HttpStatus.INTERNAL_SERVER_ERROR,
                        "Internal Server Error", null, null, null
                ));

        assertThrows(RestClientResponseException.class, () -> gitHubClient.getLatestRelease());
        verify(githubRestClient, times(1)).get();
    }

    @Test
    @DisplayName("downloadAsset should return an InputStream when download URL is valid")
    void downloadAsset_shouldReturnInputStream_whenUrlIsValid() {
        String downloadUrl = "https://github.com/MathCunha16/Devaulty/releases/download/v0.2.0/devaulty_0.2.0.deb";
        InputStream mockStream = new ByteArrayInputStream("binary content".getBytes(StandardCharsets.UTF_8));

        when(downloadRestClient.get()).thenReturn(downloadRequestHeadersUriSpec);
        when(downloadRequestHeadersUriSpec.uri(any(URI.class))).thenReturn(downloadRequestHeadersSpec);
        when(downloadRequestHeadersSpec.exchange(any(), eq(false))).thenReturn(mockStream);

        InputStream result = gitHubClient.downloadAsset(downloadUrl);

        assertNotNull(result);
        assertSame(mockStream, result);
        verify(downloadRestClient, times(1)).get();
        verify(githubRestClient, never()).get();
    }
}