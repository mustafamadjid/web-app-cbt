import { useNavigate } from "react-router";
import toast from "react-hot-toast";

import { useAuth } from "@/contexts/AuthContext";
import { buildCreatePenjadwalanUjianPayload } from "@/helper/FormData/buildCreatePenjadwalanUjianPayload";
import BuatUjianForm from "@/layouts/Form/Admin/Ujian/BuatUjianForm";
import { paths } from "@/routes/paths";
import { ApiError } from "@/services/Api/api";
import { createJadwalUjian } from "@/services/Api/features-api/Ujian/jadwalujian.service";
import type { BuatUjianFormValues } from "@/types/Ujian/BuatUjian";

const BuatUjian = () => {
  const navigate = useNavigate();
  const { user } = useAuth();

  const handleSubmit = async (values: BuatUjianFormValues) => {
    if (!user?.id_pengguna || user.id_pengguna <= 0) {
      throw new ApiError("Akun login tidak valid untuk membuat ujian.");
    }

    const payload = buildCreatePenjadwalanUjianPayload({
      values,
      idGuru: user.id_pengguna,
    });

    await createJadwalUjian(payload);
    toast.success("Jadwal ujian berhasil dibuat.");

    navigate(
      user.role === "GURU"
        ? paths.dashboard.jadwal_ujian_guru
        : paths.dashboard.jadwal_ujian,
    );
  };

  return <BuatUjianForm mode="create" onSubmit={handleSubmit} />;
};

export default BuatUjian;
