package com.devaulty.backend.adapter.in.web.security;

import com.devaulty.backend.adapter.in.web.security.dto.MasterPasswordRequest;
import com.devaulty.backend.application.port.in.security.*;
import jakarta.validation.Valid;
import org.springframework.http.ResponseEntity;
import org.springframework.web.bind.annotation.*;

import java.util.Arrays;

@RestController
@RequestMapping("/api/v1/security")
public class SecurityController implements SecurityApi{

    private final SetupMasterPasswordUseCase setupMasterPasswordUseCase;
    private final CheckMasterPasswordSetupUseCase checkMasterPasswordSetupUseCase;
    private final LockVaultUseCase lockVaultUseCase;
    private final UnlockVaultUseCase unlockVaultUseCase;
    private final GetSessionStatusUseCase getSessionStatusUseCase;

    public SecurityController(SetupMasterPasswordUseCase setupMasterPasswordUseCase, CheckMasterPasswordSetupUseCase checkMasterPasswordSetupUseCase, LockVaultUseCase lockVaultUseCase, UnlockVaultUseCase unlockVaultUseCase, GetSessionStatusUseCase getSessionStatusUseCase) {
        this.setupMasterPasswordUseCase = setupMasterPasswordUseCase;
        this.checkMasterPasswordSetupUseCase = checkMasterPasswordSetupUseCase;
        this.lockVaultUseCase = lockVaultUseCase;
        this.unlockVaultUseCase = unlockVaultUseCase;
        this.getSessionStatusUseCase = getSessionStatusUseCase;
    }

    @Override
    @PostMapping("/master-password")
    public ResponseEntity<Void> setupMasterPassword(@RequestBody @Valid MasterPasswordRequest request) {
        setupMasterPasswordUseCase.execute(request.masterPassword());
        Arrays.fill(request.masterPassword(),'\0');
        return ResponseEntity.noContent().build();
    }

    @Override
    @GetMapping("/master-password/required-status")
    public ResponseEntity<Boolean> checkMasterPasswordSetup(){
        return ResponseEntity.ok(checkMasterPasswordSetupUseCase.isSetupRequired());
    }

    @Override
    @GetMapping("/vault/status")
    public ResponseEntity<SessionStatus> getSessionStatus(){
        return ResponseEntity.ok(getSessionStatusUseCase.execute());
    }

    @Override
    @PostMapping("/vault/unlock")
    public ResponseEntity<Boolean> unlockVault(@RequestBody @Valid MasterPasswordRequest request){
        boolean result = unlockVaultUseCase.execute(request.masterPassword());
        Arrays.fill(request.masterPassword(),'\0');
        return ResponseEntity.ok(result);
    }

    @Override
    @PostMapping("/vault/lock")
    public ResponseEntity<Void> lockVault(){
        lockVaultUseCase.execute();
        return ResponseEntity.noContent().build();
    }
}
