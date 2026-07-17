package com.devaulty.backend.application.port.out.persistence;

import com.devaulty.backend.domain.model.AppSetting;

import java.util.Optional;

public interface AppSettingRepositoryPort {

    AppSetting save(AppSetting appSetting);

    Optional<AppSetting> findByKey(String key);

    boolean existsByKey(String key);

}
