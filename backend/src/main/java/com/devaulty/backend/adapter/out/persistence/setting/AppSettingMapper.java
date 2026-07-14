package com.devaulty.backend.adapter.out.persistence.setting;

import com.devaulty.backend.domain.model.AppSetting;
import org.mapstruct.Mapper;

@Mapper(componentModel = "spring")
public interface AppSettingMapper {

    AppSettingEntity toEntity(AppSetting appSetting);

    AppSetting toDomain(AppSettingEntity appSettingEntity);
}
