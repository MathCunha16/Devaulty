import { apiClient } from "@/api/client";
import type {
  CreateSnippetRequest,
  SnippetViewResponse,
  UpdateSnippetRequest,
  PagedModelSnippetViewResponse,
} from "~types/api";

export const snippetsApi = {
  getAllByProject: async (
    projectId: string,
    page = 0,
    size = 100
  ): Promise<PagedModelSnippetViewResponse> => {
    const response = await apiClient.get<PagedModelSnippetViewResponse>(
      `/projects/${projectId}/snippets`,
      {
        params: { page, size },
      }
    );
    return response.data;
  },

  getById: async (projectId: string, snippetId: string): Promise<SnippetViewResponse> => {
    const response = await apiClient.get<SnippetViewResponse>(
      `/projects/${projectId}/snippets/${snippetId}`
    );
    return response.data;
  },

  create: async (projectId: string, request: CreateSnippetRequest): Promise<SnippetViewResponse> => {
    const response = await apiClient.post<SnippetViewResponse>(
      `/projects/${projectId}/snippets`,
      request
    );
    return response.data;
  },

  update: async (
    projectId: string,
    snippetId: string,
    request: UpdateSnippetRequest
  ): Promise<SnippetViewResponse> => {
    const response = await apiClient.patch<SnippetViewResponse>(
      `/projects/${projectId}/snippets/${snippetId}`,
      request
    );
    return response.data;
  },

  delete: async (projectId: string, snippetId: string): Promise<void> => {
    await apiClient.delete(`/projects/${projectId}/snippets/${snippetId}`);
  },
};
