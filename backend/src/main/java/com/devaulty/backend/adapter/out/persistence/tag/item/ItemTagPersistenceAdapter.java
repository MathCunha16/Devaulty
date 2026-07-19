package com.devaulty.backend.adapter.out.persistence.tag.item;

import com.devaulty.backend.adapter.out.persistence.tag.SpringDataTagRepository;
import com.devaulty.backend.adapter.out.persistence.tag.TagEntity;
import com.devaulty.backend.application.port.out.persistence.ItemTagRepositoryPort;
import org.springframework.stereotype.Component;

import java.util.UUID;

@Component
public class ItemTagPersistenceAdapter implements ItemTagRepositoryPort {

    private final SpringDataItemTagRepository itemTagRepository;
    private final SpringDataTagRepository tagRepository;

    public ItemTagPersistenceAdapter(SpringDataItemTagRepository itemTagRepository, SpringDataTagRepository tagRepository) {
        this.itemTagRepository = itemTagRepository;
        this.tagRepository = tagRepository;
    }

    @Override
    public void associateTagToItem(UUID tagId, String itemType, UUID itemId) {
        TagEntity tagReference = tagRepository.getReferenceById(tagId);
        itemTagRepository.save(new ItemTagEntity(tagReference, itemType, itemId));
    }

    @Override
    public void disassembleTagFromItem(UUID tagId, String itemType, UUID itemId) {
        ItemTagId id = new ItemTagId(tagId, itemType, itemId);
        itemTagRepository.deleteById(id);
    }

    @Override
    public void removeAllTagsFromItem(String itemType, UUID itemId) {
        itemTagRepository.deleteByItemTypeAndItemId(itemType, itemId);
    }

}
