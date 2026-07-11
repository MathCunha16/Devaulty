package com.devaulty.backend.application.port.in.link;

import com.devaulty.backend.domain.model.Link;
import org.springframework.data.domain.Page;

import java.util.UUID;

public interface GetAllLinksByProjectUseCase {
    Page<Link> execute(UUID projectId, int page, int size);
}
