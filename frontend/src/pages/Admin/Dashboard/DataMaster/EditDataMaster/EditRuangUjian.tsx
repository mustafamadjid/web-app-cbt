import { useEffect, useMemo, useState } from "react";
import { useNavigate, useParams } from "react-router";

import EditRuangForm from "@/layouts/Form/Admin/DataMaster/EditRuangForm";
import { paths } from "@/routes/paths";
import { ApiError } from "@/services/Api/api";
import {
  getRuangUjianById,
  updateRuangUjianPartial,
} from "@/services/Api/features-api/DataMaster/ruang-ujian.service";
import type { RuangUjianFormValues } from "@/types/DataMaster/RuangUjian";

const buildInitialValues = (): RuangUjianFormValues => ({
  kode_ruang: "",
  nama_ruangan: "",
});

const EditRuangUjian = () => {
  const { id } = useParams();
  const navigate = useNavigate();
  const [initialValues, setInitialValues] =
    useState<RuangUjianFormValues>(buildInitialValues());
  const [loading, setLoading] = useState<boolean>(true);
  const [submitting, setSubmitting] = useState<boolean>(false);

  const ruangId = useMemo(() => Number(id), [id]);

  useEffect(() => {
    let active = true;

    const loadRuang = async () => {
      if (!id || Number.isNaN(ruangId)) {
        setLoading(false);
        return;
      }

      try {
        const data = await getRuangUjianById(ruangId);
        if (!active) return;
        setInitialValues(data);
      } finally {
        if (active) setLoading(false);
      }
    };

    loadRuang();

    return () => {
      active = false;
    };
  }, [id, ruangId]);

  const handleSubmit = async (values: RuangUjianFormValues) => {
    if (!id || Number.isNaN(ruangId)) {
      throw new ApiError("ID ruangan tidak ditemukan.");
    }

    setSubmitting(true);
    try {
      await updateRuangUjianPartial(ruangId, values, initialValues);
      navigate(paths.dashboard.data_master_ruang);
    } finally {
      setSubmitting(false);
    }
  };

  return (
    <EditRuangForm
      initialValues={initialValues}
      onSubmit={handleSubmit}
      loading={loading}
      submitting={submitting}
    />
  );
};

export default EditRuangUjian;
