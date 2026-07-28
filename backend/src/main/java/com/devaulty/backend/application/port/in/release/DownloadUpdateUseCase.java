package com.devaulty.backend.application.port.in.release;

import java.util.function.Consumer;

public interface DownloadUpdateUseCase {
    void execute(Consumer<UpdateProgressInfo> onProgress);
}
