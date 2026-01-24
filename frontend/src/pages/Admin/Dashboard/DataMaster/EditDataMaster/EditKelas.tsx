import { useEffect, useMemo, useState } from "react";
import { useNavigate, useParams } from "react-router";

import EditKelasForm from "@/layouts/Form/Admin/DataMaster/EditKelasForm";
import type { KelasFormValues } from "@/types/DataMaster/Kelas";
import { getKelasById, updateKelas } from "@/services/Api/features-api/DataMaster/kelas.service";
import { ApiError } from "@/services/Api/api";
import { paths } from "@/routes/paths";

const buildInitialValues = (): KelasFormValues => ({
  tingkat_kelas: "",
  nama_kelas: "",
});

const EditKelas = () => {
  const { id } = useParams();
  const navigate = useNavigate();
  const [initialValues, setInitialValues] =
    useState<KelasFormValues>(buildInitialValues());
  const [loading, setLoading] = useState<boolean>(true);
  const [submitting, setSubmitting] = useState<boolean>(false);

  const kelasId = useMemo(() => Number(id), [id]);

  useEffect(() => {
    let active = true;
    const loadKelas = async () => {
      if (!id || Number.isNaN(kelasId)) {
        setLoading(false);
        return;
      }

      try {
        const data = await getKelasById(kelasId);
        if (!active || !data) return;
        setInitialValues(data);
      } finally {
        if (active) setLoading(false);
      }
    };

    loadKelas();

    return () => {
      active = false;
    };
  }, [id, kelasId]);

  const handleSubmit = async (values: KelasFormValues) => {
    if (!id || Number.isNaN(kelasId)) {
      throw new ApiError("ID kelas tidak ditemukan.");
    }

    setSubmitting(true);
    try {
      await updateKelas(kelasId, values);
      navigate(paths.dashboard.data_master_kelas);
    } finally {
      setSubmitting(false);
    }
  };

  return (
    <EditKelasForm
      initialValues={initialValues}
      onSubmit={handleSubmit}
      loading={loading}
      submitting={submitting}
    />
  );
};

export default EditKelas;
