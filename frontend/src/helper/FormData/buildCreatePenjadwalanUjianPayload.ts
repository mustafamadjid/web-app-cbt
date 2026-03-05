import { toRfc3339Local } from "@/helper/dateFormatting/toRfc3339Local";
import type {
  BuatUjianFormValues,
  CreatePenjadwalanUjianPayload,
} from "@/types/Ujian/BuatUjian";

type BuildCreatePenjadwalanUjianPayloadParams = {
  values: BuatUjianFormValues;
  idGuru: number;
};

export const buildCreatePenjadwalanUjianPayload = ({
  values,
  idGuru,
}: BuildCreatePenjadwalanUjianPayloadParams): CreatePenjadwalanUjianPayload => {
  const tanggalUjian = toRfc3339Local(values.tanggal_ujian, "00:00");
  const waktuMulai = toRfc3339Local(values.tanggal_ujian, values.waktu_mulai);
  const waktuSelesai = toRfc3339Local(values.tanggal_ujian, values.waktu_selesai);

  const payload: CreatePenjadwalanUjianPayload = {
    id_bank_soal: values.id_bank_soal,
    id_kelas: values.id_kelas,
    id_guru: idGuru,
    nama_ujian: values.nama_ujian.trim(),
    deskripsi_ujian: values.deskripsi_ujian.trim(),
    acak_soal: values.acak_soal,
    id_sesi: values.id_sesi,
    id_ruangan: values.id_ruangan,
    tanggal_ujian: tanggalUjian,
    waktu_mulai: waktuMulai,
    waktu_selesai: waktuSelesai,
    status_ujian: "BELUM_MULAI",
    token: values.token.trim(),
    id_pengawas: values.id_pengawas,
  };

  if (values.kelas_scope === "SPESIFIK" && values.id_nama_kelas > 0) {
    payload.id_nama_kelas = values.id_nama_kelas;
  }

  if (payload.deskripsi_ujian === "") {
    delete payload.deskripsi_ujian;
  }

  return payload;
};
