package com.devaulty.backend.adapter.in.web.security;

import com.devaulty.backend.adapter.in.web.security.dto.MasterPasswordRequest;
import com.devaulty.backend.application.port.in.security.CheckMasterPasswordSetupUseCase;
import com.devaulty.backend.application.port.in.security.LockVaultUseCase;
import com.devaulty.backend.application.port.in.security.SetupMasterPasswordUseCase;
import com.devaulty.backend.application.port.in.security.UnlockVaultUseCase;
import jakarta.validation.Valid;
import org.springframework.http.ResponseEntity;
import org.springframework.web.bind.annotation.*;

import java.util.Arrays;

@RestController
@RequestMapping("/ap1/v1/security")
public class SecurityController {

    private final SetupMasterPasswordUseCase setupMasterPasswordUseCase;
    private final CheckMasterPasswordSetupUseCase checkMasterPasswordSetupUseCase;
    private final LockVaultUseCase lockVaultUseCase;
    private final UnlockVaultUseCase unlockVaultUseCase;

    public SecurityController(SetupMasterPasswordUseCase setupMasterPasswordUseCase, CheckMasterPasswordSetupUseCase checkMasterPasswordSetupUseCase, LockVaultUseCase lockVaultUseCase, UnlockVaultUseCase unlockVaultUseCase) {
        this.setupMasterPasswordUseCase = setupMasterPasswordUseCase;
        this.checkMasterPasswordSetupUseCase = checkMasterPasswordSetupUseCase;
        this.lockVaultUseCase = lockVaultUseCase;
        this.unlockVaultUseCase = unlockVaultUseCase;
    }

    @PostMapping("/master-password")
    public ResponseEntity<Void> setupMasterPassword(@RequestBody @Valid MasterPasswordRequest request) {
        setupMasterPasswordUseCase.execute(request.masterPassword());
        Arrays.fill(request.masterPassword(),'\0');
        return ResponseEntity.noContent().build();
    }

    @GetMapping("/master-password/required-status")
    public ResponseEntity<Boolean> checkMasterPasswordSetup(){
        return ResponseEntity.ok(checkMasterPasswordSetupUseCase.isSetupRequired());
    }

    @PostMapping("/vault/unlock")
    public ResponseEntity<Boolean> unlockVault(@RequestBody @Valid MasterPasswordRequest request){
        boolean result = unlockVaultUseCase.execute(request.masterPassword());
        Arrays.fill(request.masterPassword(),'\0');
        return ResponseEntity.ok(result);
    }

    @PostMapping("/vault/lock")
    public ResponseEntity<Void> lockVault(){
        lockVaultUseCase.execute();
        return ResponseEntity.noContent().build();
    }
}
