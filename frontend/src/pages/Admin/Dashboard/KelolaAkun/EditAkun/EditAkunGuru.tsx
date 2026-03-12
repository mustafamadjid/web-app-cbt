import { useMemo, useState } from "react";
import { useNavigate, useParams } from "react-router";

import EditAkunGuruForm from "@/layouts/Form/Admin/KelolaAkun/EditAkunGuruForm";
import { resolveImageUrl } from "@/helper/MediaUrl/resolveMediaUrl";
import type { TeacherUpdateFormValues } from "@/types/KelolaAkun/AkunGuru";
import type { ResetPasswordFormValues } from "@/types/KelolaAkun/ResetPassword";
import {
  useGetGuruById,
  updateGuru,
} from "@/services/Api/features-api/KelolaAkun/akunguru.service";
import { resetPasswordPengguna } from "@/services/Api/features-api/KelolaAkun/akun.service";
import { ApiError } from "@/services/Api/api";
import { paths } from "@/routes/paths";
import toast from "react-hot-toast";

const buildInitialValues = (): TeacherUpdateFormValues => ({
  id_pengguna: 0,
  role: "GURU",
  nama_lengkap: "",
  email: "",
  username: "",
  no_hp: "",
  jenis_kelamin: "LAKI_LAKI",
  status_akun: "AKTIF",
  nip: "",
  jabatan: "",
  bidang_studi: "",
  foto_profil: null,
});

const EditAkunGuru = () => {
  const { id } = useParams();
  const navigate = useNavigate();
  const [submitting, setSubmitting] = useState<boolean>(false);

  const guruId = useMemo(() => Number(id), [id]);
  const isGuruIdValid = Boolean(id) && !Number.isNaN(guruId);

  const { data: guruData, loading } = useGetGuruById(
    isGuruIdValid ? guruId : -1,
  );

  const initialValues = useMemo<TeacherUpdateFormValues>(() => {
    if (!guruData) return buildInitialValues();

    return {
      id_pengguna: guruData.id_pengguna,
      role: guruData.role ?? "GURU",
      nama_lengkap: guruData.nama_lengkap,
      email: guruData.email ?? "",
      username: guruData.username,
      no_hp: guruData.no_hp ?? "",
      jenis_kelamin: guruData.jenis_kelamin,
      status_akun: guruData.status_akun,
      nip: guruData.nip,
      jabatan: guruData.jabatan,
      bidang_studi: guruData.bidang_studi,
      foto_profil: null,
    };
  }, [guruData]);

  const fotoUrl = useMemo(
    () => resolveImageUrl(guruData?.foto_profil),
    [guruData?.foto_profil],
  );

  const handleSubmit = async (values: TeacherUpdateFormValues) => {
    if (!id || Number.isNaN(guruId)) {
      throw new ApiError("ID guru tidak ditemukan.");
    }

    setSubmitting(true);
    try {
      await updateGuru(guruId, values);
      toast.success("Akun guru berhasil diperbarui.");
      navigate(paths.dashboard.kelola_akun_guru);
    } finally {
      setSubmitting(false);
    }
  };

  const handleSubmitResetPassword = async (values: ResetPasswordFormValues) => {
    if (!id || Number.isNaN(guruId)) {
      throw new ApiError("ID guru tidak ditemukan.");
    }

    setSubmitting(true);
    try {
      await resetPasswordPengguna(guruId, { password: values.password });
      toast.success("Password akun guru berhasil diperbarui.");
    } finally {
      setSubmitting(false);
    }
  };

  return (
    <EditAkunGuruForm
      key={`${initialValues.id_pengguna}-${initialValues.username}-${fotoUrl}`}
      initialValues={initialValues}
      initialFotoUrl={fotoUrl}
      onSubmit={handleSubmit}
      onSubmitResetPassword={handleSubmitResetPassword}
      loading={isGuruIdValid ? loading : false}
      submitting={submitting}
    />
  );
};

// const EditAkunGuru = () => <div></div>;

export default EditAkunGuru;
