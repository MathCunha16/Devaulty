package com.devaulty.backend.adapter.out.external.github;

import com.devaulty.backend.adapter.out.external.github.dto.GitHubAssetResponse;
import com.devaulty.backend.adapter.out.external.github.dto.GitHubReleaseResponse;
import org.junit.jupiter.api.BeforeEach;
import org.junit.jupiter.api.DisplayName;
import org.junit.jupiter.api.Test;
import org.junit.jupiter.api.extension.ExtendWith;
import org.mockito.Mock;
import org.mockito.junit.jupiter.MockitoExtension;
import org.springframework.core.io.buffer.DataBuffer;
import org.springframework.core.io.buffer.DefaultDataBufferFactory;
import org.springframework.web.reactive.function.client.WebClient;
import org.springframework.web.reactive.function.client.WebClientResponseException;
import reactor.core.publisher.Flux;
import reactor.core.publisher.Mono;

import java.net.URI;
import java.nio.charset.StandardCharsets;
import java.time.Instant;
import java.util.List;

import static org.junit.jupiter.api.Assertions.*;
import static org.mockito.ArgumentMatchers.any;
import static org.mockito.ArgumentMatchers.anyString;
import static org.mockito.Mockito.*;

@ExtendWith(MockitoExtension.class)
@SuppressWarnings({"unchecked", "rawtypes"})
class GitHubClientTest {

    @Mock
    private WebClient webClient;

    @Mock
    private WebClient.RequestHeadersUriSpec requestHeadersUriSpec;

    @Mock
    private WebClient.RequestHeadersSpec requestHeadersSpec;

    @Mock
    private WebClient.ResponseSpec responseSpec;

    private GitHubClient gitHubClient;

    @BeforeEach
    void setUp() {
        gitHubClient = new GitHubClient(webClient);
    }

    @Test
    @DisplayName("getLatestRelease should return GitHubReleaseResponse when request succeeds")
    void getLatestRelease_shouldReturnReleaseResponse_whenGitHubReturnsRelease() {
        // Arrange
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

        when(webClient.get()).thenReturn(requestHeadersUriSpec);
        when(requestHeadersUriSpec.uri(anyString())).thenReturn(requestHeadersSpec);
        when(requestHeadersSpec.retrieve()).thenReturn(responseSpec);
        when(responseSpec.bodyToMono(GitHubReleaseResponse.class)).thenReturn(Mono.just(expectedResponse));

        // Act
        GitHubReleaseResponse actualResponse = gitHubClient.getLatestRelease();

        // Assert
        assertNotNull(actualResponse);
        assertEquals("v0.2.0", actualResponse.tagName());
        assertEquals("Release v0.2.0", actualResponse.name());
        assertEquals("Changelog notes", actualResponse.body());
        assertEquals(1, actualResponse.assets().size());

        verify(webClient, times(1)).get();
    }

    @Test
    @DisplayName("getLatestRelease should return null when 404 Not Found WebClientResponseException occurs")
    void getLatestRelease_shouldReturnNull_when404NotFoundOccurs() {
        // Arrange
        when(webClient.get()).thenReturn(requestHeadersUriSpec);
        when(requestHeadersUriSpec.uri(anyString())).thenReturn(requestHeadersSpec);
        when(requestHeadersSpec.retrieve()).thenReturn(responseSpec);
        when(responseSpec.bodyToMono(GitHubReleaseResponse.class))
                .thenReturn(Mono.error(WebClientResponseException.create(404, "Not Found", null, null, null)));

        // Act
        GitHubReleaseResponse actualResponse = gitHubClient.getLatestRelease();

        // Assert
        assertNull(actualResponse);
        verify(webClient, times(1)).get();
    }

    @Test
    @DisplayName("getLatestRelease should rethrow exception when 500 Server Error WebClientResponseException occurs")
    void getLatestRelease_shouldThrowException_whenServerErrorOccurs() {
        // Arrange
        when(webClient.get()).thenReturn(requestHeadersUriSpec);
        when(requestHeadersUriSpec.uri(anyString())).thenReturn(requestHeadersSpec);
        when(requestHeadersSpec.retrieve()).thenReturn(responseSpec);
        when(responseSpec.bodyToMono(GitHubReleaseResponse.class))
                .thenReturn(Mono.error(WebClientResponseException.create(500, "Internal Server Error", null, null, null)));

        // Act & Assert
        assertThrows(WebClientResponseException.class, () -> gitHubClient.getLatestRelease());
        verify(webClient, times(1)).get();
    }

    @Test
    @DisplayName("downloadAsset should return DataBuffer Flux when download URL is valid")
    void downloadAsset_shouldReturnDataBufferFlux_whenUrlIsValid() {
        // Arrange
        String downloadUrl = "https://github.com/MathCunha16/Devaulty/releases/download/v0.2.0/devaulty_0.2.0.deb";
        DefaultDataBufferFactory factory = new DefaultDataBufferFactory();
        DataBuffer mockBuffer = factory.wrap("binary content".getBytes(StandardCharsets.UTF_8));

        when(webClient.get()).thenReturn(requestHeadersUriSpec);
        when(requestHeadersUriSpec.uri(any(URI.class))).thenReturn(requestHeadersSpec);
        when(requestHeadersSpec.retrieve()).thenReturn(responseSpec);
        when(responseSpec.bodyToFlux(DataBuffer.class)).thenReturn(Flux.just(mockBuffer));

        // Act
        Flux<DataBuffer> resultFlux = gitHubClient.downloadAsset(downloadUrl);
        List<DataBuffer> buffers = resultFlux.collectList().block();

        // Assert
        assertNotNull(buffers);
        assertEquals(1, buffers.size());
        verify(webClient, times(1)).get();
    }
}
