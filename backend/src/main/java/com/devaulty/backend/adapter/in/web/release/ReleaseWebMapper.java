package com.devaulty.backend.adapter.in.web.release;

import com.devaulty.backend.adapter.in.web.release.dto.AppUpdateInfoResponse;
import com.devaulty.backend.application.port.in.release.AppUpdateInfo;
import org.mapstruct.Mapper;

import java.time.Instant;
import java.time.LocalDateTime;
import java.time.ZoneId;

@Mapper(componentModel = "spring")
public interface ReleaseWebMapper {

    AppUpdateInfoResponse toAppUpdateInfoResponse(AppUpdateInfo appUpdateInfo);

    default LocalDateTime toLocalDateTime(Instant instant) {
        if(instant == null) return null;
        return LocalDateTime.ofInstant(instant, ZoneId.systemDefault());
    }
}
