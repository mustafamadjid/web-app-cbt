import { useMemo, useState } from "react";
import { useNavigate, useParams } from "react-router";

import EditRuangForm from "@/layouts/Form/Admin/DataMaster/EditRuangForm";
import { paths } from "@/routes/paths";
import { ApiError } from "@/services/Api/api";
import {
  useGetRuangUjianById,
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
  const [submitting, setSubmitting] = useState<boolean>(false);

  const ruangId = useMemo(() => Number(id), [id]);
  const isRuangIdValid = Boolean(id) && !Number.isNaN(ruangId);

  const { data: fetchedInitialValues, loading } = useGetRuangUjianById(
    isRuangIdValid ? ruangId : -1,
  );

  const initialValues = useMemo<RuangUjianFormValues>(
    () => fetchedInitialValues ?? buildInitialValues(),
    [fetchedInitialValues],
  );

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
      loading={isRuangIdValid ? loading : false}
      submitting={submitting}
    />
  );
};

export default EditRuangUjian;
