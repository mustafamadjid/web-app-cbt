import { useEffect, useMemo, useState } from "react";
import { useNavigate, useParams } from "react-router";

import EditAkunSiswaForm from "@/layouts/Form/Admin/KelolaAkun/EditAkunSiswaForm";
import type { StudentRegisterFormValues } from "@/types/KelolaAkun/AkunSiswa";
import { getSiswaById, updateSiswa } from "@/services/Api/features-api/KelolaAkun/akunsiswa.service";
import { ApiError } from "@/services/Api/api";
import { paths } from "@/routes/paths";

const buildInitialValues = (): StudentRegisterFormValues => ({
  role: "SISWA",
  namaLengkap: "",
  username: "",
  password: "",
  jenisKelamin: "LAKI_LAKI",
  email: "",
  noHp: "",
  noAbsen: 0,
  angkatan: 0,
  tempatLahir: "",
  tanggalLahir: "",
  id_tingkat_kelas: "",
  id_nama_kelas: "",
  fotoProfil: null,
  statusAkun: "aktif",
});

const EditAkunSiswa = () => {
  const { id } = useParams();
  const navigate = useNavigate();
  const [initialValues, setInitialValues] =
    useState<StudentRegisterFormValues>(buildInitialValues());
  const [fotoUrl, setFotoUrl] = useState<string>("");
  const [loading, setLoading] = useState<boolean>(true);
  const [submitting, setSubmitting] = useState<boolean>(false);

  const siswaId = useMemo(() => Number(id), [id]);

  useEffect(() => {
    let active = true;
    const loadSiswa = async () => {
      if (!id || Number.isNaN(siswaId)) {
        setLoading(false);
        return;
      }

      try {
        const data = await getSiswaById(siswaId);
        if (!active || !data) return;

        setInitialValues({
          role: "SISWA",
          namaLengkap: data.namaLengkap,
          username: data.username,
          password: "",
          jenisKelamin: data.jenisKelamin,
          email: data.email ?? "",
          noHp: data.noHp ?? "",
          noAbsen: data.noAbsen,
          angkatan: data.angkatan,
          tempatLahir: data.tempatLahir,
          tanggalLahir: data.tanggalLahir,
          id_tingkat_kelas: data.id_tingkat_kelas,
          id_nama_kelas: data.id_nama_kelas,
          fotoProfil: null,
          statusAkun: data.statusAkun,
        });
        setFotoUrl(data.urlGambarProfil ?? "");
      } finally {
        if (active) setLoading(false);
      }
    };

    loadSiswa();

    return () => {
      active = false;
    };
  }, [id, siswaId]);

  const handleSubmit = async (values: StudentRegisterFormValues) => {
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

  return (
    <EditAkunSiswaForm
      initialValues={initialValues}
      initialFotoUrl={fotoUrl}
      onSubmit={handleSubmit}
      loading={loading}
      submitting={submitting}
    />
  );
};

export default EditAkunSiswa;
