import { useMemo, useState } from "react";
import { useNavigate, useParams } from "react-router";

import EditAkunSiswaForm from "@/layouts/Form/Admin/KelolaAkun/EditAkunSiswaForm";
import type { StudentUpdateFormValues } from "@/types/KelolaAkun/AkunSiswa";
import type { ResetPasswordFormValues } from "@/types/KelolaAkun/ResetPassword";
import {
  useGetSiswaById,
  updateSiswa,
} from "@/services/Api/features-api/KelolaAkun/akunsiswa.service";
import { resetPasswordPengguna } from "@/services/Api/features-api/KelolaAkun/akun.service";
import { ApiError } from "@/services/Api/api";
import { paths } from "@/routes/paths";
import { useGetDataKelasFull } from "@/services/Api/features-api/DataMaster/kelas.service";
import toast from "react-hot-toast";
import { resolveImageUrl } from "@/helper/MediaUrl/resolveMediaUrl";

const buildInitialValues = (): StudentUpdateFormValues => ({
  id_pengguna: 0,
  role: "SISWA",
  nama_lengkap: "",
  username: "",
  jenis_kelamin: "LAKI_LAKI",
  email: "",
  no_hp: "",
  nisn: "",
  no_absen: 0,
  angkatan: 0,
  tempat_lahir: "",
  tanggal_lahir: "",
  id_nama_kelas: "",
  foto_profil: null,
  status_akun: "AKTIF",
});

const EditAkunSiswa = () => {
  const { id } = useParams();
  const navigate = useNavigate();
  const [submitting, setSubmitting] = useState<boolean>(false);

  const siswaId = useMemo(() => Number(id), [id]);
  const isSiswaIdValid = Boolean(id) && !Number.isNaN(siswaId);

  const {
    data: siswaData,
    loading: loadingSiswa,
  } = useGetSiswaById(isSiswaIdValid ? siswaId : -1);
  const {
    data: kelasData,
    loading: loadingKelas,
  } = useGetDataKelasFull();

  const initialValues = useMemo<StudentUpdateFormValues>(() => {
    if (!siswaData) return buildInitialValues();

    const namaMatch = kelasData?.item_nama_kelas.find(
      (item) => item.nama_kelas === siswaData.nama_kelas,
    );

    return {
      id_pengguna: siswaData.id_pengguna,
      role: siswaData.role ?? "SISWA",
      nama_lengkap: siswaData.nama_lengkap,
      username: siswaData.username,
      jenis_kelamin: siswaData.jenis_kelamin,
      email: siswaData.email ?? "",
      no_hp: siswaData.no_hp ?? "",
      nisn: siswaData.nisn ?? "",
      no_absen: siswaData.no_absen,
      angkatan: siswaData.angkatan,
      tempat_lahir: siswaData.tempat_lahir,
      tanggal_lahir: siswaData.tanggal_lahir,
      id_nama_kelas: namaMatch ? String(namaMatch.id_nama_kelas) : "",
      foto_profil: null,
      status_akun: siswaData.status_akun,
    };
  }, [kelasData?.item_nama_kelas, siswaData]);

  const fotoUrl = useMemo(
    () => resolveImageUrl(siswaData?.foto_profil),
    [siswaData?.foto_profil],
  );

  const loading = isSiswaIdValid ? loadingSiswa || loadingKelas : false;

  const handleSubmit = async (values: StudentUpdateFormValues) => {
    if (!id || Number.isNaN(siswaId)) {
      throw new ApiError("ID siswa tidak ditemukan.");
    }

    setSubmitting(true);
    try {
      await updateSiswa(siswaId, values);
      navigate(paths.dashboard.kelola_akun_siswa);
    } finally {
      setSubmitting(false);
    }
  };

  const handleSubmitResetPassword = async (values: ResetPasswordFormValues) => {
    if (!id || Number.isNaN(siswaId)) {
      throw new ApiError("ID siswa tidak ditemukan.");
    }

    setSubmitting(true);
    try {
      await resetPasswordPengguna(siswaId, { password: values.password });
      toast.success("Password akun siswa berhasil diperbarui.");
    } finally {
      setSubmitting(false);
    }
  };

  return (
    <EditAkunSiswaForm
      initialValues={initialValues}
      initialFotoUrl={fotoUrl}
      onSubmit={handleSubmit}
      onSubmitResetPassword={handleSubmitResetPassword}
      loading={loading}
      submitting={submitting}
    />
  );
};

export default EditAkunSiswa;
