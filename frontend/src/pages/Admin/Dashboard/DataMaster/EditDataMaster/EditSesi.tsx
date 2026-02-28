import { useMemo, useState } from "react";
import { useNavigate, useParams } from "react-router";

import EditSesiForm from "@/layouts/Form/Admin/DataMaster/EditSesiForm";
import type { SesiFormValues } from "@/types/DataMaster/Sesi";
import {
  useGetSesiById,
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
  const [submitting, setSubmitting] = useState<boolean>(false);

  const sesiId = useMemo(() => Number(id), [id]);
  const isSesiIdValid = Boolean(id) && !Number.isNaN(sesiId);

  const { data: fetchedInitialValues, loading } = useGetSesiById(
    isSesiIdValid ? sesiId : -1,
  );

  const initialValues = useMemo<SesiFormValues>(
    () => fetchedInitialValues ?? buildInitialValues(),
    [fetchedInitialValues],
  );

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
      loading={isSesiIdValid ? loading : false}
      submitting={submitting}
    />
  );
};

export default EditSesi;
