import { api, type ApiEnvelope } from "../../api";
import { buildJsonData } from "@/helper/FormData/BuildJsonData";

import type {
  BuatUjianFormValues,
  BuatUjianSubmitResponse,
  SiswaPreviewItem,
} from "@/types/Ujian/BuatUjian";
import { DUMMY_SISWA } from "@/services/Api/features-api/KelolaAkun/akunsiswa.service";
import { getTingkatKelasById } from "@/services/Api/features-api/DataMaster/kelas.service";

const sleep = (ms: number) => new Promise((resolve) => setTimeout(resolve, ms));

export async function getUjianSiswaPreview(params: {
  // Data dari server harusnya di sorting duluan berdasarkan nama nomor absen dan kelas
  tingkatKelasId?: number;
}): Promise<SiswaPreviewItem[]> {
  await sleep(220);
  const tingkatKelas = getTingkatKelasById(params.tingkatKelasId);
  if (!tingkatKelas) return [];

  return DUMMY_SISWA.filter(
    (siswa) => siswa.id_tingkat_kelas === tingkatKelas,
  ).map((siswa) => ({
    id: siswa.id,
    nama: siswa.namaLengkap,
    username: siswa.username,
    no_absen: siswa.noAbsen,
    kelas: String(siswa.id_tingkat_kelas),
    status_akun: siswa.statusAkun,
  }));
}

export async function submitBuatUjian(values: BuatUjianFormValues) {
  const data = buildJsonData(values);
  const res = await api<ApiEnvelope<BuatUjianSubmitResponse>>("/ujian", {
    method: "POST",
    data,
  });

  return res.data;
}
