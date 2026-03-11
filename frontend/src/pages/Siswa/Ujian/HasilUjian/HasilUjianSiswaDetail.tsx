import React from "react";
import { useNavigate, useParams } from "react-router";
import { ArrowLeft, Award, CheckCircle2, CalendarDays, XCircle } from "lucide-react";
import { useGetUjianSiswaResultDetail } from "@/services/Api/features-api/Ujian/ujianSiswa.service";
import { paths } from "@/routes/paths";

const HasilUjianSiswaDetail: React.FC = () => {
  const navigate = useNavigate();
  const { id } = useParams();
  const detailId = Number(id);
  const isDetailIdValid = Number.isFinite(detailId);
  const {
    data: detail,
    loading,
  } = useGetUjianSiswaResultDetail(isDetailIdValid ? detailId : -1);

  return (
    <div className="space-y-6 px-8 py-10">
      <header className="rounded-xl border border-gray-200 bg-white p-6 shadow-sm">
        <button
          type="button"
          onClick={() => navigate(paths.dashboard.hasil_ujian_siswa)}
          className="flex items-center gap-2 text-sm font-semibold text-[#397e50]"
        >
          <ArrowLeft className="h-4 w-4" />
          Kembali ke hasil ujian
        </button>
        <h1 className="mt-4 text-2xl font-bold text-[#37513d]">
          Detail Hasil Ujian
        </h1>
        <p className="mt-2 text-sm text-gray-500">
          Rincian nilai dan performa ujian yang sudah diselesaikan.
        </p>
      </header>

      {loading ? (
        <div className="rounded-xl border border-dashed border-gray-200 bg-white p-6 text-center text-sm text-gray-500">
          Memuat detail hasil ujian...
        </div>
      ) : detail ? (
        <section className="rounded-xl border border-gray-200 bg-white p-6 shadow-sm">
          <div className="flex flex-wrap items-start justify-between gap-4">
            <div>
              <p className="text-xs font-semibold uppercase tracking-wide text-gray-400">
                {detail.mapel}
              </p>
              <h2 className="text-xl font-bold text-[#37513d]">
                {detail.nama_ujian}
              </h2>
              <div className="mt-2 flex items-center gap-2 text-sm text-gray-500">
                <CalendarDays className="h-4 w-4 text-[#397e50]" />
                <span>{detail.tgl_ujian}</span>
              </div>
            </div>
            <div className="rounded-xl border border-[#397e50]/20 bg-[#397e50]/10 px-4 py-3 text-center">
              <p className="text-xs font-semibold text-[#397e50]">Total Nilai</p>
              <p className="text-2xl font-bold text-[#37513d]">
                {detail.nilai}
              </p>
            </div>
          </div>

          <div className="mt-6 grid gap-4 md:grid-cols-3">
            <div className="rounded-xl border border-gray-200 p-4">
              <div className="flex items-center gap-2 text-sm font-semibold text-gray-500">
                <CheckCircle2 className="h-4 w-4 text-[#397e50]" />
                Jumlah Benar
              </div>
              <p className="mt-2 text-2xl font-bold text-[#37513d]">
                {detail.jumlah_benar}
              </p>
            </div>
            <div className="rounded-xl border border-gray-200 p-4">
              <div className="flex items-center gap-2 text-sm font-semibold text-gray-500">
                <XCircle className="h-4 w-4 text-[#397e50]" />
                Jumlah Salah
              </div>
              <p className="mt-2 text-2xl font-bold text-[#37513d]">
                {detail.jumlah_salah}
              </p>
            </div>
            <div className="rounded-xl border border-gray-200 p-4">
              <div className="flex items-center gap-2 text-sm font-semibold text-gray-500">
                <Award className="h-4 w-4 text-[#397e50]" />
                Status Nilai
              </div>
              <p className="mt-2 text-2xl font-bold text-[#37513d]">
                {detail.nilai >= 75 ? "Tuntas" : "Perlu Remedial"}
              </p>
            </div>
          </div>
        </section>
      ) : (
        <div className="rounded-xl border border-dashed border-gray-200 bg-white p-6 text-center text-sm text-gray-500">
          {!isDetailIdValid
            ? "ID detail hasil ujian tidak valid."
            : "Detail hasil ujian tidak ditemukan."}
        </div>
      )}
    </div>
  );
};

export default HasilUjianSiswaDetail;
