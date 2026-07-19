package com.devaulty.backend.adapter.out.persistence.tag.item;

import org.hibernate.annotations.JdbcTypeCode;
import org.hibernate.type.SqlTypes;

import java.io.Serializable;
import java.util.Objects;
import java.util.UUID;

public class ItemTagId implements Serializable {

    @JdbcTypeCode(SqlTypes.VARCHAR)
    private UUID tag;

    private String itemType;

    @JdbcTypeCode(SqlTypes.VARCHAR)
    private UUID itemId;

    public ItemTagId() {
        // Empty constructor
    }

    public ItemTagId(UUID tag, String itemType, UUID itemId) {
        this.tag = tag;
        this.itemType = itemType;
        this.itemId = itemId;
    }

    @Override
    public boolean equals(Object o) {
        if (this == o) return true;
        if (o == null || getClass() != o.getClass()) return false;
        ItemTagId that = (ItemTagId) o;
        return Objects.equals(tag, that.tag)
                && Objects.equals(itemType, that.itemType)
                && Objects.equals(itemId, that.itemId);
    }

    @Override
    public int hashCode() {
        return Objects.hash(tag, itemType, itemId);
    }
}