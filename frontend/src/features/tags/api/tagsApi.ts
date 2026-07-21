import { apiClient } from "@/api/client";
import type {
  CreateTagRequest,
  UpdateTagRequest,
  TagViewResponse,
} from "~types/api";

export const tagsApi = {
  getAllByProject: async (projectId: string): Promise<TagViewResponse[]> => {
    const response = await apiClient.get<TagViewResponse[]>(
      `/projects/${projectId}/tags`
    );
    return response.data;
  },

  create: async (projectId: string, request: CreateTagRequest): Promise<TagViewResponse> => {
    const response = await apiClient.post<TagViewResponse>(
      `/projects/${projectId}/tags`,
      request
    );
    return response.data;
  },

  update: async (
    projectId: string,
    tagId: string,
    request: UpdateTagRequest
  ): Promise<TagViewResponse> => {
    const response = await apiClient.patch<TagViewResponse>(
      `/projects/${projectId}/tags/${tagId}`,
      request
    );
    return response.data;
  },

  delete: async (projectId: string, tagId: string): Promise<void> => {
    await apiClient.delete(`/projects/${projectId}/tags/${tagId}`);
  },

  search: async (projectId: string, name: string): Promise<TagViewResponse[]> => {
    const response = await apiClient.get<TagViewResponse[]>(
      `/projects/${projectId}/tags/search`,
      { params: { name } }
    );
    return response.data;
  },

  associate: async (
    projectId: string,
    itemType: string,
    itemId: string,
    tagId: string
  ): Promise<void> => {
    await apiClient.put(
      `/projects/${projectId}/items/${itemType}/${itemId}/tags/${tagId}`
    );
  },

  disassociate: async (
    projectId: string,
    itemType: string,
    itemId: string,
    tagId: string
  ): Promise<void> => {
    await apiClient.delete(
      `/projects/${projectId}/items/${itemType}/${itemId}/tags/${tagId}`
    );
  },
};
