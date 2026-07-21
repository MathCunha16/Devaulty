import { apiClient } from "@/api/client";
import type {
  CreateProblemRequest,
  ProblemViewResponse,
  UpdateProblemRequest,
  UpdateProblemStatusRequest,
  PagedModelProblemSummaryResponse,
} from "~types/api";

export const problemsApi = {
  getAllByProject: async (
    projectId: string,
    page = 0,
    size = 100
  ): Promise<PagedModelProblemSummaryResponse> => {
    const response = await apiClient.get<PagedModelProblemSummaryResponse>(
      `/projects/${projectId}/problems`,
      {
        params: { page, size },
      }
    );
    return response.data;
  },

  getById: async (projectId: string, problemId: string): Promise<ProblemViewResponse> => {
    const response = await apiClient.get<ProblemViewResponse>(
      `/projects/${projectId}/problems/${problemId}`
    );
    return response.data;
  },

  create: async (projectId: string, request: CreateProblemRequest): Promise<ProblemViewResponse> => {
    const response = await apiClient.post<ProblemViewResponse>(
      `/projects/${projectId}/problems`,
      request
    );
    return response.data;
  },

  update: async (
    projectId: string,
    problemId: string,
    request: UpdateProblemRequest
  ): Promise<ProblemViewResponse> => {
    const response = await apiClient.patch<ProblemViewResponse>(
      `/projects/${projectId}/problems/${problemId}`,
      request
    );
    return response.data;
  },

  delete: async (projectId: string, problemId: string): Promise<void> => {
    await apiClient.delete(`/projects/${projectId}/problems/${problemId}`);
  },

  updateStatus: async (
    projectId: string,
    problemId: string,
    request: UpdateProblemStatusRequest
  ): Promise<ProblemViewResponse> => {
    const response = await apiClient.patch<ProblemViewResponse>(
      `/projects/${projectId}/problems/${problemId}/status`,
      request
    );
    return response.data;
  },
};
