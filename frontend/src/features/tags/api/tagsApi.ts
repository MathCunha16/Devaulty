import { apiClient } from "@/api/client";
import type {
  CreateTagRequest,
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

  delete: async (projectId: string, tagId: string): Promise<void> => {
    await apiClient.delete(`/projects/${projectId}/tags/${tagId}`);
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
