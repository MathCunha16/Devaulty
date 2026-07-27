package com.devaulty.backend.adapter.in.web.release;

import com.devaulty.backend.application.port.in.release.InstallUpdateUseCase;
import com.devaulty.backend.application.port.out.external.release.ReleasePort;
import com.devaulty.backend.application.port.out.external.release.dto.LatestReleaseInfo;
import com.devaulty.backend.application.port.out.external.release.dto.ReleaseAssetInfo;
import com.devaulty.backend.infrastructure.BaseIntegrationTest;
import com.devaulty.backend.infrastructure.security.AppTokenContext;
import org.junit.jupiter.api.DisplayName;
import org.junit.jupiter.api.Test;
import org.springframework.http.MediaType;
import org.springframework.test.context.bean.override.mockito.MockitoBean;
import org.springframework.test.web.servlet.MvcResult;

import java.io.ByteArrayInputStream;
import java.io.InputStream;
import java.nio.charset.StandardCharsets;
import java.time.Instant;
import java.util.List;

import static org.hamcrest.Matchers.*;
import static org.mockito.ArgumentMatchers.anyString;
import static org.mockito.Mockito.when;
import static org.springframework.test.web.servlet.request.MockMvcRequestBuilders.get;
import static org.springframework.test.web.servlet.request.MockMvcRequestBuilders.post;
import static org.springframework.test.web.servlet.request.MockMvcRequestBuilders.asyncDispatch;
import static org.springframework.test.web.servlet.result.MockMvcResultMatchers.*;

class ReleaseControllerIT extends BaseIntegrationTest {

    @MockitoBean
    private ReleasePort releasePort;

    @MockitoBean
    private InstallUpdateUseCase installUpdateUseCase;

    @Test
    @DisplayName("GET /api/v1/releases/current-app-version should return 200 OK with actual version payload")
    void getAppInfo_shouldReturn200OK_withActualVersion() throws Exception {
        mockMvc.perform(get("/api/v1/releases/current-app-version")
                        .header(AppTokenContext.HEADER_NAME, AppTokenContext.PROCESS_TOKEN)
                        .accept(MediaType.APPLICATION_JSON))
                .andExpect(status().isOk())
                .andExpect(jsonPath("$.actualVersion", not(emptyOrNullString())));
    }

    @Test
    @DisplayName("GET /api/v1/releases/check should return 200 OK with update response payload")
    void checkUpdates_shouldReturn200OK_withUpdateInfo() throws Exception {
        // Arrange
        String currentOs = System.getProperty("os.name").toLowerCase();
        String extension = currentOs.contains("win") ? ".msi" : currentOs.contains("mac") ? ".dmg" : ".deb";

        ReleaseAssetInfo asset = new ReleaseAssetInfo(
                "devaulty_0.2.0" + extension,
                "https://github.com/MathCunha16/Devaulty/releases/download/v0.2.0/devaulty_0.2.0" + extension,
                50000000L,
                "application/octet-stream"
        );
        LatestReleaseInfo latestRelease = new LatestReleaseInfo(
                "v0.2.0",
                "Release v0.2.0",
                "Changelog content",
                "https://github.com/MathCunha16/Devaulty/releases/tag/v0.2.0",
                Instant.now(),
                false,
                List.of(asset)
        );

        when(releasePort.getLatestRelease()).thenReturn(latestRelease);

        // Act & Assert
        mockMvc.perform(get("/api/v1/releases/check")
                        .header(AppTokenContext.HEADER_NAME, AppTokenContext.PROCESS_TOKEN)
                        .accept(MediaType.APPLICATION_JSON))
                .andExpect(status().isOk())
                .andExpect(jsonPath("$.updateAvailable").value(true))
                .andExpect(jsonPath("$.currentVersion").isNotEmpty())
                .andExpect(jsonPath("$.latestVersion").value("v0.2.0"))
                .andExpect(jsonPath("$.releaseTitle").value("Release v0.2.0"))
                .andExpect(jsonPath("$.releaseNotes").value("Changelog content"))
                .andExpect(jsonPath("$.downloadUrl").value(containsString(extension)))
                .andExpect(jsonPath("$.downloadSizeInBytes").value(50000000L));
    }

    @Test
    @DisplayName("POST /api/v1/releases/download-and-install should return 200 OK with SSE stream")
    void downloadUpdate_shouldReturn200OK_withSSEStream() throws Exception {
        // Arrange
        String currentOs = System.getProperty("os.name").toLowerCase();
        String extension = currentOs.contains("win") ? ".msi" : currentOs.contains("mac") ? ".dmg" : ".deb";
        String downloadUrl = "https://github.com/MathCunha16/Devaulty/releases/download/v0.2.0/devaulty_0.2.0" + extension;

        ReleaseAssetInfo asset = new ReleaseAssetInfo("devaulty_0.2.0" + extension, downloadUrl, 11L, "application/octet-stream");
        LatestReleaseInfo latestRelease = new LatestReleaseInfo(
                "v0.2.0",
                "Release v0.2.0",
                "Changelog content",
                "https://github.com/MathCunha16/Devaulty/releases/tag/v0.2.0",
                Instant.now(),
                false,
                List.of(asset)
        );

        when(releasePort.getLatestRelease()).thenReturn(latestRelease);

        InputStream fakeAssetStream = new ByteArrayInputStream("hello world".getBytes(StandardCharsets.UTF_8));
        when(releasePort.downloadAsset(anyString())).thenReturn(fakeAssetStream);

        MvcResult mvcResult = mockMvc.perform(post("/api/v1/releases/download-and-install")
                        .header(AppTokenContext.HEADER_NAME, AppTokenContext.PROCESS_TOKEN)
                        .accept(MediaType.TEXT_EVENT_STREAM))
                .andExpect(request().asyncStarted())
                .andReturn();
        mvcResult.getAsyncResult();

        mockMvc.perform(asyncDispatch(mvcResult))
                .andExpect(status().isOk())
                .andExpect(content().contentTypeCompatibleWith(MediaType.TEXT_EVENT_STREAM));
    }
}