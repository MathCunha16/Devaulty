import { apiClient } from "@/api/client";
import type {
  CreateProjectRequest,
  ProjectViewResponse,
  UpdateProjectRequest,
  PagedModelProjectViewResponse,
} from "~types/api";

export const projectsApi = {
  getAll: async (page = 0, size = 100): Promise<PagedModelProjectViewResponse> => {
    const response = await apiClient.get<PagedModelProjectViewResponse>("/projects", {
      params: { page, size },
    });
    return response.data;
  },

  getById: async (id: string): Promise<ProjectViewResponse> => {
    const response = await apiClient.get<ProjectViewResponse>(`/projects/${id}`);
    return response.data;
  },

  create: async (request: CreateProjectRequest): Promise<ProjectViewResponse> => {
    const response = await apiClient.post<ProjectViewResponse>("/projects", request);
    return response.data;
  },

  update: async (id: string, request: UpdateProjectRequest): Promise<ProjectViewResponse> => {
    const response = await apiClient.patch<ProjectViewResponse>(`/projects/${id}`, request);
    return response.data;
  },

  delete: async (id: string): Promise<void> => {
    await apiClient.delete(`/projects/${id}`);
  },

  archive: async (id: string): Promise<void> => {
    await apiClient.patch(`/projects/${id}/archive`);
  },

  unarchive: async (id: string): Promise<void> => {
    await apiClient.patch(`/projects/${id}/unarchive`);
  },
};
