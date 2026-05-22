import {
  calculateDuration,
} from "@/helper/CalculateDuration/calculateDuration";
import {
  createValidator,
  matchesPattern,
  minNumber,
  requiredString,
  requiredValue,
} from "@/helper/validate/validateForm";
import type { BuatUjianFormValues } from "@/types/Ujian/BuatUjian";

export const sectionTitle = "text-sm font-semibold text-slate-800";
export const helperText = "text-xs text-slate-500";

export const selectBaseClass =
  "w-full rounded-lg border border-slate-200 bg-white px-3 py-2 text-sm outline-none transition focus:border-[#397e50] focus:ring-2 focus:ring-[#397e50] disabled:bg-slate-50 disabled:text-slate-500";

const timePatternRule = matchesPattern(
  /^([01]\d|2[0-3]):([0-5]\d)$/,
  "Gunakan format 24 jam HH:mm, contoh 07:30.",
);

export const validateBuatUjianForm = createValidator<BuatUjianFormValues>({
  nama_ujian: [requiredString("Nama ujian wajib diisi.")],
  deskripsi_ujian: [requiredString("Deskripsi ujian wajib diisi.")],
  id_kelas: [minNumber(1, "Tingkat kelas wajib dipilih.")],
  kelas_scope: [requiredValue("Cakupan kelas wajib dipilih.")],
  id_nama_kelas: [
    (value, currentValues) =>
      currentValues.kelas_scope === "SPESIFIK" && value === 0
        ? "Nama kelas wajib dipilih."
        : null,
  ],
  id_bank_soal: [minNumber(1, "Bank soal wajib dipilih.")],
  tanggal_ujian: [requiredString("Tanggal ujian wajib diisi.")],
  waktu_mulai: [requiredString("Waktu mulai wajib diisi."), timePatternRule],
  waktu_selesai: [
    requiredString("Waktu selesai wajib diisi."),
    timePatternRule,
    (_, currentValues) =>
      currentValues.waktu_mulai &&
      currentValues.waktu_selesai &&
      calculateDuration(currentValues.waktu_mulai, currentValues.waktu_selesai) <=
        0
        ? "Waktu selesai harus setelah waktu mulai."
        : null,
  ],
  id_ruangan: [minNumber(1, "Ruang ujian wajib dipilih.")],
  id_pengawas: [minNumber(1, "Guru pengawas wajib dipilih.")],
  id_sesi: [minNumber(1, "Sesi ujian wajib dipilih.")],
  token: [
    requiredString("Token ujian wajib diisi."),
    (value) =>
      value.trim().length > 30 ? "Token ujian maksimal 30 karakter." : null,
  ],
});
