import { useQuery, useMutation, useQueryClient } from "@tanstack/react-query";
import { tagsApi } from "../api/tagsApi";
import type { CreateTagRequest, UpdateTagRequest } from "~types/api";

export const tagsKeys = {
  all: (projectId: string) => ["projects", projectId, "tags"] as const,
  search: (projectId: string, name: string) =>
    ["projects", projectId, "tags", "search", name] as const,
};

export const useTagsQuery = (projectId: string) => {
  return useQuery({
    queryKey: tagsKeys.all(projectId),
    queryFn: () => tagsApi.getAllByProject(projectId),
  });
};

export const useSearchTagsQuery = (projectId: string, name: string, enabled = true) => {
  return useQuery({
    queryKey: tagsKeys.search(projectId, name),
    queryFn: () => tagsApi.search(projectId, name),
    enabled: enabled && !!name.trim(),
  });
};

export const useCreateTagMutation = (projectId: string) => {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: (request: CreateTagRequest) => tagsApi.create(projectId, request),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: tagsKeys.all(projectId) });
    },
  });
};

export const useUpdateTagMutation = (projectId: string) => {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: ({ tagId, request }: { tagId: string; request: UpdateTagRequest }) =>
      tagsApi.update(projectId, tagId, request),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: tagsKeys.all(projectId) });
      // Invalidate related entities
      queryClient.invalidateQueries({ queryKey: ["projects", projectId, "problems"] });
      queryClient.invalidateQueries({ queryKey: ["projects", projectId, "snippets"] });
      queryClient.invalidateQueries({ queryKey: ["projects", projectId, "notes"] });
      queryClient.invalidateQueries({ queryKey: ["projects", projectId, "links"] });
      queryClient.invalidateQueries({ queryKey: ["projects", projectId, "credentials"] });
    },
  });
};

export const useDeleteTagMutation = (projectId: string) => {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: (tagId: string) => tagsApi.delete(projectId, tagId),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: tagsKeys.all(projectId) });
      queryClient.invalidateQueries({ queryKey: ["projects", projectId, "problems"] });
      queryClient.invalidateQueries({ queryKey: ["projects", projectId, "snippets"] });
      queryClient.invalidateQueries({ queryKey: ["projects", projectId, "notes"] });
      queryClient.invalidateQueries({ queryKey: ["projects", projectId, "links"] });
      queryClient.invalidateQueries({ queryKey: ["projects", projectId, "credentials"] });
    },
  });
};

interface AssociationParams {
  itemType: string;
  itemId: string;
  tagId: string;
}

export const useAssociateTagMutation = (projectId: string) => {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: ({ itemType, itemId, tagId }: AssociationParams) =>
      tagsApi.associate(projectId, itemType, itemId, tagId),
    onSuccess: (_, variables) => {
      const typeLower = variables.itemType.toLowerCase();
      let itemKey = "snippets";
      if (typeLower === "problem") itemKey = "problems";
      else if (typeLower === "note") itemKey = "notes";
      else if (typeLower === "link") itemKey = "links";
      else if (typeLower === "credential") itemKey = "credentials";
      queryClient.invalidateQueries({ queryKey: ["projects", projectId, itemKey] });
    },
  });
};

export const useDisassociateTagMutation = (projectId: string) => {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: ({ itemType, itemId, tagId }: AssociationParams) =>
      tagsApi.disassociate(projectId, itemType, itemId, tagId),
    onSuccess: (_, variables) => {
      const typeLower = variables.itemType.toLowerCase();
      let itemKey = "snippets";
      if (typeLower === "problem") itemKey = "problems";
      else if (typeLower === "note") itemKey = "notes";
      else if (typeLower === "link") itemKey = "links";
      else if (typeLower === "credential") itemKey = "credentials";
      queryClient.invalidateQueries({ queryKey: ["projects", projectId, itemKey] });
    },
  });
};
