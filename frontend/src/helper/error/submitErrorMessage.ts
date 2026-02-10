import { ApiError } from "@/services/Api/api";

type ErrorMessageOptions = {
  defaultMessage: string;
  codeMap?: Record<string, string>;
  messageMap?: Record<string, string>;
};

const COMMON_CODE_MAP: Record<string, string> = {
  USERNAME_TAKEN: "Username sudah digunakan.",
  NIP_TAKEN: "NIP sudah digunakan.",
  NISN_TAKEN: "NISN sudah digunakan.",
  NO_FIELD_TO_UPDATE: "Tidak ada data yang diubah.",
  INVALID_INPUT: "Data yang diinput tidak valid.",
  NOT_FOUND: "Data tidak ditemukan.",
  FORBIDDEN: "Anda tidak memiliki akses untuk melakukan aksi ini.",
};

export function getSubmitErrorMessage(
  error: unknown,
  { defaultMessage, codeMap = {}, messageMap = {} }: ErrorMessageOptions,
): string {
  if (!(error instanceof ApiError)) {
    return defaultMessage;
  }

  const normalizedMessage = error.message.trim().toLowerCase();
  const normalizedMessageMap = Object.entries(messageMap).reduce<Record<string, string>>(
    (acc, [key, value]) => {
      acc[key.trim().toLowerCase()] = value;
      return acc;
    },
    {},
  );

  return (
    (error.code ? codeMap[error.code] : undefined) ??
    (error.code ? COMMON_CODE_MAP[error.code] : undefined) ??
    normalizedMessageMap[normalizedMessage] ??
    defaultMessage
  );
}
