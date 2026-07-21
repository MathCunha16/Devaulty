import { apiClient } from "@/api/client";
import type {
  CreateLinkRequest,
  UpdateLinkRequest,
  LinkViewResponse,
  PagedModelLinkViewResponse,
} from "~types/api";

export const linksApi = {
  getAllByProject: async (
    projectId: string,
    page = 0,
    size = 100
  ): Promise<PagedModelLinkViewResponse> => {
    const response = await apiClient.get<PagedModelLinkViewResponse>(
      `/projects/${projectId}/links`,
      { params: { page, size } }
    );
    return response.data;
  },

  getById: async (projectId: string, linkId: string): Promise<LinkViewResponse> => {
    const response = await apiClient.get<LinkViewResponse>(
      `/projects/${projectId}/links/${linkId}`
    );
    return response.data;
  },

  create: async (projectId: string, request: CreateLinkRequest): Promise<LinkViewResponse> => {
    const response = await apiClient.post<LinkViewResponse>(
      `/projects/${projectId}/links`,
      request
    );
    return response.data;
  },

  update: async (
    projectId: string,
    linkId: string,
    request: UpdateLinkRequest
  ): Promise<LinkViewResponse> => {
    const response = await apiClient.patch<LinkViewResponse>(
      `/projects/${projectId}/links/${linkId}`,
      request
    );
    return response.data;
  },

  delete: async (projectId: string, linkId: string): Promise<void> => {
    await apiClient.delete(`/projects/${projectId}/links/${linkId}`);
  },
};
