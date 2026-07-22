import { useQuery, useMutation, useQueryClient } from "@tanstack/react-query";
import { securityApi } from "../api/securityApi";
import type { MasterPasswordRequest } from "~types/api";

export const securityKeys = {
  setupStatus: ["security", "master-password", "setup-status"] as const,
  vaultStatus: ["security", "vault", "status"] as const,
};

export const useMasterPasswordSetupStatusQuery = (enabled = true) => {
  return useQuery({
    queryKey: securityKeys.setupStatus,
    queryFn: securityApi.checkMasterPasswordSetup,
    enabled,
    staleTime: 1000 * 60 * 5, // 5 minutes
  });
};

export const useVaultStatusQuery = (enabled = true) => {
  return useQuery({
    queryKey: securityKeys.vaultStatus,
    queryFn: securityApi.getSessionStatus,
    enabled,
    refetchInterval: (query) => {
      const data = query.state.data;
      if (data?.active) {
        return 10000; // Refetch every 10 seconds while vault session is active
      }
      return false;
    },
  });
};

export const useSetupMasterPasswordMutation = () => {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: (request: MasterPasswordRequest) => securityApi.setupMasterPassword(request),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: securityKeys.setupStatus });
      queryClient.invalidateQueries({ queryKey: securityKeys.vaultStatus });
    },
  });
};

export const useUnlockVaultMutation = () => {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: (request: MasterPasswordRequest) => securityApi.unlockVault(request),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: securityKeys.vaultStatus });
    },
  });
};

export const useLockVaultMutation = () => {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: () => securityApi.lockVault(),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: securityKeys.vaultStatus });
      // Wipes cached credential queries from memory upon lock
      queryClient.removeQueries({
        predicate: (query) => query.queryKey.includes("credentials"),
      });
    },
  });
};
