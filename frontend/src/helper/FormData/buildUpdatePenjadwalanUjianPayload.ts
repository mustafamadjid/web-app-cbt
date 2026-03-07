import { toRfc3339Local } from "@/helper/dateFormatting/toRfc3339Local";
import type {
  BuatUjianFormValues,
  UpdatePenjadwalanUjianPayload,
} from "@/types/Ujian/BuatUjian";

const normalizeText = (value: string) => value.trim();
const normalizeToken = (value: string) => value.trim().toUpperCase();

export const buildUpdatePenjadwalanUjianPayload = (
  values: BuatUjianFormValues,
  initialValues: BuatUjianFormValues,
): UpdatePenjadwalanUjianPayload => {
  const payload: UpdatePenjadwalanUjianPayload = {};

  if (values.id_bank_soal !== initialValues.id_bank_soal) {
    payload.id_bank_soal = values.id_bank_soal;
  }

  if (values.id_kelas !== initialValues.id_kelas) {
    payload.id_kelas = values.id_kelas;
  }

  if (values.kelas_scope === "SPESIFIK") {
    if (values.id_nama_kelas !== initialValues.id_nama_kelas) {
      payload.id_nama_kelas = values.id_nama_kelas;
    }
  } else if (initialValues.kelas_scope === "SPESIFIK") {
    payload.id_nama_kelas = 0;
  }

  if (normalizeText(values.nama_ujian) !== normalizeText(initialValues.nama_ujian)) {
    payload.nama_ujian = normalizeText(values.nama_ujian);
  }

  if (
    normalizeText(values.deskripsi_ujian) !==
    normalizeText(initialValues.deskripsi_ujian)
  ) {
    payload.deskripsi_ujian = normalizeText(values.deskripsi_ujian);
  }

  if (values.acak_soal !== initialValues.acak_soal) {
    payload.acak_soal = values.acak_soal;
  }

  if (values.id_sesi !== initialValues.id_sesi) {
    payload.id_sesi = values.id_sesi;
  }

  if (values.id_ruangan !== initialValues.id_ruangan) {
    payload.id_ruangan = values.id_ruangan;
  }

  if (values.id_pengawas !== initialValues.id_pengawas) {
    payload.id_pengawas = values.id_pengawas;
  }

  if (values.tanggal_ujian !== initialValues.tanggal_ujian) {
    payload.tanggal_ujian = toRfc3339Local(values.tanggal_ujian, "00:00");
  }

  if (values.waktu_mulai !== initialValues.waktu_mulai) {
    payload.waktu_mulai = toRfc3339Local(values.tanggal_ujian, values.waktu_mulai);
  }

  if (values.waktu_selesai !== initialValues.waktu_selesai) {
    payload.waktu_selesai = toRfc3339Local(values.tanggal_ujian, values.waktu_selesai);
  }

  if (normalizeToken(values.token) !== normalizeToken(initialValues.token)) {
    payload.token = normalizeText(values.token);
  }

  return payload;
};
