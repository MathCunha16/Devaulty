package com.devaulty.backend.infrastructure.configuration;

import com.devaulty.backend.application.impl.note.*;
import com.devaulty.backend.application.port.in.note.*;
import com.devaulty.backend.application.port.out.persistence.NoteRepositoryPort;
import com.devaulty.backend.application.port.out.persistence.ProjectRepositoryPort;
import org.springframework.context.annotation.Bean;
import org.springframework.context.annotation.Configuration;

@Configuration
public class NoteBeanConfig {

    @Bean
    public CreateNoteUseCase createNoteUseCase(
            NoteRepositoryPort noteRepositoryPort,
            ProjectRepositoryPort projectRepositoryPort
    ){
        return new CreateNoteImpl(
                noteRepositoryPort,
                projectRepositoryPort
        );
    }

    @Bean
    public GetAllNotesByProjectUseCase getAllNotesByProjectUseCase(
            NoteRepositoryPort noteRepositoryPort,
            ProjectRepositoryPort projectRepositoryPort
    ){
        return new GetAllNotesByProjectImpl(
                noteRepositoryPort,
                projectRepositoryPort
        );
    }

    @Bean
    public GetNoteByIdUseCase getNoteByIdUseCase(
            NoteRepositoryPort noteRepositoryPort,
            ProjectRepositoryPort projectRepositoryPort
    ){
        return new GetNoteByIdImpl(
                noteRepositoryPort,
                projectRepositoryPort
        );
    }

    @Bean
    public DeleteNoteUseCase deleteNoteUseCase(
            NoteRepositoryPort noteRepositoryPort,
            ProjectRepositoryPort projectRepositoryPort
    ){
        return new DeleteNoteImpl(
                noteRepositoryPort,
                projectRepositoryPort
        );
    }

    @Bean
    public ArchiveNoteUseCase archiveNoteUseCase(
            NoteRepositoryPort noteRepositoryPort,
            ProjectRepositoryPort projectRepositoryPort
    ){
        return new ArchiveNoteImpl(
                noteRepositoryPort,
                projectRepositoryPort
        );
    }

    @Bean
    public UnarchiveNoteUseCase unarchiveNoteUseCase(
            NoteRepositoryPort noteRepositoryPort,
            ProjectRepositoryPort projectRepositoryPort
    ){
        return new UnarchiveNoteImpl(
                noteRepositoryPort,
                projectRepositoryPort
        );
    }

    @Bean
    public UpdateNoteUseCase updateNoteUseCase(
            NoteRepositoryPort noteRepositoryPort,
            ProjectRepositoryPort projectRepositoryPort
    ){
        return new UpdateNoteImpl(
                noteRepositoryPort,
                projectRepositoryPort
        );
    }
}
