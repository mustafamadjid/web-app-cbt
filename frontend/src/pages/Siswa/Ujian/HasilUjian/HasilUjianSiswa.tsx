import React from "react";
import { Link } from "react-router";
import {
  ArrowRight,
  CalendarDays,
  Clock3,
  RefreshCw,
  ShieldCheck,
  SquareLibrary,
  Tag,
  Users,
} from "lucide-react";

import { useAuth } from "@/contexts/AuthContext";
import { paths } from "@/routes/paths";
import { useGetHasilUjianSiswaList } from "@/services/Api/features-api/Ujian/hasilUjianSiswa.service";

const HasilUjianSiswa: React.FC = () => {
  const { user, status: authStatus } = useAuth();
  const siswaId =
    user?.role === "SISWA" && typeof user.id_pengguna === "number"
      ? user.id_pengguna
      : -1;

  const {
    data: items,
    loading,
    error,
    refetch,
  } = useGetHasilUjianSiswaList(siswaId, authStatus === "authenticated");

  const hasilUjian = items ?? [];
  const canFetch = authStatus === "authenticated" && siswaId > 0;

  return (
    <div className="mx-auto max-w-7xl px-4 py-10 sm:px-8">
      <div className="space-y-6">
        <header className="rounded-2xl border border-slate-200 bg-white p-6 shadow-sm">
          <div className="flex flex-col gap-4 lg:flex-row lg:items-start lg:justify-between">
            <div className="space-y-2">
              <div className="inline-flex items-center gap-2 rounded-full   px-3 py-1 text-xs font-bold uppercase tracking-[0.18em] text-green-700">
                <ShieldCheck size={14} />
                Hasil Ujian Siswa
              </div>
              <h1 className="text-2xl font-bold text-slate-800">
                Daftar ujian yang sudah kamu submit
              </h1>
              <p className="max-w-3xl text-sm text-slate-500">
                Halaman ini menampilkan ujian yang sudah diserahkan berdasarkan
                attempt yang berhasil disubmit.
              </p>
            </div>

            <div className="flex flex-wrap items-center gap-2">
              <span className="inline-flex items-center gap-2 rounded-full border border-slate-200 bg-slate-50 px-3 py-1 text-xs font-semibold text-slate-600">
                <Users size={14} />
                Siswa: {user?.username ?? "-"}
              </span>
              <span className="inline-flex items-center gap-2 rounded-full border px-3 py-1 text-xs font-semibold text-emerald-700">
                <SquareLibrary size={14} />
                Total: {hasilUjian.length}
              </span>
              <button
                type="button"
                onClick={() => void refetch()}
                disabled={loading || !canFetch}
                className="inline-flex items-center cursor-pointer gap-2 rounded-xl border border-slate-200 px-3 py-2 text-xs font-semibold text-slate-600 transition hover:border-emerald-400 hover:text-emerald-700 disabled:cursor-not-allowed disabled:opacity-60"
              >
                <RefreshCw
                  size={14}
                  className={loading ? "animate-spin" : ""}
                />
                Refresh
              </button>
            </div>
          </div>
        </header>

        {!canFetch ? (
          <div className="rounded-2xl border border-dashed border-slate-200 bg-white p-6 text-center text-sm text-slate-500">
            Memuat data pengguna...
          </div>
        ) : error ? (
          <div className="rounded-2xl border border-rose-200 bg-rose-50 p-6 text-sm font-medium text-rose-700">
            <div className="flex flex-col gap-3 sm:flex-row sm:items-center sm:justify-between">
              <span>{error}</span>
              <button
                type="button"
                onClick={() => void refetch()}
                className="inline-flex items-center gap-2 self-start rounded-xl border border-rose-200 bg-white px-4 py-2 text-xs font-semibold text-rose-700 transition hover:border-rose-300 hover:bg-rose-100"
              >
                <RefreshCw size={14} />
                Coba Lagi
              </button>
            </div>
          </div>
        ) : loading ? (
          <div className="rounded-2xl border border-dashed border-slate-200 bg-white p-6 text-sm text-slate-500">
            Memuat daftar ujian yang sudah disubmit...
          </div>
        ) : hasilUjian.length === 0 ? (
          <div className="rounded-2xl border border-dashed border-slate-200 bg-white p-6 text-center text-sm text-slate-500">
            Belum ada ujian yang kamu submit.
          </div>
        ) : (
          <section className="grid gap-6 lg:grid-cols-2">
            {hasilUjian.map((item) => {
              const detailPath =
                paths.dashboard.hasil_ujian_detail_siswa.replace(
                  ":id",
                  String(item.id),
                ).replace(":attemptId", String(item.id_attempt));

              return (
                <article
                  key={`${item.id_attempt}-${item.id}`}
                  className="group overflow-hidden rounded-2xl border border-slate-200 bg-white shadow-sm transition hover:border-[#397e50]"
                >
                  <div className="h-1.5 " />

                  <div className="p-6">
                    <div className="flex flex-col gap-4 sm:flex-row sm:items-start sm:justify-between">
                      <div className="space-y-3">
                        <div>
                          <h2 className="text-xl font-bold text-slate-800">
                            {item.nama_ujian}
                          </h2>
                          <p className="mt-1 text-sm font-medium text-emerald-700">
                            Kelas {item.tingkat_kelas} - {item.nama_kelas}
                          </p>
                        </div>

                        <p className="max-w-2xl text-sm text-slate-500">
                          {item.deskripsi_ujian || "Tidak ada deskripsi ujian."}
                        </p>
                      </div>

                      <div className="rounded-2xl border border-slate-200 bg-slate-50 px-4 py-3">
                        <p className="text-xs font-semibold uppercase tracking-[0.18em] text-slate-400">
                          Pengawas
                        </p>
                        <p className="mt-1 text-sm font-bold text-slate-700">
                          {item.pengawas_nama_lengkap || item.pengawas_ujian}
                        </p>
                        
                      </div>
                    </div>

                    <div className="mt-6 grid gap-3 rounded-2xl border border-slate-100 bg-slate-50/60 p-4 sm:grid-cols-2 xl:grid-cols-4">
                      <div className="space-y-1">
                        <p className="flex items-center gap-1.5 text-2xs font-bold uppercase tracking-wider text-slate-400">
                          <CalendarDays
                            size={12}
                            className="text-emerald-600"
                          />
                          Tanggal
                        </p>
                        <p className="text-sm font-bold text-slate-700">
                          {item.tgl_ujian}
                        </p>
                      </div>

                      <div className="space-y-1">
                        <p className="flex items-center gap-1.5 text-2xs font-bold uppercase tracking-wider text-slate-400">
                          <Clock3 size={12} className="text-emerald-600" />
                          Waktu
                        </p>
                        <p className="text-sm font-bold text-slate-700">
                          {item.waktu_mulai} - {item.waktu_selesai}
                        </p>
                      </div>

                      <div className="space-y-1">
                        <p className="flex items-center gap-1.5 text-2xs font-bold uppercase tracking-wider text-slate-400">
                          <Users size={12} className="text-emerald-600" />
                          Sesi
                        </p>
                        <p className="text-sm font-bold text-slate-700">
                          {item.nama_sesi}
                        </p>
                      </div>

                      <div className="space-y-1">
                        <p className="flex items-center gap-1.5 text-2xs font-bold uppercase tracking-wider text-slate-400">
                          <Tag size={12} className="text-emerald-600" />
                          Ruang
                        </p>
                        <p className="text-sm font-bold text-slate-700">
                          {item.ruang_ujian}
                        </p>
                      </div>
                    </div>

                    <div className="mt-6 flex flex-col gap-3 border-t border-slate-100 pt-5 sm:flex-row sm:items-center sm:justify-between">
                      <Link
                        to={detailPath}
                        className="inline-flex items-center justify-center gap-2 rounded-xl bg-[#397e50] px-5 py-2.5 text-xs font-bold uppercase tracking-widest text-white transition hover:bg-emerald-700 hover:shadow-lg hover:shadow-emerald-900/20"
                      >
                        Lihat Detail
                        <ArrowRight size={14} />
                      </Link>
                    </div>
                  </div>
                </article>
              );
            })}
          </section>
        )}
      </div>
    </div>
  );
};

export default HasilUjianSiswa;
