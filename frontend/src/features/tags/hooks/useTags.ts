import { useQuery, useMutation, useQueryClient } from "@tanstack/react-query";
import { tagsApi } from "../api/tagsApi";
import type { CreateTagRequest } from "~types/api";

export const tagsKeys = {
  all: (projectId: string) => ["projects", projectId, "tags"] as const,
};

export const useTagsQuery = (projectId: string) => {
  return useQuery({
    queryKey: tagsKeys.all(projectId),
    queryFn: () => tagsApi.getAllByProject(projectId),
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

export const useDeleteTagMutation = (projectId: string) => {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: (tagId: string) => tagsApi.delete(projectId, tagId),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: tagsKeys.all(projectId) });
      // Invalidate all items that might be using this tag
      queryClient.invalidateQueries({ queryKey: ["projects", projectId, "problems"] });
      queryClient.invalidateQueries({ queryKey: ["projects", projectId, "snippets"] });
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
      // Invalidate specific item and the lists
      const itemKey = variables.itemType.toLowerCase() === "problem" ? "problems" : "snippets";
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
      const itemKey = variables.itemType.toLowerCase() === "problem" ? "problems" : "snippets";
      queryClient.invalidateQueries({ queryKey: ["projects", projectId, itemKey] });
    },
  });
};
