export type ErrorAction =
  | "fetch"
  | "create"
  | "update"
  | "delete"
  | "submit"
  | "login";

export type ErrorMessageContext = {
  action?: ErrorAction;
  entity?: string;
  fallbackMessage?: string;
  fieldLabels?: Record<string, string>;
};

const DEFAULT_FIELD_LABELS: Record<string, string> = {
  username: "Username",
  password: "Password",
  email: "Email",
  no_hp: "Nomor HP",
  nip: "NIP",
  nisn: "NISN",
  kode_sesi: "Kode sesi",
  nama_sesi: "Nama sesi",
  kode_ruang: "Kode ruang ujian",
  nama_ruangan: "Nama ruangan",
  kode_mapel: "Kode mata pelajaran",
  nama_mapel: "Nama mata pelajaran",
  deskripsi: "Deskripsi",
  deskripsi_ujian: "Deskripsi ujian",
  nama_ujian: "Nama ujian",
  judul_pengumuman: "Judul pengumuman",
  isi_pengumuman: "Isi pengumuman",
  tanggal_rilis_pengumuman: "Tanggal rilis pengumuman",
  tanggal_selesai_pengumuman: "Tanggal selesai pengumuman",
  dokumen_pengumuman: "Dokumen pengumuman",
  logo_sekolah: "Logo sekolah",
  email_sekolah: "Email sekolah",
  no_telp_sekolah: "Nomor telepon sekolah",
  kepala_sekolah: "Nama kepala sekolah",
  waka_sekolah: "Nama wakil kepala sekolah",
  nama_sekolah: "Nama sekolah",
  alamat_sekolah: "Alamat sekolah",
  token: "Token ujian",
  token_ujian: "Token ujian",
  waktu_mulai: "Waktu mulai",
  waktu_selesai: "Waktu selesai",
  id_kelas: "Kelas",
  id_mapel: "Mata pelajaran",
  id_bank_soal: "Bank soal",
  id_ruangan: "Ruang ujian",
  id_sesi: "Sesi ujian",
  id_pengawas: "Guru pengawas",
  id_nama_kelas: "Nama kelas",
  tingkat_kelas: "Tingkat kelas",
  nama_kelas: "Nama kelas",
  file: "File",
};

const CODE_MESSAGES: Record<string, string> = {
  INVALID_CREDENTIALS: "Username atau password yang Anda masukkan salah.",
  HAS_SESSION:
    "Login gagal. Silakan logout terlebih dahulu pada perangkat sebelumnya.",
  SESSION_EXPIRED: "Sesi Anda telah berakhir. Silakan login kembali.",
  UNAUTHENTICATED: "Sesi Anda telah berakhir. Silakan login kembali.",
  INVALID_TOKEN: "Sesi Anda tidak valid. Silakan login kembali.",
  FORBIDDEN: "Anda tidak memiliki akses untuk melakukan tindakan ini.",
  DELETE_RESTRICTED:
    "Data tidak bisa dihapus karena masih digunakan di bagian lain.",
  RATE_LIMITED:
    "Terlalu banyak permintaan. Silakan tunggu sebentar lalu coba lagi.",
  INVALID_TOKEN_UJIAN: "Token ujian salah.",
  SISWA_NOT_ALLOWED: "Anda tidak diizinkan mengikuti ujian ini.",
  UJIAN_ATTEMPT_TIME_EXPIRED: "Ujian telah selesai.",
  DOUBLE_ATTEMPT_NOT_ALLOWED: "Anda sudah mengikuti ujian ini.",
  USERNAME_TAKEN: "Username sudah terdaftar. Gunakan username lain.",
  EMAIL_TAKEN: "Email sudah terdaftar. Gunakan email lain.",
  NO_HP_TAKEN: "Nomor HP sudah terdaftar. Gunakan nomor lain.",
  NISN_TAKEN: "NISN sudah terdaftar. Gunakan NISN lain.",
  NIP_TAKEN: "NIP sudah terdaftar. Gunakan NIP lain.",
};

function titleCase(value: string) {
  return value.charAt(0).toUpperCase() + value.slice(1);
}

function normalizeFieldKey(value: string) {
  return value.trim().toLowerCase().replace(/[\s-]+/g, "_");
}

function resolveFieldLabel(
  rawField: string,
  fieldLabels?: Record<string, string>,
): string | null {
  const normalizedField = normalizeFieldKey(rawField);
  const label =
    fieldLabels?.[normalizedField] ??
    fieldLabels?.[rawField] ??
    DEFAULT_FIELD_LABELS[normalizedField];

  if (label) return label;

  if (!/^[a-z_]+$/.test(normalizedField)) {
    return null;
  }

  return titleCase(normalizedField.replace(/_/g, " "));
}

function getBackendRawMessage(error: unknown): string | null {
  if (error instanceof Error) {
    const apiLikeError = error as Error & { rawMessage?: string };
    return apiLikeError.rawMessage?.trim() || error.message.trim() || null;
  }

  return null;
}

function buildEntityLabel(entity?: string) {
  if (!entity) return "data";
  return entity.trim().toLowerCase();
}

function buildActionFallback(
  action: ErrorAction | undefined,
  entity: string,
): string {
  switch (action) {
    case "create":
      return `Gagal menambahkan ${entity}.`;
    case "update":
      return `Gagal memperbarui ${entity}.`;
    case "delete":
      return `Gagal menghapus ${entity}.`;
    case "login":
      return "Login gagal. Silakan coba lagi.";
    case "submit":
      return `Gagal memproses ${entity}.`;
    case "fetch":
      return `Gagal memuat ${entity}.`;
    default:
      return "Terjadi kendala pada sistem. Silakan coba lagi.";
  }
}

