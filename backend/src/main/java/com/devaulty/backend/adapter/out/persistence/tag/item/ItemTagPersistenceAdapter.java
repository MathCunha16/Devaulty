package com.devaulty.backend.adapter.out.persistence.tag.item;

import com.devaulty.backend.adapter.out.persistence.tag.SpringDataTagRepository;
import com.devaulty.backend.adapter.out.persistence.tag.TagEntity;
import com.devaulty.backend.adapter.out.persistence.tag.TagMapper;
import com.devaulty.backend.application.port.out.persistence.ItemTagRepositoryPort;
import com.devaulty.backend.domain.model.Tag;
import org.springframework.stereotype.Component;

import java.util.List;
import java.util.Map;
import java.util.UUID;
import java.util.stream.Collectors;

@Component
public class ItemTagPersistenceAdapter implements ItemTagRepositoryPort {

    private final SpringDataItemTagRepository itemTagRepository;
    private final SpringDataTagRepository tagRepository;
    private final TagMapper tagMapper;

    public ItemTagPersistenceAdapter(SpringDataItemTagRepository itemTagRepository, SpringDataTagRepository tagRepository, TagMapper tagMapper) {
        this.itemTagRepository = itemTagRepository;
        this.tagRepository = tagRepository;
        this.tagMapper = tagMapper;
    }

    @Override
    public void associateTagToItem(UUID tagId, String itemType, UUID itemId) {
        TagEntity tagReference = tagRepository.getReferenceById(tagId);
        itemTagRepository.save(new ItemTagEntity(tagReference, itemType, itemId));
    }

    @Override
    public void disassembleTagFromItem(UUID projectId, UUID tagId, String itemType, UUID itemId) {
        itemTagRepository.deleteByTagIdAndItemTypeAndItemIdAndProjectId(tagId, itemType, itemId, projectId);
    }

    @Override
    public void removeAllTagsFromItem(String itemType, UUID itemId) {
        itemTagRepository.deleteByItemTypeAndItemId(itemType, itemId);
    }

    @Override
    public List<Tag> findTagsForItem(String itemType, UUID projectId, UUID itemId) {
        List<ItemTagEntity> entities = itemTagRepository.findByItemTypeAndItemIdAndProjectId(itemType, itemId, projectId);
        return entities.stream()
                .map(itemTagEntity -> tagMapper.toDomain(itemTagEntity.getTag()))
                .toList();
    }

    @Override
    public Map<UUID, List<Tag>> findTagsForItems(String itemType, UUID projectId, List<UUID> itemIds) {
        if (itemIds == null || itemIds.isEmpty()) {
            return Map.of();
        }
        return itemTagRepository.findByItemTypeAndItemIdInAndProjectId(itemType, itemIds, projectId)
                .stream()
                .collect(Collectors.groupingBy(
                        ItemTagEntity::getItemId,
                        Collectors.mapping(e -> tagMapper.toDomain(e.getTag()), Collectors.toList())
                ));
    }
}
