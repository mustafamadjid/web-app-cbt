import React from "react";
import { ArrowLeft, RefreshCw } from "lucide-react";
import { Link, useParams } from "react-router";
import HasilJawabanUjianContent from "@/components/features/Ujian/HasilJawabanUjianContent";
import { useGetHasilJawabanUjian } from "@/services/Api/features-api/Ujian/hasilJawabanUjian.service";
import { paths } from "@/routes/paths";

const HasilUjianSiswaDetail: React.FC = () => {
  const { id, attemptId } = useParams();
  const parsedHasilId = Number(id);
  const parsedAttemptId = Number(attemptId);
  const isHasilIdValid = Number.isInteger(parsedHasilId) && parsedHasilId > 0;
  const isAttemptIdValid =
    Number.isInteger(parsedAttemptId) && parsedAttemptId > 0;
  const { data, loading, error, refetch } = useGetHasilJawabanUjian(
    parsedAttemptId,
    isAttemptIdValid,
  );

  const backPath = paths.dashboard.hasil_ujian_siswa;

  return (
    <div className="mx-auto flex max-w-7xl flex-col gap-6 px-4 py-8 sm:px-6 lg:px-8">
      <div className="flex flex-col gap-4 rounded-2xl border border-slate-200 bg-white p-6 shadow-sm sm:flex-row sm:items-center sm:justify-between">
        <div>
          <Link
            to={backPath}
            className="group inline-flex cursor-pointer items-center gap-2 text-sm font-medium text-slate-500 transition-colors hover:text-[#397e50]"
          >
            <ArrowLeft
              size={16}
              className="transition-transform group-hover:-translate-x-1"
            />
            Kembali ke hasil ujian
          </Link>
        </div>

        <button
          type="button"
          onClick={() => void refetch()}
          disabled={!isAttemptIdValid || loading}
          className="inline-flex cursor-pointer items-center justify-center gap-2 rounded-xl border border-slate-200 px-4 py-2 text-sm font-semibold text-slate-600 transition hover:border-[#397e50] hover:text-[#397e50] disabled:cursor-not-allowed disabled:opacity-60"
        >
          <RefreshCw size={16} className={loading ? "animate-spin" : ""} />
          Refresh
        </button>
      </div>

      {!isHasilIdValid || !isAttemptIdValid ? (
        <div className="rounded-2xl border border-rose-200 bg-rose-50 p-6 text-sm font-medium text-rose-700">
          {!isHasilIdValid
            ? "Data hasil ujian tidak valid."
            : "Data attempt ujian tidak valid."}
        </div>
      ) : loading ? (
        <div className="rounded-2xl border border-dashed border-slate-200 bg-white p-6 text-center text-sm text-slate-500">
          Memuat hasil jawaban ujian...
        </div>
      ) : error ? (
        <div className="rounded-2xl border border-rose-200 bg-rose-50 p-6 text-sm text-rose-700">
          Gagal memuat hasil jawaban ujian: {error}
        </div>
      ) : data ? (
        <HasilJawabanUjianContent
          title="Detail Hasil Ujian"
          subtitle="Tinjau semua jawaban ujian beserta hasil penilaiannya."
          nilaiAkhir={data.nilai_akhir}
          hasilJawabanUjian={data.hasil_jawaban}
          canGradeEssay={false}
        />
      ) : (
        <div className="rounded-2xl border border-dashed border-slate-200 bg-white p-6 text-center text-sm text-slate-500">
          Hasil jawaban ujian tidak ditemukan.
        </div>
      )}
    </div>
  );
};

export default HasilUjianSiswaDetail;
