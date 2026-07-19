package com.devaulty.backend.adapter.out.persistence.tag.item;

import com.devaulty.backend.adapter.out.persistence.tag.TagEntity;
import jakarta.persistence.*;
import org.hibernate.annotations.JdbcTypeCode;
import org.hibernate.type.SqlTypes;

import java.util.Objects;
import java.util.UUID;

@Entity
@Table(name = "item_tags")
@IdClass(ItemTagId.class)
public class ItemTagEntity {

    @Id
    @ManyToOne(fetch = FetchType.LAZY)
    @JoinColumn(name = "tag_id")
    @JdbcTypeCode(SqlTypes.VARCHAR)
    private TagEntity tag;

    @Id
    @Column(name = "item_type")
    private String itemType;

    @Id
    @Column(name = "item_id")
    @JdbcTypeCode(SqlTypes.VARCHAR)
    private UUID itemId;

    public ItemTagEntity() {
        // Empty constructor
    }

    public ItemTagEntity(TagEntity tag, String itemType, UUID itemId) {
        this.tag = tag;
        this.itemType = itemType;
        this.itemId = itemId;
    }

    public TagEntity getTag() {
        return tag;
    }

    public void setTag(TagEntity tag) {
        this.tag = tag;
    }

    public String getItemType() {
        return itemType;
    }

    public void setItemType(String itemType) {
        this.itemType = itemType;
    }

    public UUID getItemId() {
        return itemId;
    }

    public void setItemId(UUID itemId) {
        this.itemId = itemId;
    }

    @Override
    public boolean equals(Object o) {
        if (o == null || getClass() != o.getClass()) return false;
        ItemTagEntity that = (ItemTagEntity) o;
        return Objects.equals(tag, that.tag)
                && Objects.equals(itemType, that.itemType)
                && Objects.equals(itemId, that.itemId);
    }

    @Override
    public int hashCode() {
        return Objects.hash(tag, itemType, itemId);
    }
}