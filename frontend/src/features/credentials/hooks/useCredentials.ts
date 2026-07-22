import { useQuery, useMutation, useQueryClient } from "@tanstack/react-query";
import { credentialsApi } from "../api/credentialsApi";
import type { CreateCredentialRequest, UpdateCredentialRequest } from "~types/api";

export const credentialsKeys = {
  all: (projectId: string) => ["projects", projectId, "credentials"] as const,
  detail: (projectId: string, credentialId: string) =>
    ["projects", projectId, "credentials", credentialId] as const,
};

export const useCredentialsQuery = (projectId: string, enabled = true) => {
  return useQuery({
    queryKey: credentialsKeys.all(projectId),
    queryFn: () => credentialsApi.getAllByProject(projectId),
    enabled: enabled && !!projectId,
    retry: false, // Do not retry if vault is locked (423) or master password missing (403)
  });
};

export const useCredentialQuery = (projectId: string, credentialId: string, enabled = true) => {
  return useQuery({
    queryKey: credentialsKeys.detail(projectId, credentialId),
    queryFn: () => credentialsApi.getById(projectId, credentialId),
    enabled: enabled && !!projectId && !!credentialId,
    retry: false,
    gcTime: 0, // Do not retain decrypted credentials in cache memory when unmounted
  });
};

export const useCreateCredentialMutation = (projectId: string) => {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: (request: CreateCredentialRequest) =>
      credentialsApi.create(projectId, request),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: credentialsKeys.all(projectId) });
    },
  });
};

export const useUpdateCredentialMutation = (projectId: string, credentialId: string) => {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: (request: UpdateCredentialRequest) =>
      credentialsApi.update(projectId, credentialId, request),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: credentialsKeys.all(projectId) });
      queryClient.invalidateQueries({ queryKey: credentialsKeys.detail(projectId, credentialId) });
    },
  });
};

export const useDeleteCredentialMutation = (projectId: string) => {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: (credentialId: string) => credentialsApi.delete(projectId, credentialId),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: credentialsKeys.all(projectId) });
    },
  });
};
