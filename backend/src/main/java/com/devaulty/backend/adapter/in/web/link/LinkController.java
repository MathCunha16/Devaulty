package com.devaulty.backend.adapter.in.web.link;

import com.devaulty.backend.adapter.in.web.link.dto.CreateLinkRequest;
import com.devaulty.backend.adapter.in.web.link.dto.LinkViewResponse;
import com.devaulty.backend.adapter.in.web.link.dto.UpdateLinkRequest;
import com.devaulty.backend.adapter.in.web.util.UriLocationBuilderHelper;
import com.devaulty.backend.application.port.in.link.*;
import com.devaulty.backend.domain.model.Link;
import com.devaulty.backend.application.port.in.tag.item.GetTagsForItemUseCase;
import com.devaulty.backend.application.port.in.tag.item.GetTagsForItemsUseCase;
import com.devaulty.backend.domain.model.Tag;
import jakarta.validation.Valid;
import org.springframework.data.domain.Page;
import org.springframework.http.ResponseEntity;
import org.springframework.web.bind.annotation.*;

import java.net.URI;
import java.util.List;
import java.util.Map;
import java.util.UUID;

@RestController
@RequestMapping("/api/v1/projects/{projectId}/links")
public class LinkController implements LinkApi{

    private final CreateLinkUseCase createLinkUseCase;
    private final GetAllLinksByProjectUseCase getAllLinksByProjectUseCase;
    private final GetLinkByIdUseCase getLinkByIdUseCase;
    private final DeleteLinkUseCase deleteLinkUseCase;
    private final UpdateLinkUseCase updateLinkUseCase;
    private final GetTagsForItemUseCase getTagsForItemUseCase;
    private final GetTagsForItemsUseCase getTagsForItemsUseCase;
    private final LinkWebMapper webMapper;
    private final UriLocationBuilderHelper uriLocationBuilderHelper;

    private static final String ITEM_TYPE = "link";

    public LinkController(CreateLinkUseCase createLinkUseCase, GetAllLinksByProjectUseCase getAllLinksByProjectUseCase, GetLinkByIdUseCase getLinkByIdUseCase, DeleteLinkUseCase deleteLinkUseCase, UpdateLinkUseCase updateLinkUseCase, GetTagsForItemUseCase getTagsForItemUseCase, GetTagsForItemsUseCase getTagsForItemsUseCase, LinkWebMapper webMapper, UriLocationBuilderHelper uriLocationBuilderHelper) {
        this.createLinkUseCase = createLinkUseCase;
        this.getAllLinksByProjectUseCase = getAllLinksByProjectUseCase;
        this.getLinkByIdUseCase = getLinkByIdUseCase;
        this.deleteLinkUseCase = deleteLinkUseCase;
        this.updateLinkUseCase = updateLinkUseCase;
        this.getTagsForItemUseCase = getTagsForItemUseCase;
        this.getTagsForItemsUseCase = getTagsForItemsUseCase;
        this.webMapper = webMapper;
        this.uriLocationBuilderHelper = uriLocationBuilderHelper;
    }

    @Override
    @PostMapping
    public ResponseEntity<LinkViewResponse> create(
            @PathVariable UUID projectId,
            @RequestBody @Valid CreateLinkRequest request
    ){
        CreateLinkCommand command = webMapper.toCreateLinkCommand(request, projectId);
        Link link = createLinkUseCase.execute(command);
        URI uri = uriLocationBuilderHelper.buildLocationUri(link.getId());
        return ResponseEntity.created(uri).body(webMapper.toViewResponse(link, List.of()));
    }

    @Override
    @GetMapping
    public ResponseEntity<Page<LinkViewResponse>> getAllByProject(
            @PathVariable UUID projectId,
            @RequestParam(defaultValue = "0") int page,
            @RequestParam(defaultValue = "10") int size
    ){
        Page<Link> links = getAllLinksByProjectUseCase.execute(projectId, page, size);
        List<UUID> ids = links.getContent().stream().map(Link::getId).toList();
        Map<UUID, List<Tag>> tagsByItem = getTagsForItemsUseCase.execute(ITEM_TYPE, projectId, ids);
        Page<LinkViewResponse> response = links.map(s ->
                webMapper.toViewResponse(s, tagsByItem.getOrDefault(s.getId(), List.of())));
        return ResponseEntity.ok(response);
    }

    @Override
    @GetMapping("/{linkId}")
    public ResponseEntity<LinkViewResponse> getById(
            @PathVariable UUID projectId,
            @PathVariable UUID linkId
    ){
        Link link = getLinkByIdUseCase.execute(projectId, linkId);
        List<Tag> tags = getTagsForItemUseCase.execute(ITEM_TYPE, projectId, linkId);
        return ResponseEntity.ok(webMapper.toViewResponse(link, tags));
    }

    @Override
    @PatchMapping("/{linkId}")
    public ResponseEntity<LinkViewResponse> update(
            @PathVariable UUID projectId,
            @PathVariable UUID linkId,
            @RequestBody @Valid UpdateLinkRequest request
    ){
        UpdateLinkCommand command = webMapper.toUpdateLinkCommand(request, projectId, linkId);
        Link link = updateLinkUseCase.execute(command);
        List<Tag> tags = getTagsForItemUseCase.execute(ITEM_TYPE, projectId, linkId);
        return ResponseEntity.ok(webMapper.toViewResponse(link, tags));
    }

    @Override
    @DeleteMapping("/{linkId}")
    public ResponseEntity<Void> delete(
            @PathVariable UUID projectId,
            @PathVariable UUID linkId
    ){
        deleteLinkUseCase.execute(projectId, linkId);
        return ResponseEntity.noContent().build();
    }
}
