package com.devaulty.backend.application.port.in.release;

import reactor.core.publisher.Flux;

public interface DownloadUpdateUseCase {
    Flux<UpdateProgressInfo> execute();
}
