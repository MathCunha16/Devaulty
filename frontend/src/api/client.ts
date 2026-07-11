import axios, { AxiosError } from "axios";
import type { ApiErrorResponse } from "../types/api";

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

export const apiClient = axios.create({
  baseURL: "http://localhost:8080/api/v1",
  headers: {
    "Content-Type": "application/json",
  },
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
