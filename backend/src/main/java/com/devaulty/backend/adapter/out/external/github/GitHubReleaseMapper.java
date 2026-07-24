package com.devaulty.backend.adapter.out.external.github;

import com.devaulty.backend.adapter.out.external.github.dto.GitHubAssetResponse;
import com.devaulty.backend.adapter.out.external.github.dto.GitHubReleaseResponse;
import com.devaulty.backend.application.port.out.external.release.dto.LatestReleaseInfo;
import com.devaulty.backend.application.port.out.external.release.dto.ReleaseAssetInfo;
import org.mapstruct.Mapper;
import org.mapstruct.Mapping;

@Mapper(componentModel = "spring")
public interface GitHubReleaseMapper {

    @Mapping(target = "isPreRelease", source = "preRelease")
    @Mapping(target = "name", source = "name", defaultValue = "")
    @Mapping(target = "body", source = "body", defaultValue = "Version notes not informed")
    LatestReleaseInfo toDomain(GitHubReleaseResponse dto);

    @Mapping(target = "fileName", source = "name")
    @Mapping(target = "downloadUrl", source = "browserDownloadUrl")
    @Mapping(target = "sizeInBytes", source = "size")
    ReleaseAssetInfo toAssetDomain(GitHubAssetResponse dto);

}
