package com.devaulty.backend.application.port.in.release;

import java.nio.file.Path;

public interface InstallUpdateUseCase {
    void execute(Path installerFile);
}
