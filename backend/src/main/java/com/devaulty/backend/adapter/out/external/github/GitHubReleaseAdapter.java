package com.devaulty.backend.adapter.out.external.github;

import com.devaulty.backend.application.port.out.external.release.ReleasePort;
import com.devaulty.backend.application.port.out.external.release.dto.LatestReleaseInfo;
import org.springframework.stereotype.Component;

@Component
public class GitHubReleaseAdapter implements ReleasePort {

    private final GitHubReleaseClient gitHubReleaseClient;
    private final GitHubReleaseMapper mapper;

    public GitHubReleaseAdapter(GitHubReleaseClient gitHubReleaseClient, GitHubReleaseMapper mapper) {
        this.gitHubReleaseClient = gitHubReleaseClient;
        this.mapper = mapper;
    }

    @Override
    public LatestReleaseInfo getLatestRelease() {
        return mapper.toDomain(gitHubReleaseClient.getLatestRelease());
    }
}
