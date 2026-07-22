import { apiClient } from "@/api/client";
import type {
  CreateCredentialRequest,
  UpdateCredentialRequest,
  CredentialViewResponse,
  PagedModelCredentialSummaryResponse,
} from "~types/api";

export const credentialsApi = {
  getAllByProject: async (
    projectId: string,
    page = 0,
    size = 100
  ): Promise<PagedModelCredentialSummaryResponse> => {
    const response = await apiClient.get<PagedModelCredentialSummaryResponse>(
      `/projects/${projectId}/credentials`,
      { params: { page, size } }
    );
    return response.data;
  },

  getById: async (
    projectId: string,
    credentialId: string
  ): Promise<CredentialViewResponse> => {
    const response = await apiClient.get<CredentialViewResponse>(
      `/projects/${projectId}/credentials/${credentialId}`
    );
    return response.data;
  },

  create: async (
    projectId: string,
    request: CreateCredentialRequest
  ): Promise<CredentialViewResponse> => {
    const response = await apiClient.post<CredentialViewResponse>(
      `/projects/${projectId}/credentials`,
      request
    );
    return response.data;
  },

  update: async (
    projectId: string,
    credentialId: string,
    request: UpdateCredentialRequest
  ): Promise<CredentialViewResponse> => {
    const response = await apiClient.patch<CredentialViewResponse>(
      `/projects/${projectId}/credentials/${credentialId}`,
      request
    );
    return response.data;
  },

  delete: async (projectId: string, credentialId: string): Promise<void> => {
    await apiClient.delete(`/projects/${projectId}/credentials/${credentialId}`);
  },
};
