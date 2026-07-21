import { useQuery, useMutation, useQueryClient } from "@tanstack/react-query";
import { notesApi } from "../api/notesApi";
import type { CreateNoteRequest, UpdateNoteRequest } from "~types/api";

export const notesKeys = {
  all: (projectId: string) => ["projects", projectId, "notes"] as const,
  detail: (projectId: string, noteId: string) =>
    ["projects", projectId, "notes", "detail", noteId] as const,
};

export const useNotesQuery = (projectId: string, page = 0, size = 100) => {
  return useQuery({
    queryKey: [...notesKeys.all(projectId), { page, size }],
    queryFn: () => notesApi.getAllByProject(projectId, page, size),
  });
};

export const useNoteQuery = (projectId: string, noteId: string) => {
  return useQuery({
    queryKey: notesKeys.detail(projectId, noteId),
    queryFn: () => notesApi.getById(projectId, noteId),
    enabled: !!noteId,
  });
};

export const useCreateNoteMutation = (projectId: string) => {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: (request: CreateNoteRequest) => notesApi.create(projectId, request),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: notesKeys.all(projectId) });
    },
  });
};

export const useUpdateNoteMutation = (projectId: string, noteId: string) => {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: (request: UpdateNoteRequest) => notesApi.update(projectId, noteId, request),
    onSuccess: (data) => {
      queryClient.invalidateQueries({ queryKey: notesKeys.all(projectId) });
      queryClient.invalidateQueries({ queryKey: notesKeys.detail(projectId, noteId) });
      queryClient.setQueryData(notesKeys.detail(projectId, noteId), data);
    },
  });
};

export const useDeleteNoteMutation = (projectId: string) => {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: (noteId: string) => notesApi.delete(projectId, noteId),
    onSuccess: (_, noteId) => {
      queryClient.invalidateQueries({ queryKey: notesKeys.all(projectId) });
      queryClient.invalidateQueries({ queryKey: notesKeys.detail(projectId, noteId) });
    },
  });
};

export const useArchiveNoteMutation = (projectId: string) => {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: (noteId: string) => notesApi.archive(projectId, noteId),
    onSuccess: (_, noteId) => {
      queryClient.invalidateQueries({ queryKey: notesKeys.all(projectId) });
      queryClient.invalidateQueries({ queryKey: notesKeys.detail(projectId, noteId) });
    },
  });
};

export const useUnarchiveNoteMutation = (projectId: string) => {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: (noteId: string) => notesApi.unarchive(projectId, noteId),
    onSuccess: (_, noteId) => {
      queryClient.invalidateQueries({ queryKey: notesKeys.all(projectId) });
      queryClient.invalidateQueries({ queryKey: notesKeys.detail(projectId, noteId) });
    },
  });
};
