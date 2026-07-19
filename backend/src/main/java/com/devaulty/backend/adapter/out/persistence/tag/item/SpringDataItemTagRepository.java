package com.devaulty.backend.adapter.out.persistence.tag.item;

import org.springframework.data.jpa.repository.JpaRepository;
import org.springframework.data.jpa.repository.Modifying;
import org.springframework.data.jpa.repository.Query;
import org.springframework.data.repository.query.Param;
import org.springframework.transaction.annotation.Transactional;

import java.util.List;
import java.util.UUID;

public interface SpringDataItemTagRepository extends JpaRepository<ItemTagEntity, ItemTagId> {

    @Transactional
    void deleteByItemTypeAndItemId(String itemType, UUID itemId);

    @Query("SELECT it FROM ItemTagEntity it WHERE it.itemType = :itemType AND it.itemId = :itemId AND it.tag.project.id = :projectId")
    List<ItemTagEntity> findByItemTypeAndItemIdAndProjectId(
            @Param("itemType") String itemType,
            @Param("itemId") UUID itemId,
            @Param("projectId") UUID projectId
    );

    @Query("SELECT it FROM ItemTagEntity it WHERE it.itemType = :itemType AND it.itemId IN :itemIds AND it.tag.project.id = :projectId")
    List<ItemTagEntity> findByItemTypeAndItemIdInAndProjectId(
            @Param("itemType") String itemType,
            @Param("itemIds") List<UUID> itemIds,
            @Param("projectId") UUID projectId
    );

    @Modifying
    @Transactional
    @Query("DELETE FROM ItemTagEntity it WHERE it.tag.id = :tagId AND it.itemType = :itemType AND it.itemId = :itemId AND it.tag.project.id = :projectId")
    void deleteByTagIdAndItemTypeAndItemIdAndProjectId(
            @Param("tagId") UUID tagId,
            @Param("itemType") String itemType,
            @Param("itemId") UUID itemId,
            @Param("projectId") UUID projectId
    );
}
