import { useQuery, useSuspenseQuery, useMutation, useQueryClient } from "@tanstack/react-query";
import { problemsApi } from "../api/problemsApi";
import type {
  CreateProblemRequest,
  UpdateProblemRequest,
  ProblemStatus,
} from "~types/api";

export const problemsKeys = {
  all: (projectId: string) => ["projects", projectId, "problems"] as const,
  detail: (projectId: string, problemId: string) =>
    ["projects", projectId, "problems", problemId] as const,
};

export const useProblemsQuery = (projectId: string, page = 0, size = 100) => {
  return useSuspenseQuery({
    queryKey: [...problemsKeys.all(projectId), { page, size }],
    queryFn: () => problemsApi.getAllByProject(projectId, page, size),
  });
};

export const useProblemQuery = (projectId: string, problemId: string) => {
  return useQuery({
    queryKey: problemsKeys.detail(projectId, problemId),
    queryFn: () => problemsApi.getById(projectId, problemId),
    enabled: !!problemId,
  });
};

export const useCreateProblemMutation = (projectId: string) => {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: (request: CreateProblemRequest) => problemsApi.create(projectId, request),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: problemsKeys.all(projectId) });
    },
  });
};

export const useUpdateProblemMutation = (projectId: string, problemId: string) => {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: (request: UpdateProblemRequest) =>
      problemsApi.update(projectId, problemId, request),
    onSuccess: (data) => {
      queryClient.invalidateQueries({ queryKey: problemsKeys.all(projectId) });
      queryClient.invalidateQueries({ queryKey: problemsKeys.detail(projectId, problemId) });
      queryClient.setQueryData(problemsKeys.detail(projectId, problemId), data);
    },
  });
};

export const useUpdateProblemStatusMutation = (projectId: string) => {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: ({ problemId, status }: { problemId: string; status: ProblemStatus }) =>
      problemsApi.updateStatus(projectId, problemId, { status }),
    onSuccess: (data, variables) => {
      queryClient.invalidateQueries({ queryKey: problemsKeys.all(projectId) });
      queryClient.invalidateQueries({ queryKey: problemsKeys.detail(projectId, variables.problemId) });
      queryClient.setQueryData(problemsKeys.detail(projectId, variables.problemId), data);
    },
  });
};

export const useDeleteProblemMutation = (projectId: string) => {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: (problemId: string) => problemsApi.delete(projectId, problemId),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: problemsKeys.all(projectId) });
    },
  });
};
