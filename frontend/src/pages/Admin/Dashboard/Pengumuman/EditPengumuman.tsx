import { useMemo, useState } from "react";
import { useNavigate, useParams } from "react-router";

import EditPengumumanForm from "@/layouts/Form/Admin/Pengumuman/EditPengumumanForm";
import { useAuth } from "@/contexts/AuthContext";
import { paths } from "@/routes/paths";
import { ApiError } from "@/services/Api/api";
import {
  useGetPengumumanById,
  updatePengumumanPartial,
} from "@/services/Api/features-api/pengumuman/pengumuman.service";
import type { PengumumanFormValues } from "@/types/Widget/Pengumuman";

const buildInitialValues = (): PengumumanFormValues => ({
  judul_pengumuman: "",
  isi_pengumuman: "",
  tanggal_rilis_pengumuman: "",
  tanggal_selesai_pengumuman: "",
  dokumen_pengumuman: null,
});

const EditPengumuman = () => {
  const { id } = useParams();
  const navigate = useNavigate();
  const { user } = useAuth();

  const [submitting, setSubmitting] = useState(false);

  const pengumumanId = useMemo(() => Number(id), [id]);
  const isPengumumanIdValid = Boolean(id) && !Number.isNaN(pengumumanId);
  const goBackPath =
    user?.role === "GURU"
      ? paths.dashboard.pengumuman_guru
      : paths.dashboard.pengumuman_admin;

  const { data: pengumumanDetail, loading } = useGetPengumumanById(
    isPengumumanIdValid ? pengumumanId : -1,
  );

  const initialValues = useMemo<PengumumanFormValues>(() => {
    if (!pengumumanDetail) return buildInitialValues();

    return {
      judul_pengumuman: pengumumanDetail.judul_pengumuman,
      isi_pengumuman: pengumumanDetail.isi_pengumuman,
      tanggal_rilis_pengumuman: pengumumanDetail.tanggal_rilis_pengumuman,
      tanggal_selesai_pengumuman: pengumumanDetail.tanggal_selesai_pengumuman,
      dokumen_pengumuman: null,
    };
  }, [pengumumanDetail]);

  const dokumenLama = useMemo(
    () => pengumumanDetail?.dokumen_pengumuman ?? "",
    [pengumumanDetail],
  );

  const handleSubmit = async (values: PengumumanFormValues) => {
    if (!id || Number.isNaN(pengumumanId)) {
      throw new ApiError("ID pengumuman tidak ditemukan.");
    }

    setSubmitting(true);
    try {
      await updatePengumumanPartial(pengumumanId, values, initialValues);
      navigate(goBackPath);
    } finally {
      setSubmitting(false);
    }
  };

  return (
    <EditPengumumanForm
      initialValues={initialValues}
      dokumenLama={dokumenLama}
      onSubmit={handleSubmit}
      loading={isPengumumanIdValid ? loading : false}
      submitting={submitting}
    />
  );
};

export default EditPengumuman;
