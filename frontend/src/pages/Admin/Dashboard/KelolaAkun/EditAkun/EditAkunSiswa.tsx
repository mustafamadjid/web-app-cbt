import { useEffect, useMemo, useState } from "react";
import { useNavigate, useParams } from "react-router";

import EditAkunSiswaForm from "@/layouts/Form/Admin/KelolaAkun/EditAkunSiswaForm";
import type { StudentUpdateFormValues } from "@/types/KelolaAkun/AkunSiswa";
import { getSiswaById, updateSiswa } from "@/services/Api/features-api/KelolaAkun/akunsiswa.service";
import { ApiError } from "@/services/Api/api";
import { paths } from "@/routes/paths";
import { GetDataKelasFull } from "@/services/Api/features-api/DataMaster/kelas.service";

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
  id_tingkat_kelas: "",
  id_nama_kelas: "",
  foto_profil: null,
  status_akun: "AKTIF",
});

const EditAkunSiswa = () => {
  const { id } = useParams();
  const navigate = useNavigate();
  const [initialValues, setInitialValues] =
    useState<StudentUpdateFormValues>(buildInitialValues());
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
        const [data, kelas] = await Promise.all([
          getSiswaById(siswaId),
          GetDataKelasFull(),
        ]);
        if (!active || !data) return;

        const tingkatMatch = kelas.item_tingkat_kelas.find(
          (item) => item.tingkat_kelas === data.tingkat_kelas,
        );
        const namaMatch = kelas.item_nama_kelas.find((item) => {
          if (item.nama_kelas !== data.nama_kelas) return false;
          if (tingkatMatch) return item.id_tingkat_kelas === tingkatMatch.id_tingkat_kelas;
          return true;
        });

        setInitialValues({
          id_pengguna: data.id_pengguna,
          role: data.role ?? "SISWA",
          nama_lengkap: data.nama_lengkap,
          username: data.username,
          jenis_kelamin: data.jenis_kelamin,
          email: data.email ?? "",
          no_hp: data.no_hp ?? "",
          nisn: data.nisn ?? "",
          no_absen: data.no_absen,
          angkatan: data.angkatan,
          tempat_lahir: data.tempat_lahir,
          tanggal_lahir: data.tanggal_lahir,
          id_tingkat_kelas: tingkatMatch?.id_tingkat_kelas ?? "",
          id_nama_kelas: namaMatch ? String(namaMatch.id_nama_kelas) : "",
          foto_profil: null,
          status_akun: data.status_akun,
        });
        setFotoUrl(
          data.foto_profil
            ? `${import.meta.env.VITE_API_URL}${data.foto_profil}`
            : "",
        );
      } finally {
        if (active) setLoading(false);
      }
    };

    loadSiswa();

    return () => {
      active = false;
    };
  }, [id, siswaId]);

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
