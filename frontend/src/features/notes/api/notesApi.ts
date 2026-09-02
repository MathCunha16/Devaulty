import { apiClient } from "@/api/client";
import type {
  CreateNoteRequest,
  UpdateNoteRequest,
  NoteViewResponse,
  PagedModelNoteSummaryResponse,
} from "~types/api";

export const notesApi = {
  getAllByProject: async (
    projectId: string,
    page = 0,
    size = 100
  ): Promise<PagedModelNoteSummaryResponse> => {
    const response = await apiClient.get<PagedModelNoteSummaryResponse>(
      `/projects/${projectId}/notes`,
      { params: { page, size } }
    );
    return response.data;
  },

  getById: async (projectId: string, noteId: string): Promise<NoteViewResponse> => {
    const response = await apiClient.get<NoteViewResponse>(
      `/projects/${projectId}/notes/${noteId}`
    );
    return response.data;
  },

  create: async (projectId: string, request: CreateNoteRequest): Promise<NoteViewResponse> => {
    const response = await apiClient.post<NoteViewResponse>(
      `/projects/${projectId}/notes`,
      request
    );
    return response.data;
  },

  update: async (
    projectId: string,
    noteId: string,
    request: UpdateNoteRequest
  ): Promise<NoteViewResponse> => {
    const response = await apiClient.patch<NoteViewResponse>(
      `/projects/${projectId}/notes/${noteId}`,
      request
    );
    return response.data;
  },

  delete: async (projectId: string, noteId: string): Promise<void> => {
    await apiClient.delete(`/projects/${projectId}/notes/${noteId}`);
  },

  archive: async (projectId: string, noteId: string): Promise<void> => {
    await apiClient.patch(`/projects/${projectId}/notes/${noteId}/archive`);
  },

  unarchive: async (projectId: string, noteId: string): Promise<void> => {
    await apiClient.patch(`/projects/${projectId}/notes/${noteId}/unarchive`);
  },
};
