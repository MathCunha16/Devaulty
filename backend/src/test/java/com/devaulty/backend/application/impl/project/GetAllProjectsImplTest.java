package com.devaulty.backend.application.impl.project;

import com.devaulty.backend.application.port.out.persistence.ProjectRepositoryPort;
import com.devaulty.backend.domain.model.Project;
import org.junit.jupiter.api.Test;
import org.junit.jupiter.api.extension.ExtendWith;
import org.mockito.InjectMocks;
import org.mockito.Mock;
import org.mockito.junit.jupiter.MockitoExtension;
import org.springframework.data.domain.Page;
import org.springframework.data.domain.PageImpl;

import java.util.Collections;
import java.util.List;
import java.util.UUID;

import static org.junit.jupiter.api.Assertions.*;
import static org.mockito.Mockito.*;

@ExtendWith(MockitoExtension.class)
class GetAllProjectsImplTest {

    @Mock
    private ProjectRepositoryPort projectRepositoryPort;

    @InjectMocks
    private GetAllProjectsImpl getAllProjectsUseCase;

    @Test
    void shouldReturnPageOfProjects() {
        // Arrange
        int page = 0;
        int size = 10;
        Project project = new Project(UUID.randomUUID(), "Test Project", "Desc", "icon", "#fff", false);
        List<Project> projectList = Collections.singletonList(project);
        Page<Project> expectedPage = new PageImpl<>(projectList);

        when(projectRepositoryPort.findAll(page, size)).thenReturn(expectedPage);

        // Act
        Page<Project> result = getAllProjectsUseCase.execute(page, size);

        // Assert
        assertNotNull(result);
        assertEquals(expectedPage, result);
        assertEquals(1, result.getTotalElements());
        assertEquals(project, result.getContent().getFirst());
        verify(projectRepositoryPort, times(1)).findAll(page, size);
    }
}
