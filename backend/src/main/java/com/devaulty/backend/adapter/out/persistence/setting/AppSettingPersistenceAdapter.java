package com.devaulty.backend.adapter.out.persistence.setting;

import com.devaulty.backend.application.port.out.persistence.AppSettingRepositoryPort;
import com.devaulty.backend.domain.model.AppSetting;
import org.springframework.stereotype.Component;

import java.util.Optional;

@Component
public class AppSettingPersistenceAdapter implements AppSettingRepositoryPort {

    private final SpringDataAppSettingRepository appSettingRepository;
    private final AppSettingMapper appSettingMapper;

    public AppSettingPersistenceAdapter(SpringDataAppSettingRepository appSettingRepository, AppSettingMapper appSettingMapper) {
        this.appSettingRepository = appSettingRepository;
        this.appSettingMapper = appSettingMapper;
    }

    @Override
    public AppSetting save(AppSetting appSetting) {
        AppSettingEntity appSettingEntity = appSettingMapper.toEntity(appSetting);
        return appSettingMapper.toDomain(appSettingRepository.save(appSettingEntity));
    }

    @Override
    public Optional<AppSetting> findByKey(String key) {
        Optional<AppSettingEntity> entity = appSettingRepository.findById(key);
        return entity.map(appSettingMapper::toDomain);
    }

    @Override
    public boolean existsByKey(String key) {
        return appSettingRepository.existsById(key);
    }
}
