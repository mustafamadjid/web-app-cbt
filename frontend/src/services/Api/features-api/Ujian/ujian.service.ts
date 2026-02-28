import { api, type ApiEnvelope } from "../../api";
import { useFetch, usePost } from "@/hooks/fetch";
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
    (siswa) => siswa.tingkat_kelas === tingkatKelas,
  ).map((siswa) => ({
    id: siswa.id_pengguna,
    nama: siswa.nama_lengkap,
    username: siswa.username,
    no_absen: siswa.no_absen,
    kelas: String(siswa.tingkat_kelas),
    status_akun: siswa.status_akun,
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

// =====================
// Hook Wrappers
// =====================

export function useGetUjianSiswaPreview(params: { tingkatKelasId?: number }) {
  return useFetch(
    () => getUjianSiswaPreview(params),
    [params.tingkatKelasId],
  );
}

export function useSubmitBuatUjian() {
  return usePost((values: BuatUjianFormValues) => submitBuatUjian(values));
}
