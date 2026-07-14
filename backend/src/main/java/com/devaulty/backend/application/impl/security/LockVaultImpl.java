package com.devaulty.backend.application.impl.security;

import com.devaulty.backend.application.port.out.security.MasterKeySessionPort;
import com.devaulty.backend.application.port.in.security.LockVaultUseCase;

public class LockVaultImpl implements LockVaultUseCase {

    private final MasterKeySessionPort sessionHolder;

    public LockVaultImpl(MasterKeySessionPort sessionHolder) {
        this.sessionHolder = sessionHolder;
    }

    @Override
    public void execute() {
        // Purges key bytes completely from RAM
        sessionHolder.clear();
    }
}
