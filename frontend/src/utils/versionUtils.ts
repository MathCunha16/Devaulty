/**
 * Normalizes a version string so it consistently starts with a single lowercase 'v'.
 * Avoids double 'vv' if the input already contains a leading 'v' or 'V'.
 * Example: "0.1.4-alpha" -> "v0.1.4-alpha", "v0.1.4-alpha" -> "v0.1.4-alpha"
 */
export const formatVersionTag = (ver?: string): string => {
  if (!ver) return "";
  const trimmed = ver.trim();
  if (trimmed.toLowerCase().startsWith("v")) {
    return trimmed;
  }
  return `v${trimmed}`;
};
