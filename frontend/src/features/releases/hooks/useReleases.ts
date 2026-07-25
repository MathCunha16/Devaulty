import { useQuery, useMutation } from "@tanstack/react-query";
import { releasesApi } from "../api/releasesApi";

export const releaseKeys = {
  version: ["releases", "current-version"] as const,
  check: ["releases", "check"] as const,
};

export const useCurrentVersionQuery = () => {
  return useQuery({
    queryKey: releaseKeys.version,
    queryFn: releasesApi.getCurrentVersion,
    staleTime: 1000 * 60 * 60, // 1 hour
    retry: 1,
  });
};

export const useCheckUpdatesQuery = (enabled = true) => {
  return useQuery({
    queryKey: releaseKeys.check,
    queryFn: releasesApi.checkUpdates,
    enabled,
    staleTime: 1000 * 60 * 15, // 15 minutes
    retry: false,
  });
};

export const useCheckUpdatesMutation = () => {
  return useMutation({
    mutationFn: releasesApi.checkUpdates,
  });
};
