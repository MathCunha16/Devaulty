import { useQuery, useMutation, useQueryClient } from "@tanstack/react-query";
import { drawingsApi } from "../api/drawingsApi";
import type {
  CreateDrawingRequest,
  Drawing,
  UpdateDrawingRequest,
} from "~types/api";

export const drawingsKeys = {
  all: (projectId: string) => ["projects", projectId, "drawings"] as const,
  detail: (projectId: string, drawingId: string) =>
    ["projects", projectId, "drawings", "detail", drawingId] as const,
};

export const useDrawingsQuery = (projectId: string, page = 0, size = 100) => {
  return useQuery({
    queryKey: [...drawingsKeys.all(projectId), { page, size }],
    queryFn: () => drawingsApi.getAllByProject(projectId, page, size),
    enabled: !!projectId,
  });
};

export const useDrawingQuery = (projectId: string, drawingId: string) => {
  return useQuery({
    queryKey: drawingsKeys.detail(projectId, drawingId),
    queryFn: () => drawingsApi.getById(projectId, drawingId),
    enabled: !!projectId && !!drawingId,
  });
};

export const useCreateDrawingMutation = (projectId: string) => {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: (request: CreateDrawingRequest) =>
      drawingsApi.create(projectId, request),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: drawingsKeys.all(projectId) });
    },
  });
};

export const useUpdateDrawingMutation = (projectId: string, drawingId: string) => {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: (request: UpdateDrawingRequest) =>
      drawingsApi.update(projectId, drawingId, request),
    onSuccess: (data: Drawing) => {
      queryClient.invalidateQueries({ queryKey: drawingsKeys.all(projectId) });
      queryClient.setQueryData(drawingsKeys.detail(projectId, drawingId), data);
    },
  });
};

export const useDeleteDrawingMutation = (projectId: string) => {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: (drawingId: string) => drawingsApi.delete(projectId, drawingId),
    onSuccess: (_, drawingId) => {
      queryClient.invalidateQueries({ queryKey: drawingsKeys.all(projectId) });
      queryClient.removeQueries({
        queryKey: drawingsKeys.detail(projectId, drawingId),
      });
    },
  });
};
