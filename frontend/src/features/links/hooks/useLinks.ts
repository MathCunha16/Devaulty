import { useQuery, useMutation, useQueryClient } from "@tanstack/react-query";
import { linksApi } from "../api/linksApi";
import type { CreateLinkRequest, UpdateLinkRequest } from "~types/api";

export const linksKeys = {
  all: (projectId: string) => ["projects", projectId, "links"] as const,
  detail: (projectId: string, linkId: string) =>
    ["projects", projectId, "links", "detail", linkId] as const,
};

export const useLinksQuery = (projectId: string, page = 0, size = 100) => {
  return useQuery({
    queryKey: [...linksKeys.all(projectId), { page, size }],
    queryFn: () => linksApi.getAllByProject(projectId, page, size),
  });
};

export const useLinkQuery = (projectId: string, linkId: string) => {
  return useQuery({
    queryKey: linksKeys.detail(projectId, linkId),
    queryFn: () => linksApi.getById(projectId, linkId),
    enabled: !!linkId,
  });
};

export const useCreateLinkMutation = (projectId: string) => {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: (request: CreateLinkRequest) => linksApi.create(projectId, request),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: linksKeys.all(projectId) });
    },
  });
};

export const useUpdateLinkMutation = (projectId: string, linkId: string) => {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: (request: UpdateLinkRequest) => linksApi.update(projectId, linkId, request),
    onSuccess: (data) => {
      queryClient.invalidateQueries({ queryKey: linksKeys.all(projectId) });
      queryClient.invalidateQueries({ queryKey: linksKeys.detail(projectId, linkId) });
      queryClient.setQueryData(linksKeys.detail(projectId, linkId), data);
    },
  });
};

export const useDeleteLinkMutation = (projectId: string) => {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: (linkId: string) => linksApi.delete(projectId, linkId),
    onSuccess: (_, linkId) => {
      queryClient.invalidateQueries({ queryKey: linksKeys.all(projectId) });
      queryClient.invalidateQueries({ queryKey: linksKeys.detail(projectId, linkId) });
    },
  });
};
