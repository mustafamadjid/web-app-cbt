import { useEffect, useMemo, useState } from "react";
import { useNavigate, useParams } from "react-router";

import EditSesiForm from "@/layouts/Form/Admin/DataMaster/EditSesiForm";
import type { SesiFormValues } from "@/types/DataMaster/Sesi";
import {
  getSesiById,
  updateSesiPartial,
} from "@/services/Api/features-api/DataMaster/sesi.service";
import { ApiError } from "@/services/Api/api";
import { paths } from "@/routes/paths";

const buildInitialValues = (): SesiFormValues => ({
  kode_sesi: "",
  nama_sesi: "",
});

const EditSesi = () => {
  const { id } = useParams();
  const navigate = useNavigate();
  const [initialValues, setInitialValues] =
    useState<SesiFormValues>(buildInitialValues());
  const [loading, setLoading] = useState<boolean>(true);
  const [submitting, setSubmitting] = useState<boolean>(false);

  const sesiId = useMemo(() => Number(id), [id]);

  useEffect(() => {
    let active = true;
    const loadSesi = async () => {
      if (!id || Number.isNaN(sesiId)) {
        setLoading(false);
        return;
      }

      try {
        const data = await getSesiById(sesiId);
        if (!active || !data) return;
        setInitialValues(data);
      } finally {
        if (active) setLoading(false);
      }
    };

    loadSesi();

    return () => {
      active = false;
    };
  }, [id, sesiId]);

  const handleSubmit = async (values: SesiFormValues) => {
    if (!id || Number.isNaN(sesiId)) {
      throw new ApiError("ID sesi tidak ditemukan.");
    }

    setSubmitting(true);
    try {
      await updateSesiPartial(sesiId, values, initialValues);
      navigate(paths.dashboard.data_master_sesi);
    } finally {
      setSubmitting(false);
    }
  };

  return (
    <EditSesiForm
      key={`${initialValues.kode_sesi}-${initialValues.nama_sesi}`}
      initialValues={initialValues}
      onSubmit={handleSubmit}
      loading={loading}
      submitting={submitting}
    />
  );
};

export default EditSesi;
