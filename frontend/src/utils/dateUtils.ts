export const formatUpdatedDate = (
  updatedAt?: string | null,
  createdAt?: string | null
): string | null => {
  if (!updatedAt) return null;

  const updatedTime = new Date(updatedAt).getTime();
  if (isNaN(updatedTime) || updatedTime <= 0 || new Date(updatedAt).getFullYear() <= 1970) {
    return null;
  }

  if (createdAt) {
    const createdTime = new Date(createdAt).getTime();
    if (!isNaN(createdTime) && Math.abs(updatedTime - createdTime) < 2000) {
      // Created and Updated at virtually the same time, return null to suppress
      return null;
    }
  }

  return new Date(updatedAt).toLocaleDateString();
};
