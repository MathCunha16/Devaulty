package com.devaulty.backend.adapter.out.external.github;

import com.devaulty.backend.application.port.out.external.release.ReleasePort;
import com.devaulty.backend.application.port.out.external.release.dto.LatestReleaseInfo;
import org.springframework.stereotype.Component;

import java.io.InputStream;

@Component
public class GitHubReleaseAdapter implements ReleasePort {

    private final GitHubClient gitHubClient;
    private final GitHubReleaseMapper mapper;

    public GitHubReleaseAdapter(GitHubClient gitHubClient, GitHubReleaseMapper mapper) {
        this.gitHubClient = gitHubClient;
        this.mapper = mapper;
    }

    @Override
    public LatestReleaseInfo getLatestRelease() {
        return mapper.toDomain(gitHubClient.getLatestRelease());
    }

    @Override
    public InputStream downloadAsset(String downloadUrl) {
        return gitHubClient.downloadAsset(downloadUrl);
    }
}