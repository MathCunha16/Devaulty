import { useSuspenseQuery, useMutation, useQueryClient } from "@tanstack/react-query";
import { snippetsApi } from "../api/snippetsApi";
import type { CreateSnippetRequest, UpdateSnippetRequest } from "~types/api";

export const snippetsKeys = {
  all: (projectId: string) => ["projects", projectId, "snippets"] as const,
  detail: (projectId: string, snippetId: string) =>
    ["projects", projectId, "snippets", snippetId] as const,
};

export const useSnippetsQuery = (projectId: string, page = 0, size = 100) => {
  return useSuspenseQuery({
    queryKey: [...snippetsKeys.all(projectId), { page, size }],
    queryFn: () => snippetsApi.getAllByProject(projectId, page, size),
  });
};

export const useSnippetQuery = (projectId: string, snippetId: string) => {
  return useSuspenseQuery({
    queryKey: snippetsKeys.detail(projectId, snippetId),
    queryFn: () => snippetsApi.getById(projectId, snippetId),
  });
};

export const useCreateSnippetMutation = (projectId: string) => {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: (request: CreateSnippetRequest) => snippetsApi.create(projectId, request),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: snippetsKeys.all(projectId) });
    },
  });
};

export const useUpdateSnippetMutation = (projectId: string, snippetId: string) => {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: (request: UpdateSnippetRequest) =>
      snippetsApi.update(projectId, snippetId, request),
    onSuccess: (data) => {
      queryClient.invalidateQueries({ queryKey: snippetsKeys.all(projectId) });
      queryClient.setQueryData(snippetsKeys.detail(projectId, snippetId), data);
    },
  });
};

export const useDeleteSnippetMutation = (projectId: string) => {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: (snippetId: string) => snippetsApi.delete(projectId, snippetId),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: snippetsKeys.all(projectId) });
    },
  });
};
