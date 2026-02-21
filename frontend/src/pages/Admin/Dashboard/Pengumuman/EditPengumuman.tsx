import { useEffect, useMemo, useState } from "react";
import { useNavigate, useParams } from "react-router";

import EditPengumumanForm from "@/layouts/Form/Admin/Pengumuman/EditPengumumanForm";
import { useAuth } from "@/contexts/AuthContext";
import { paths } from "@/routes/paths";
import { ApiError } from "@/services/Api/api";
import {
  getPengumumanById,
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

  const [initialValues, setInitialValues] =
    useState<PengumumanFormValues>(buildInitialValues());
  const [dokumenLama, setDokumenLama] = useState("");
  const [loading, setLoading] = useState(true);
  const [submitting, setSubmitting] = useState(false);

  const pengumumanId = useMemo(() => Number(id), [id]);
  const goBackPath =
    user?.role === "GURU"
      ? paths.dashboard.pengumuman_guru
      : paths.dashboard.pengumuman_admin;

  useEffect(() => {
    let active = true;
    const fetchDetail = async () => {
      if (!id || Number.isNaN(pengumumanId)) {
        setLoading(false);
        return;
      }

      try {
        const data = await getPengumumanById(pengumumanId);
        if (!active) return;

        setInitialValues({
          judul_pengumuman: data.judul_pengumuman,
          isi_pengumuman: data.isi_pengumuman,
          tanggal_rilis_pengumuman: data.tanggal_rilis_pengumuman,
          tanggal_selesai_pengumuman: data.tanggal_selesai_pengumuman,
          dokumen_pengumuman: null,
        });
        setDokumenLama(data.dokumen_pengumuman ?? "");
      } finally {
        if (active) {
          setLoading(false);
        }
      }
    };

    fetchDetail();
    return () => {
      active = false;
    };
  }, [id, pengumumanId]);

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
      loading={loading}
      submitting={submitting}
    />
  );
};

export default EditPengumuman;
