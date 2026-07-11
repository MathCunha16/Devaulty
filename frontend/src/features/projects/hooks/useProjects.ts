import { useSuspenseQuery, useMutation, useQueryClient } from "@tanstack/react-query";
import { projectsApi } from "../api/projectsApi";
import type { CreateProjectRequest, UpdateProjectRequest } from "~types/api";

export const projectsKeys = {
  all: ["projects"] as const,
  detail: (id: string) => ["projects", id] as const,
};

export const useProjectsQuery = (page = 0, size = 100) => {
  return useSuspenseQuery({
    queryKey: [...projectsKeys.all, { page, size }],
    queryFn: () => projectsApi.getAll(page, size),
  });
};

export const useProjectQuery = (id: string) => {
  return useSuspenseQuery({
    queryKey: projectsKeys.detail(id),
    queryFn: () => projectsApi.getById(id),
  });
};

export const useCreateProjectMutation = () => {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: (request: CreateProjectRequest) => projectsApi.create(request),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: projectsKeys.all });
    },
  });
};

export const useUpdateProjectMutation = (id: string) => {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: (request: UpdateProjectRequest) => projectsApi.update(id, request),
    onSuccess: (data) => {
      queryClient.invalidateQueries({ queryKey: projectsKeys.all });
      queryClient.setQueryData(projectsKeys.detail(id), data);
    },
  });
};

export const useArchiveProjectMutation = () => {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: (id: string) => projectsApi.archive(id),
    onSuccess: (_, id) => {
      queryClient.invalidateQueries({ queryKey: projectsKeys.all });
      queryClient.invalidateQueries({ queryKey: projectsKeys.detail(id) });
    },
  });
};

export const useUnarchiveProjectMutation = () => {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: (id: string) => projectsApi.unarchive(id),
    onSuccess: (_, id) => {
      queryClient.invalidateQueries({ queryKey: projectsKeys.all });
      queryClient.invalidateQueries({ queryKey: projectsKeys.detail(id) });
    },
  });
};

export const useDeleteProjectMutation = () => {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: (id: string) => projectsApi.delete(id),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: projectsKeys.all });
    },
  });
};