export function isSensitiveBackendMessage(message: string): boolean {
  const normalized = message.trim().toLowerCase();

  if (!normalized) return true;

  return [
    "constraint violation",
    "sql",
    "database",
    "column",
    "relation",
    "postgres",
    "duplicate key",
    "foreign key",
    "syntax error",
    "internal server error",
    "invalid api response envelope",
    "api returned null data",
  ].some((term) => normalized.includes(term));
}

function mapRawMessage(
  rawMessage: string,
  context: ErrorMessageContext,
): string | null {
  const normalized = rawMessage.trim().toLowerCase();
  const entity = buildEntityLabel(context.entity);

  if (!normalized) return null;

  if (normalized === "file too large") {
    return "Ukuran file terlalu besar.";
  }

  if (normalized === "no fields to update") {
    return "Tidak ada perubahan data untuk disimpan.";
  }

  if (
    normalized === "data not found" ||
    normalized === "not found" ||
    normalized === "active attempt not found"
  ) {
    return `Data ${entity} tidak ditemukan.`;
  }

  if (normalized === "bad request: invalid date format") {
    return "Format tanggal belum sesuai.";
  }

  if (normalized === "bad request: missing fields") {
    return "Periksa kembali data yang wajib diisi.";
  }

  if (normalized === "bad request: invalid dokumen_pengumuman") {
    return "Dokumen pengumuman harus berupa PDF atau DOCX.";
  }

  if (normalized === "bad request: invalid token ujian") {
    return "Token ujian salah.";
  }

  const requiredFieldMatch = normalized.match(
    /(?:bad request:\s*)?([a-z_]+)\s+is required$/,
  );
  if (requiredFieldMatch) {
    const label = resolveFieldLabel(requiredFieldMatch[1], context.fieldLabels);
    return label ? `${label} wajib diisi.` : "Periksa kembali data yang diisi.";
  }

  const alreadyExistsMatch = normalized.match(
    /(?:bad request:\s*)?([a-z_]+)\s+already exist(?:s)?$/,
  );
  if (alreadyExistsMatch) {
    const label = resolveFieldLabel(alreadyExistsMatch[1], context.fieldLabels);
    return label
      ? `${label} sudah digunakan. Gunakan ${label.toLowerCase()} lain.`
      : "Data yang dimasukkan sudah terdaftar.";
  }

  if (
    normalized.startsWith("bad request: invalid") ||
    normalized === "invalid input"
  ) {
    if (isSensitiveBackendMessage(normalized)) {
      return "Periksa kembali data yang diisi.";
    }

    const invalidFieldMatch = normalized.match(
      /(?:bad request:\s*)?invalid\s+([a-z_]+)$/,
    );
    if (invalidFieldMatch) {
      const label = resolveFieldLabel(invalidFieldMatch[1], context.fieldLabels);
      return label
        ? `${label} yang dipilih atau diisi belum valid.`
        : "Periksa kembali data yang diisi.";
    }

    return "Periksa kembali data yang diisi.";
  }

  return null;
}

export function getUserFriendlyErrorMessage(
  error: unknown,
  context: ErrorMessageContext = {},
): string {
  if (typeof error === "string" && error.trim()) {
    return getUserFriendlyErrorMessage(new Error(error), context);
  }

  const entity = buildEntityLabel(context.entity);
  const fallbackMessage =
    context.fallbackMessage ?? buildActionFallback(context.action, entity);

  if (error instanceof Error) {
    const apiLikeError = error as Error & {
      code?: string;
      status?: number;
      rawMessage?: string;
    };

    if (apiLikeError.code && CODE_MESSAGES[apiLikeError.code]) {
      return CODE_MESSAGES[apiLikeError.code];
    }

    if (apiLikeError.code === "NOT_FOUND") {
      return `Data ${entity} tidak ditemukan.`;
    }

    if (
      apiLikeError.code === "CONFLICT" &&
      !isSensitiveBackendMessage(apiLikeError.rawMessage ?? error.message)
    ) {
      const rawMessage = getBackendRawMessage(error);
      if (rawMessage) {
        const mapped = mapRawMessage(rawMessage, context);
        if (mapped) return mapped;
      }

      return "Data yang dimasukkan sudah terdaftar.";
    }

    if (
      apiLikeError.code === "BAD_REQUEST" ||
      apiLikeError.code === "INVALID_INPUT"
    ) {
      const rawMessage = getBackendRawMessage(error);
      if (rawMessage) {
        const mapped = mapRawMessage(rawMessage, context);
        if (mapped) return mapped;
      }

      return "Periksa kembali data yang diisi.";
    }

    const rawMessage = getBackendRawMessage(error);
    if (rawMessage && !isSensitiveBackendMessage(rawMessage)) {
      const mapped = mapRawMessage(rawMessage, context);
      if (mapped) return mapped;
    }

    if (apiLikeError.status === 401) {
      return "Sesi Anda telah berakhir. Silakan login kembali.";
    }

    if (apiLikeError.status === 403) {
      return "Anda tidak memiliki akses untuk melakukan tindakan ini.";
    }

    if (apiLikeError.status === 404) {
      return `Data ${entity} tidak ditemukan.`;
    }

    return fallbackMessage;
  }

  return fallbackMessage;
}
