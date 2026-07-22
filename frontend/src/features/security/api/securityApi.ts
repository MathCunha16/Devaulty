import { apiClient } from "@/api/client";
import type { MasterPasswordRequest, SessionStatus } from "~types/api";

export const securityApi = {
  /**
   * Returns true if master password setup is required (not configured yet), false if configured.
   */
  checkMasterPasswordSetup: async (): Promise<boolean> => {
    const response = await apiClient.get<boolean>("/security/master-password/required-status");
    return response.data;
  },

  /**
   * Configures the initial master password for the vault across all projects.
   */
  setupMasterPassword: async (request: MasterPasswordRequest): Promise<void> => {
    await apiClient.post("/security/master-password", request);
  },

  /**
   * Retrieves the current vault session status (active status and remaining seconds).
   */
  getSessionStatus: async (): Promise<SessionStatus> => {
    const response = await apiClient.get<SessionStatus>("/security/vault/status");
    return response.data;
  },

  /**
   * Unlocks the vault using master password.
   */
  unlockVault: async (request: MasterPasswordRequest): Promise<boolean> => {
    const response = await apiClient.post<boolean>("/security/vault/unlock", request);
    return response.data;
  },

  /**
   * Locks the vault, clearing cryptographic keys from RAM session.
   */
  lockVault: async (): Promise<void> => {
    await apiClient.post("/security/vault/lock");
  },
};
