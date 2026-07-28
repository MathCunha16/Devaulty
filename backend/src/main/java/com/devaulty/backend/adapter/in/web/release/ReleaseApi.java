package com.devaulty.backend.adapter.in.web.release;

import com.devaulty.backend.adapter.in.web.exception.ApiErrorResponse;
import com.devaulty.backend.adapter.in.web.release.dto.AppUpdateInfoResponse;
import com.devaulty.backend.adapter.in.web.release.dto.CurrentVersionResponse;
import com.devaulty.backend.adapter.in.web.release.dto.UpdateDownloadProgressResponse;
import io.swagger.v3.oas.annotations.Operation;
import io.swagger.v3.oas.annotations.media.Content;
import io.swagger.v3.oas.annotations.media.Schema;
import io.swagger.v3.oas.annotations.responses.ApiResponse;
import io.swagger.v3.oas.annotations.responses.ApiResponses;
import io.swagger.v3.oas.annotations.tags.Tag;
import org.springframework.http.MediaType;
import org.springframework.http.ResponseEntity;
import org.springframework.web.bind.annotation.GetMapping;
import org.springframework.web.bind.annotation.PostMapping;
import org.springframework.web.bind.annotation.RequestMapping;
import org.springframework.web.servlet.mvc.method.annotation.SseEmitter;

@Tag(name = "Releases", description = "Endpoints for checking desktop application updates, streaming asset downloads, and invoking native OS installation")
@RequestMapping("/api/v1/releases")
public interface ReleaseApi {

    @Operation(
            summary = "Check for application updates",
            description = "Fetches the latest official relea" +
                    "se published on GitHub Releases and compares its tag version with the local running application version (`app.version`). Dynamically detects the current host OS (Linux .deb vs .rpm, Windows .msi, macOS .dmg) and resolves the matching asset download URL and file size."
    )
    @ApiResponses(value = {
            @ApiResponse(
                    responseCode = "200",
                    description = "Update check completed successfully. Returns update availability status, current vs latest version info, changelog notes, and host-specific download metadata.",
                    content = @Content(mediaType = MediaType.APPLICATION_JSON_VALUE, schema = @Schema(implementation = AppUpdateInfoResponse.class))
            ),
            @ApiResponse(
                    responseCode = "403",
                    description = "Forbidden. Request missing or containing an invalid `X-Devaulty-Internal-Token` header.",
                    content = @Content(mediaType = MediaType.APPLICATION_JSON_VALUE, schema = @Schema(implementation = ApiErrorResponse.class))
            ),
            @ApiResponse(
                    responseCode = "500",
                    description = "Internal Server Error. Unexpected error occurred while querying GitHub API or reading local system properties.",
                    content = @Content(mediaType = MediaType.APPLICATION_JSON_VALUE, schema = @Schema(implementation = ApiErrorResponse.class))
            )
    })
    @GetMapping("/check")
    ResponseEntity<AppUpdateInfoResponse> checkUpdates();

    @Operation(
            summary = "Stream update download and trigger native OS installation",
            description = "Establishes a Server-Sent Events (SSE) HTTP stream (`text/event-stream`) to download the latest installer binary into the local platform temp directory (`~/.config/devaulty/temp` on Linux, `%LOCALAPPDATA%\\devaulty\\temp` on Windows, `~/Library/Caches/devaulty/temp` on macOS). Emits real-time progress events (`DOWNLOADING` status with percentage and byte counts). Once download reaches 100%, transitions to `INSTALLING` status, launches the native OS installer/restart script detached from the current process, and gracefully shuts down the running Devaulty application instance. " +
                    "Backed by a blocking download pipeline (JDK `HttpClient`) executed on a dedicated background thread, decoupled from the servlet request thread pool."
    )
    @ApiResponses(value = {
            @ApiResponse(
                    responseCode = "200",
                    description = "Server-Sent Events (SSE) stream established successfully. Emits real-time `UpdateDownloadProgressResponse` events until the stream completes or errors out.",
                    content = @Content(mediaType = MediaType.TEXT_EVENT_STREAM_VALUE, schema = @Schema(implementation = UpdateDownloadProgressResponse.class))
            ),
            @ApiResponse(
                    responseCode = "403",
                    description = "Forbidden. Request missing or containing an invalid `X-Devaulty-Internal-Token` header.",
                    content = @Content(mediaType = MediaType.APPLICATION_JSON_VALUE, schema = @Schema(implementation = ApiErrorResponse.class))
            ),
            @ApiResponse(
                    responseCode = "500",
                    description = "Internal Server Error. Streaming network failure, file I/O error during installer save, or native process execution failure. Reported as a terminal SSE error event, since HTTP headers are already committed by the time streaming begins.",
                    content = @Content(mediaType = MediaType.APPLICATION_JSON_VALUE, schema = @Schema(implementation = ApiErrorResponse.class))
            )
    })
    @PostMapping(value = "/download-and-install", produces = MediaType.TEXT_EVENT_STREAM_VALUE)
    SseEmitter downloadUpdate();

    @Operation(
            summary = "Get current running application version",
            description = "Returns the local application version (`app.version`) instantly from DevaultyProperties without querying external services."
    )
    @ApiResponses(value = {
            @ApiResponse(
                    responseCode = "200",
                    description = "Returns current application version string.",
                    content = @Content(mediaType = MediaType.APPLICATION_JSON_VALUE, schema = @Schema(implementation = CurrentVersionResponse.class))
            ),
            @ApiResponse(
                    responseCode = "403",
                    description = "Forbidden. Request missing or containing an invalid `X-Devaulty-Internal-Token` header.",
                    content = @Content(mediaType = MediaType.APPLICATION_JSON_VALUE, schema = @Schema(implementation = ApiErrorResponse.class))
            )
    })
    @GetMapping("/current-app-version")
    ResponseEntity<CurrentVersionResponse> getAppInfo();
}