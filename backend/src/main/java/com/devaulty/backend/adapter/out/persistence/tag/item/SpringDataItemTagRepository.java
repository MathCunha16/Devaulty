package com.devaulty.backend.adapter.out.persistence.tag.item;

import org.springframework.data.jpa.repository.JpaRepository;
import org.springframework.transaction.annotation.Transactional;

import java.util.UUID;

public interface SpringDataItemTagRepository extends JpaRepository<ItemTagEntity, ItemTagId> {

    @Transactional
    void deleteByItemTypeAndItemId(String itemType, UUID itemId);
}
