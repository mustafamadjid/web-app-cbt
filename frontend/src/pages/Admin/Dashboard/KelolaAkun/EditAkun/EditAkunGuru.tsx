import { useEffect, useMemo, useState } from "react";
import { useNavigate, useParams } from "react-router";

import EditAkunGuruForm from "@/layouts/Form/Admin/KelolaAkun/EditAkunGuruForm";
import type { DataUpdateGuru, TeacherRegisterFormValues } from "@/types/KelolaAkun/AkunGuru";
import { getGuruById, updateGuru } from "@/services/Api/features-api/KelolaAkun/akunguru.service";
import { ApiError } from "@/services/Api/api";
import { paths } from "@/routes/paths";

const buildInitialValues = (): DataUpdateGuru => ({
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
  foto_profil: "",
});

const EditAkunGuru = () => {
  const { id } = useParams();
  const navigate = useNavigate();
  const [initialValues, setInitialValues] =
    useState<DataUpdateGuru>(buildInitialValues());
  const [fotoUrl, setFotoUrl] = useState<string>("");
  const [loading, setLoading] = useState<boolean>(true);
  const [submitting, setSubmitting] = useState<boolean>(false);

  const guruId = useMemo(() => Number(id), [id]);

  useEffect(() => {
    let active = true;
    const loadGuru = async () => {
      if (!id || Number.isNaN(guruId)) {
        setLoading(false);
        return;
      }

      try {
        const data = await getGuruById(guruId);
        if (!active || !data) return;

        setInitialValues({
          id_pengguna: data.id_pengguna,
          role: data.role ?? "GURU",
          nama_lengkap: data.nama_lengkap,
          email: data.email,
          username: data.username,
          no_hp: data.no_hp,
          jenis_kelamin: data.jenis_kelamin,
          status_akun: data.status_akun,
          nip: data.nip,
          jabatan: data.jabatan,
          bidang_studi: data.bidang_studi,
          foto_profil: "",
        });
        setFotoUrl(data.foto_profil ?? "");
      } finally {
        if (active) setLoading(false);
      }
    };

    loadGuru();

    return () => {
      active = false;
    };
  }, [guruId, id]);

  const handleSubmit = async (values: TeacherRegisterFormValues) => {
    if (!id || Number.isNaN(guruId)) {
      throw new ApiError("ID guru tidak ditemukan.");
    }

    setSubmitting(true);
    try {
      await updateGuru(guruId, values);
      navigate(paths.dashboard.kelola_akun_guru);
    } finally {
      setSubmitting(false);
    }
  };

  return (
    <EditAkunGuruForm
      initialValues={initialValues}
      initialFotoUrl={fotoUrl}
      onSubmit={handleSubmit}
      loading={loading}
      submitting={submitting}
    />
  );
};

// const EditAkunGuru = () => <div></div>;

export default EditAkunGuru;
