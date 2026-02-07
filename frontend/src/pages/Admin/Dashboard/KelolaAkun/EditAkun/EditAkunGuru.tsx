import { useEffect, useMemo, useState } from "react";
import { useNavigate, useParams } from "react-router";

import EditAkunGuruForm from "@/layouts/Form/Admin/KelolaAkun/EditAkunGuruForm";
import type { TeacherRegisterFormValues } from "@/types/KelolaAkun/AkunGuru";
// import { getGuruById, updateGuru } from "@/services/Api/features-api/KelolaAkun/akunguru.service";
import { ApiError } from "@/services/Api/api";
import { paths } from "@/routes/paths";

// const buildInitialValues = (): TeacherRegisterFormValues => ({
//   role: "GURU",
//   namaLengkap: "",
//   email: "",
//   username: "",
//   password: "",
//   noHp: "",
//   jenisKelamin: "LAKI_LAKI",
//   statusAkun: "aktif",
//   nip: "",
//   jabatan: "",
//   bidangStudi: "",
//   fotoProfil: null,
// });

// const EditAkunGuru = () => {
//   const { id } = useParams();
//   const navigate = useNavigate();
//   const [initialValues, setInitialValues] =
//     useState<TeacherRegisterFormValues>(buildInitialValues());
//   const [fotoUrl, setFotoUrl] = useState<string>("");
//   const [loading, setLoading] = useState<boolean>(true);
//   const [submitting, setSubmitting] = useState<boolean>(false);

//   const guruId = useMemo(() => Number(id), [id]);

//   useEffect(() => {
//     let active = true;
//     const loadGuru = async () => {
//       if (!id || Number.isNaN(guruId)) {
//         setLoading(false);
//         return;
//       }

//       try {
//         const data = await getGuruById(guruId);
//         if (!active || !data) return;

//         setInitialValues({
//           role: data.role ?? "GURU",
//           namaLengkap: data.namaLengkap,
//           email: data.email,
//           username: data.username,
//           password: "",
//           noHp: data.noHp,
//           jenisKelamin: data.jenisKelamin,
//           statusAkun: data.statusAkun,
//           nip: data.nip,
//           jabatan: data.jabatan,
//           bidangStudi: data.bidangStudi,
//           fotoProfil: null,
//         });
//         setFotoUrl(data.urlGambarProfil ?? "");
//       } finally {
//         if (active) setLoading(false);
//       }
//     };

//     loadGuru();

//     return () => {
//       active = false;
//     };
//   }, [guruId, id]);

//   const handleSubmit = async (values: TeacherRegisterFormValues) => {
//     if (!id || Number.isNaN(guruId)) {
//       throw new ApiError("ID guru tidak ditemukan.");
//     }

//     setSubmitting(true);
//     try {
//       await updateGuru(guruId, values);
//       navigate(paths.dashboard.kelola_akun_guru);
//     } finally {
//       setSubmitting(false);
//     }
//   };

//   return (
//     <EditAkunGuruForm
//       initialValues={initialValues}
//       initialFotoUrl={fotoUrl}
//       onSubmit={handleSubmit}
//       loading={loading}
//       submitting={submitting}
//     />
//   );
// };

const EditAkunGuru = () => <div></div>;

export default EditAkunGuru;
