import React from "react";
import toast from "react-hot-toast";
import { ArrowLeft, RefreshCw } from "lucide-react";
import { Link, useParams } from "react-router";

import HasilJawabanUjianContent from "@/components/features/Ujian/HasilJawabanUjianContent";
import { useAuth } from "@/contexts/AuthContext";
import { paths } from "@/routes/paths";
import {
  useGetHasilJawabanUjian,
  useSubmitKoreksiEssay,
} from "@/services/Api/features-api/Ujian/hasilJawabanUjian.service";
import type { SubmitKoreksiEssayRequest } from "@/types/Ujian/SubmitKoreksiEssay";

const KoreksiHasilUjian: React.FC = () => {
  const { user } = useAuth();
  const { id, attemptId } = useParams();

  const parsedHasilId = Number(id);
  const parsedAttemptId = Number(attemptId);
  const isHasilIdValid = Number.isInteger(parsedHasilId) && parsedHasilId > 0;
  const isAttemptIdValid =
    Number.isInteger(parsedAttemptId) && parsedAttemptId > 0;
  const canGradeEssay = user?.role !== "SISWA";
  const backPath = isHasilIdValid
    ? (
        user?.role === "GURU"
          ? paths.dashboard.hasil_ujian_detail_guru
          : paths.dashboard.hasil_ujian_detail
      ).replace(":id", String(parsedHasilId))
    : user?.role === "GURU"
      ? paths.dashboard.hasil_ujian_guru
      : paths.dashboard.hasil_ujian;

  const {
    data,
    loading,
    error,
    refetch,
  } = useGetHasilJawabanUjian(parsedAttemptId, isAttemptIdValid);
  const {
    execute: executeSubmitKoreksiEssay,
    loading: submitLoading,
    error: submitError,
  } = useSubmitKoreksiEssay();

  const handleSubmitEssayCorrection = async (
    payload: SubmitKoreksiEssayRequest,
  ) => {
    try {
      await executeSubmitKoreksiEssay(payload);
      toast.success("Koreksi essay berhasil disimpan.");
      await refetch();
    } catch (submitErr) {
      const message =
        submitErr instanceof Error
          ? submitErr.message
          : "Gagal menyimpan koreksi essay.";

      toast.error(message);
    }
  };

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
            Kembali ke detail hasil ujian
          </Link>
          <p className="mt-3 text-xs font-semibold uppercase tracking-wide text-slate-400">
            Attempt ID
          </p>
          <h1 className="mt-1 text-2xl font-bold text-slate-800">
            {isAttemptIdValid ? parsedAttemptId : "-"}
          </h1>
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

      {!isAttemptIdValid ? (
        <div className="rounded-2xl border border-rose-200 bg-rose-50 p-6 text-sm font-medium text-rose-700">
          ID attempt tidak valid.
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
        <>
          {submitError && (
            <div className="rounded-2xl border border-rose-200 bg-rose-50 p-4 text-sm text-rose-700">
              Gagal menyimpan koreksi essay: {submitError}
            </div>
          )}

          <HasilJawabanUjianContent
            title={canGradeEssay ? "Koreksi Hasil Ujian" : "Detail Hasil Ujian"}
            subtitle={
              canGradeEssay
                ? "Tinjau semua jawaban siswa dan koreksi essay yang belum dinilai."
                : "Tinjau semua jawaban ujian beserta hasil penilaiannya."
            }
            nilaiAkhir={data.nilai_akhir}
            hasilJawabanUjian={data.hasil_jawaban}
            canGradeEssay={canGradeEssay}
            submitEssayCorrection={
              canGradeEssay ? handleSubmitEssayCorrection : undefined
            }
            submitDisabled={submitLoading}
            submitLoading={submitLoading}
          />
        </>
      ) : (
        <div className="rounded-2xl border border-dashed border-slate-200 bg-white p-6 text-center text-sm text-slate-500">
          Hasil jawaban ujian tidak ditemukan.
        </div>
      )}
    </div>
  );
};

export default KoreksiHasilUjian;
