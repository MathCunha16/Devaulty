import axios, { AxiosError } from "axios";
import type { ApiErrorResponse } from "../types/api";

declare global {
  interface Window {
    DEVAULTY_INTERNAL_TOKEN?: string;
    DEVAULTY_API_BASE_URL?: string;
  }
}

export class ApiError extends Error {
  status: number;
  payload: ApiErrorResponse | null;

  constructor(message: string, status: number, payload: ApiErrorResponse | null = null) {
    super(message);
    this.name = "ApiError";
    this.status = status;
    this.payload = payload;
  }
}

export const getInternalToken = (): string | undefined => {
  if (typeof window !== "undefined" && window.DEVAULTY_INTERNAL_TOKEN) {
    return window.DEVAULTY_INTERNAL_TOKEN;
  }
  // Strictly in development mode (import.meta.env.DEV), fallback to dev secret token.
  // In production builds, Vite tree-shakes and completely eliminates this fallback.
  if (import.meta.env.DEV) {
    return "dev-secret-token";
  }
  return undefined;
};

export const getApiBaseUrl = (): string => {
  if (typeof window !== "undefined" && window.DEVAULTY_API_BASE_URL) {
    return window.DEVAULTY_API_BASE_URL;
  }
  if (import.meta.env.VITE_API_BASE_URL) {
    return import.meta.env.VITE_API_BASE_URL;
  }
  // Always fall back to localhost:8080 — never use window.location.origin
  // because in a Tauri webview that would return "tauri://localhost" which
  // is not a valid HTTP base URL for fetch/axios requests.
  return "http://localhost:8080/api/v1";
};

export const apiClient = axios.create({
  baseURL: getApiBaseUrl(),
  headers: {
    "Content-Type": "application/json",
  },
});

apiClient.interceptors.request.use((config) => {
  config.baseURL = getApiBaseUrl();
  const internalToken = getInternalToken();

  if (internalToken) {
    config.headers["X-Devaulty-Internal-Token"] = internalToken;
  }

  return config;
});

apiClient.interceptors.response.use(
  (response) => response,
  (error: AxiosError<ApiErrorResponse>) => {
    if (error.response) {
      // The request was made and the server responded with a status code
      // that falls out of the range of 2xx
      const status = error.response.status;
      const data = error.response.data;
      const message = data?.message || error.message || "Request failed";
      
      throw new ApiError(message, status, data);
    } else if (error.request) {
      // The request was made but no response was received
      throw new ApiError("No response received from backend server", 503, null);
    } else {
      // Something happened in setting up the request that triggered an Error
      throw new ApiError(error.message || "Request setup error", 500, null);
    }
  }
);
