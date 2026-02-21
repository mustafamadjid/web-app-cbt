export type MediaKind = "image" | "document";

const PROTOCOL_REGEX = /^[a-zA-Z][a-zA-Z\d+\-.]*:/;

const normalizeUploadPath = (path: string, kind: MediaKind): string => {
  const normalizedPath = path.startsWith("/") ? path : `/${path}`;

  if (normalizedPath.startsWith("/uploads/pengumuman/")) {
    return `/uploads/document/${normalizedPath.slice("/uploads/".length)}`;
  }

  if (
    normalizedPath.startsWith("/uploads/") &&
    !normalizedPath.startsWith("/uploads/image/") &&
    !normalizedPath.startsWith("/uploads/document/")
  ) {
    const suffix = normalizedPath.slice("/uploads/".length);
    const prefix =
      kind === "document" ? "/uploads/document" : "/uploads/image";
    return `${prefix}/${suffix}`;
  }

  return normalizedPath;
};

export const resolveMediaUrl = (
  path: string | null | undefined,
  kind: MediaKind,
): string => {
  if (!path) return "";

  const trimmedPath = path.trim();
  if (!trimmedPath) return "";

  if (PROTOCOL_REGEX.test(trimmedPath) || trimmedPath.startsWith("//")) {
    return trimmedPath;
  }

  const apiBaseUrl = (import.meta.env.VITE_API_URL ?? "").replace(/\/+$/, "");
  const normalizedPath = normalizeUploadPath(trimmedPath, kind);

  return apiBaseUrl ? `${apiBaseUrl}${normalizedPath}` : normalizedPath;
};

export const resolveImageUrl = (path: string | null | undefined) =>
  resolveMediaUrl(path, "image");

export const resolveDocumentUrl = (path: string | null | undefined) =>
  resolveMediaUrl(path, "document");
