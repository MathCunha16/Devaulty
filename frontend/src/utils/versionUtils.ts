/**
 * Normalizes a version string so it consistently starts with a single lowercase 'v'.
 * Trims surrounding whitespace and strips any leading 'v' or 'V' (including multiples like "vv0.1" or " V0.1 ").
 *
 * Examples:
 * - "0.1.4-alpha" -> "v0.1.4-alpha"
 * - "v0.1.4-alpha" -> "v0.1.4-alpha"
 * - "vv0.1.4-alpha" -> "v0.1.4-alpha"
 * - " V0.1.4 " -> "v0.1.4"
 * - "" / "   " -> ""
 */
export const formatVersionTag = (ver?: string): string => {
  if (!ver) return "";
  const trimmed = ver.trim();
  if (!trimmed) return "";
  const cleanVersion = trimmed.replace(/^v+/i, "");
  return `v${cleanVersion}`;
};
