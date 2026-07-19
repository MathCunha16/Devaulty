package com.devaulty.backend.adapter.in.web.tag.item;

import com.devaulty.backend.application.port.in.tag.item.AssociateTagToItemUseCase;
import com.devaulty.backend.application.port.in.tag.item.RemoveTagFromItemUseCase;
import org.springframework.http.ResponseEntity;
import org.springframework.web.bind.annotation.*;

import java.util.UUID;

@RestController
@RequestMapping("/api/v1/projects/{projectId}/items/{itemType}/{itemId}/tags")
public class ItemTagController implements ItemTagApi{

    private final AssociateTagToItemUseCase associateTagToItemUseCase;
    private final RemoveTagFromItemUseCase removeTagFromItemUseCase;

    public ItemTagController(AssociateTagToItemUseCase associateTagToItemUseCase, RemoveTagFromItemUseCase removeTagFromItemUseCase) {
        this.associateTagToItemUseCase = associateTagToItemUseCase;
        this.removeTagFromItemUseCase = removeTagFromItemUseCase;
    }

    @Override
    @PutMapping("/{tagId}")
    public ResponseEntity<Void> associate(
            @PathVariable UUID projectId,
            @PathVariable String itemType,
            @PathVariable UUID itemId,
            @PathVariable UUID tagId
    ){
        associateTagToItemUseCase.execute(projectId, itemType, itemId, tagId);
        return ResponseEntity.noContent().build();
    }

    @Override
    @DeleteMapping("/{tagId}")
    public ResponseEntity<Void> disassociate(
            @PathVariable UUID projectId,
            @PathVariable String itemType,
            @PathVariable UUID itemId,
            @PathVariable UUID tagId
    ){
        removeTagFromItemUseCase.execute(projectId, itemType, itemId, tagId);
        return ResponseEntity.noContent().build();
    }
}
