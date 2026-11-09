import { apiClient } from "@/api/client";
import type {
  Drawing,
  CreateDrawingRequest,
  UpdateDrawingRequest,
  PagedModelDrawingSummaryResponse,
} from "~types/api";

export const drawingsApi = {
  getAllByProject: async (
    projectId: string,
    page = 0,
    size = 100
  ): Promise<PagedModelDrawingSummaryResponse> => {
    const response = await apiClient.get<PagedModelDrawingSummaryResponse>(
      `/projects/${projectId}/excalidraw-drawings`,
      { params: { page, size } }
    );
    return response.data;
  },

  getById: async (projectId: string, drawingId: string): Promise<Drawing> => {
    const response = await apiClient.get<Drawing>(
      `/projects/${projectId}/excalidraw-drawings/${drawingId}`
    );
    return response.data;
  },

  create: async (
    projectId: string,
    request: CreateDrawingRequest
  ): Promise<Drawing> => {
    const response = await apiClient.post<Drawing>(
      `/projects/${projectId}/excalidraw-drawings`,
      request
    );
    return response.data;
  },

  update: async (
    projectId: string,
    drawingId: string,
    request: UpdateDrawingRequest
  ): Promise<Drawing> => {
    const response = await apiClient.patch<Drawing>(
      `/projects/${projectId}/excalidraw-drawings/${drawingId}`,
      request
    );
    return response.data;
  },

  delete: async (projectId: string, drawingId: string): Promise<void> => {
    await apiClient.delete(
      `/projects/${projectId}/excalidraw-drawings/${drawingId}`
    );
  },
};
