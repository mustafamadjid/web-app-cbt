import { useEffect, useMemo, useRef, useState } from "react";
import { Link, useParams } from "react-router";
import { ArrowLeft, Users } from "lucide-react";
import {
  getHasilUjianDetail,
  type HasilUjianDetailResponse,
} from "@/services/Api/features-api/Ujian/hasilUjian.service";
import type { HasilUjianSiswa } from "@/types/Ujian/HasilUjian";
import { paths } from "@/routes/paths";

const HasilUjianDetail = () => {
  const { id } = useParams();
  const [loading, setLoading] = useState(false);
  const [errorMsg, setErrorMsg] = useState("");
  const [statistik, setStatistik] = useState<
    HasilUjianDetailResponse["statistik"] | null
  >(null);
  const [daftarSiswa, setDaftarSiswa] = useState<HasilUjianSiswa[]>([]);
  const requestSeq = useRef(0);

  const ujianId = useMemo(() => Number(id), [id]);

  useEffect(() => {
    if (!ujianId || Number.isNaN(ujianId)) return;
    const seq = ++requestSeq.current;
    (async () => {
      try {
        setLoading(true);
        setErrorMsg("");
        const data = await getHasilUjianDetail(ujianId);
        if (seq !== requestSeq.current) return;
        setStatistik(data.statistik);
        setDaftarSiswa(data.siswa);
      } catch {
        if (seq !== requestSeq.current) return;
        setErrorMsg("Gagal memuat detail hasil ujian.");
        setStatistik(null);
        setDaftarSiswa([]);
      } finally {
        if (seq !== requestSeq.current) return;
        setLoading(false);
      }
    })();
  }, [ujianId]);

  const buildKelasLabel = (siswa: HasilUjianSiswa) =>
    `Kelas ${siswa.id_tingkat_kelas} • ${siswa.id_nama_kelas}`;

  return (
    <div className="mx-auto flex max-w-7xl flex-col gap-6 px-4 py-10 sm:px-8">
      <div className="rounded-2xl border border-slate-200 bg-white p-6 shadow-sm">
        <div className="flex flex-col gap-4 sm:flex-row sm:items-center sm:justify-between">
          <div>
            <Link
              to={paths.dashboard.hasil_ujian}
              className="mb-3 inline-flex items-center gap-2 text-sm font-semibold text-[#397e50] hover:text-[#2d633f]"
            >
              <ArrowLeft size={16} />
              Kembali ke daftar hasil
            </Link>
            <h1 className="text-2xl font-bold text-slate-800">
              Detail Hasil Ujian
            </h1>
            <p className="text-sm text-slate-500">
              Statistik nilai dan daftar peserta yang mengikuti ujian.
            </p>
          </div>
          <div className="flex items-center gap-2 rounded-xl border border-emerald-100 bg-emerald-50 px-4 py-2 text-sm font-semibold text-emerald-700">
            <Users size={16} />
            Total peserta: {statistik?.jumlah_peserta ?? "-"}
          </div>
        </div>
      </div>

      {errorMsg ? (
        <div className="rounded-2xl border border-rose-200 bg-rose-50 p-6 text-sm font-semibold text-rose-700">
          {errorMsg}
        </div>
      ) : null}

      {loading ? (
        <div className="rounded-2xl border border-slate-200 bg-white p-6 text-sm text-slate-500">
          Memuat detail hasil ujian...
        </div>
      ) : null}

      <div className="grid gap-4 md:grid-cols-2 xl:grid-cols-4">
        <div className="rounded-2xl border border-slate-200 bg-white p-5 shadow-sm">
          <p className="text-xs font-semibold uppercase tracking-wider text-slate-400">
            Nilai Tertinggi
          </p>
          <p className="mt-3 text-2xl font-bold text-slate-800">
            {statistik?.nilai_tertinggi ?? "-"}
          </p>
        </div>
        <div className="rounded-2xl border border-slate-200 bg-white p-5 shadow-sm">
          <p className="text-xs font-semibold uppercase tracking-wider text-slate-400">
            Nilai Terendah
          </p>
          <p className="mt-3 text-2xl font-bold text-slate-800">
            {statistik?.nilai_terendah ?? "-"}
          </p>
        </div>
        <div className="rounded-2xl border border-slate-200 bg-white p-5 shadow-sm">
          <p className="text-xs font-semibold uppercase tracking-wider text-slate-400">
            Nilai Rata-rata
          </p>
          <p className="mt-3 text-2xl font-bold text-slate-800">
            {statistik?.rata_rata ?? "-"}
          </p>
        </div>
        <div className="rounded-2xl border border-slate-200 bg-white p-5 shadow-sm">
          <p className="text-xs font-semibold uppercase tracking-wider text-slate-400">
            Jumlah Peserta
          </p>
          <p className="mt-3 text-2xl font-bold text-slate-800">
            {statistik?.jumlah_peserta ?? "-"}
          </p>
        </div>
      </div>

      <div className="overflow-hidden rounded-2xl border border-slate-200 bg-white shadow-sm">
        <div className="border-b border-slate-100 px-6 py-4">
          <h2 className="text-lg font-bold text-slate-800">
            Daftar Siswa Peserta
          </h2>
          <p className="text-sm text-slate-500">
            Data siswa beserta ringkasan nilai ujian.
          </p>
        </div>

        <div className="overflow-x-auto">
          <table className="min-w-full divide-y divide-slate-100">
            <thead className="bg-slate-50">
              <tr className="text-left text-xs font-semibold uppercase tracking-wider text-slate-400">
                <th className="px-6 py-3">#</th>
                <th className="px-6 py-3">Siswa</th>
                <th className="px-6 py-3">Jenis Kelamin</th>
                <th className="px-6 py-3">No Absen</th>
                <th className="px-6 py-3">Angkatan</th>
                <th className="px-6 py-3">Tempat Lahir</th>
                <th className="px-6 py-3">Tanggal Lahir</th>
                <th className="px-6 py-3">Kelas</th>
                <th className="px-6 py-3">Nilai</th>
                <th className="px-6 py-3">Benar</th>
                <th className="px-6 py-3">Salah</th>
                <th className="px-6 py-3">Kosong</th>
                <th className="px-6 py-3">Aksi</th>
              </tr>
            </thead>
            <tbody className="divide-y divide-slate-100 text-sm text-slate-600">
              {daftarSiswa.length === 0 ? (
                <tr>
                  <td
                    colSpan={13}
                    className="px-6 py-6 text-center text-sm text-slate-500"
                  >
                    Belum ada data siswa.
                  </td>
                </tr>
              ) : (
                daftarSiswa.map((siswa, index) => (
                  <tr key={siswa.id} className="hover:bg-slate-50">
                    <td className="px-6 py-4 font-semibold text-slate-700">
                      {index + 1}
                    </td>
                    <td className="px-6 py-4">
                      <div className="flex items-center gap-3">
                        <img
                          src={siswa.urlGambarProfil}
                          alt={siswa.namaLengkap}
                          className="h-10 w-10 rounded-full object-cover"
                        />
                        <div>
                          <p className="font-semibold text-slate-800">
                            {siswa.namaLengkap}
                          </p>
                          <p className="text-xs text-slate-400">{siswa.role}</p>
                        </div>
                      </div>
                    </td>
                    <td className="px-6 py-4">{siswa.jenisKelamin}</td>
                    <td className="px-6 py-4">{siswa.noAbsen}</td>
                    <td className="px-6 py-4">{siswa.angkatan}</td>
                    <td className="px-6 py-4">{siswa.tempatLahir}</td>
                    <td className="px-6 py-4">{siswa.tanggalLahir}</td>
                    <td className="px-6 py-4">{buildKelasLabel(siswa)}</td>
                    <td className="px-6 py-4 font-semibold text-slate-700">
                      {siswa.nilai ?? "-"}
                    </td>
                    <td className="px-6 py-4">{siswa.jumlah_benar ?? "-"}</td>
                    <td className="px-6 py-4">{siswa.jumlah_salah ?? "-"}</td>
                    <td className="px-6 py-4">{siswa.jumlah_kosong ?? "-"}</td>
                    <td className="px-6 py-4">
                      <Link
                        to={paths.dashboard.hasil_ujian_siswa_detail
                          .replace(":id", String(ujianId))
                          .replace(":siswaId", String(siswa.id))}
                        className="inline-flex items-center rounded-lg border border-[#397e50]/30 bg-[#397e50]/10 px-3 py-1.5 text-xs font-semibold text-[#397e50] hover:border-[#397e50] hover:bg-[#397e50] hover:text-white"
                      >
                        Selengkapnya
                      </Link>
                    </td>
                  </tr>
                ))
              )}
            </tbody>
          </table>
        </div>
      </div>
    </div>
  );
};

export default HasilUjianDetail;
