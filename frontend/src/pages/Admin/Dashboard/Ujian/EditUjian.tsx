import { Link, useNavigate, useParams } from "react-router";
import { useMemo } from "react";
import toast from "react-hot-toast";
import { AlertCircle, ArrowLeft } from "lucide-react";

import { useAuth } from "@/contexts/AuthContext";
import BuatUjianForm from "@/layouts/Form/Admin/Ujian/BuatUjianForm";
import { paths } from "@/routes/paths";
import { ApiError } from "@/services/Api/api";
import {
  useGetUjianEditFormData,
  useUpdateUjianPartial,
} from "@/services/Api/features-api/Ujian/jadwalujian.service";
import type { BuatUjianFormValues } from "@/types/Ujian/BuatUjian";

const EditUjian = () => {
  const { id } = useParams();
  const navigate = useNavigate();
  const { user } = useAuth();

  const jadwalId = useMemo(() => Number(id), [id]);
  const isJadwalIdValid = Boolean(id) && Number.isFinite(jadwalId) && jadwalId > 0;
  const goBackPath =
    user?.role === "GURU"
      ? paths.dashboard.jadwal_ujian_guru
      : paths.dashboard.jadwal_ujian;
  const { execute: executeUpdateUjianPartial } = useUpdateUjianPartial();

  const {
    data: initialData,
    loading,
    error,
  } = useGetUjianEditFormData(jadwalId, isJadwalIdValid);

  const handleSubmit = async (values: BuatUjianFormValues) => {
    if (!initialData) {
      throw new ApiError("Data ujian tidak ditemukan.");
    }

    const updated = await executeUpdateUjianPartial({
      idUjian: initialData.id_ujian,
      values,
      initialValues: initialData.values,
    });

    if (updated) {
      toast.success("Jadwal ujian berhasil diperbarui.");
    } else {
      toast("Tidak ada perubahan untuk disimpan.");
    }

    navigate(goBackPath);
  };

  if (!isJadwalIdValid || (!!error && !loading && !initialData)) {
    return (
      <div className="mx-auto max-w-3xl px-4 py-10">
        <div className="rounded-2xl border border-rose-200 bg-rose-50 p-8 text-center">
          <AlertCircle size={36} className="mx-auto text-rose-600" />
          <h1 className="mt-4 text-lg font-semibold text-rose-700">Gagal memuat data ujian</h1>
          <p className="mt-2 text-sm text-rose-600">
            {!isJadwalIdValid ? "ID jadwal ujian tidak valid." : error}
          </p>
          <Link
            to={goBackPath}
            className="mt-6 inline-flex items-center gap-2 text-sm font-medium text-rose-700 underline"
          >
            <ArrowLeft size={16} />
            Kembali ke jadwal ujian
          </Link>
        </div>
      </div>
    );
  }

  return (
    <BuatUjianForm
      mode="edit"
      initialValues={initialData?.values}
      initialSelectedMapelId={initialData?.selected_mapel_id ?? 0}
      loadingInitialData={loading}
      onSubmit={handleSubmit}
    />
  );
};

export default EditUjian;
