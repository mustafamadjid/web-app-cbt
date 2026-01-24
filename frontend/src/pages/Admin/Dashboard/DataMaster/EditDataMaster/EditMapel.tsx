import { useEffect, useMemo, useState } from "react";
import { useNavigate, useParams } from "react-router";

import EditMapelForm from "@/layouts/Form/Admin/DataMaster/EditMapelForm";
import type { MataPelajaranFormValues } from "@/types/DataMaster/MataPelajaran";
import { getMataPelajaranById, updateMataPelajaran } from "@/services/Api/features-api/DataMaster/mapel.service";
import { ApiError } from "@/services/Api/api";
import { paths } from "@/routes/paths";

const buildInitialValues = (): MataPelajaranFormValues => ({
  kelasId: "",
  kodeMapel: "",
  namaMapel: "",
  deskripsiMapel: "",
});

const EditMapel = () => {
  const { id } = useParams();
  const navigate = useNavigate();
  const [initialValues, setInitialValues] =
    useState<MataPelajaranFormValues>(buildInitialValues());
  const [loading, setLoading] = useState<boolean>(true);
  const [submitting, setSubmitting] = useState<boolean>(false);

  const mapelId = useMemo(() => Number(id), [id]);

  useEffect(() => {
    let active = true;
    const loadMapel = async () => {
      if (!id || Number.isNaN(mapelId)) {
        setLoading(false);
        return;
      }

      try {
        const data = await getMataPelajaranById(mapelId);
        if (!active || !data) return;
        setInitialValues(data);
      } finally {
        if (active) setLoading(false);
      }
    };

    loadMapel();

    return () => {
      active = false;
    };
  }, [id, mapelId]);

  const handleSubmit = async (values: MataPelajaranFormValues) => {
    if (!id || Number.isNaN(mapelId)) {
      throw new ApiError("ID mata pelajaran tidak ditemukan.");
    }

    setSubmitting(true);
    try {
      await updateMataPelajaran(mapelId, values);
      navigate(paths.dashboard.data_master_mapel);
    } finally {
      setSubmitting(false);
    }
  };

  return (
    <EditMapelForm
      initialValues={initialValues}
      onSubmit={handleSubmit}
      loading={loading}
      submitting={submitting}
    />
  );
};

export default EditMapel;
