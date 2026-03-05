import { useEffect, useMemo, useState } from "react";
import { useNavigate } from "react-router";
import toast from "react-hot-toast";

import DatePicker from "@/components/common/DateInput/DatePicker";
import TimePicker from "@/components/common/DateInput/TimePicker";
import InputField from "@/components/common/Input/InputField";
import { useAuth } from "@/contexts/AuthContext";
import { buildCreatePenjadwalanUjianPayload } from "@/helper/FormData/buildCreatePenjadwalanUjianPayload";
import { calculateDuration } from "@/helper/CalculateDuration/calculateDuration";
import { createSetField } from "@/helper/setField/setField";
import {
  createValidator,
  matchesPattern,
  minNumber,
  requiredString,
  requiredValue,
} from "@/helper/validate/validateForm";
import { paths } from "@/routes/paths";
import { ApiError } from "@/services/Api/api";
import { useGetBankSoal } from "@/services/Api/features-api/BankSoal/banksoal.service";
import { useGetDataKelasFull } from "@/services/Api/features-api/DataMaster/kelas.service";
import { useGetRuangUjian } from "@/services/Api/features-api/DataMaster/ruang-ujian.service";
import { useGetSesi } from "@/services/Api/features-api/DataMaster/sesi.service";
import { useGetAllGuru } from "@/services/Api/features-api/KelolaAkun/akunguru.service";
import { useGetListSiswa } from "@/services/Api/features-api/KelolaAkun/akunsiswa.service";
import { createJadwalUjian } from "@/services/Api/features-api/Ujian/jadwalujian.service";
import type { FullDataKelas, NamaKelas, TingkatKelas } from "@/types/DataMaster/Kelas";
import type { DataGuru } from "@/types/KelolaAkun/AkunGuru";
import type { DataAkunSiswa } from "@/types/KelolaAkun/AkunSiswa";
import type { BuatUjianFormValues } from "@/types/Ujian/BuatUjian";
import type { BankSoalItem } from "@/types/BankSoal/BankSoal";
import type { RuangUjianRow } from "@/types/DataMaster/RuangUjian";
import type { SesiRow } from "@/types/DataMaster/Sesi";

const initialValues: BuatUjianFormValues = {
  nama_ujian: "",
  deskripsi_ujian: "",
  id_kelas: 0,
  kelas_scope: "SEMUA",
  id_nama_kelas: 0,
  id_bank_soal: 0,
  tanggal_ujian: "",
  waktu_mulai: "",
  waktu_selesai: "",
  id_ruangan: 0,
  acak_soal: true,
  id_pengawas: 0,
  id_sesi: 0,
  token: "",
};

const sectionTitle = "text-sm font-semibold text-slate-800";
const helperText = "text-xs text-slate-500";

const selectBaseClass =
  "w-full rounded-lg border border-slate-200 bg-white px-3 py-2 text-sm outline-none transition focus:border-[#397e50] focus:ring-2 focus:ring-[#397e50] disabled:bg-slate-50 disabled:text-slate-500";

const timePatternRule = matchesPattern(
  /^([01]\d|2[0-3]):([0-5]\d)$/,
  "Gunakan format 24 jam HH:mm, contoh 07:30.",
);

const BuatUjianForm = () => {
  const navigate = useNavigate();
  const { user } = useAuth();
  const [values, setValues] = useState<BuatUjianFormValues>(initialValues);
  const [touched, setTouched] = useState<Record<string, boolean>>({});
  const [submitError, setSubmitError] = useState<string | null>(null);
  const [submitting, setSubmitting] = useState(false);

  const setField = createSetField(setValues);

  const onBlur = (name: keyof BuatUjianFormValues) => {
    setTouched((prev) => ({ ...prev, [name]: true }));
  };

  const {
    data: kelasData,
    loading: loadingKelas,
    error: kelasError,
  } = useGetDataKelasFull();
  const kelasFull: FullDataKelas = kelasData ?? {
    item_tingkat_kelas: [],
    item_nama_kelas: [],
  };
  const tingkatKelasOptions: TingkatKelas[] = kelasFull.item_tingkat_kelas;
  const namaKelasOptions: NamaKelas[] = useMemo(
    () =>
      kelasFull.item_nama_kelas.filter(
        (item) => item.id_tingkat_kelas === values.id_kelas,
      ),
    [kelasFull.item_nama_kelas, values.id_kelas],
  );

  const {
    data: bankSoalData,
    loading: loadingBankSoal,
    error: bankSoalError,
  } = useGetBankSoal(
    {
      id_kelas: values.id_kelas > 0 ? values.id_kelas : undefined,
      limit: 100,
      offset: 0,
    },
    values.id_kelas > 0,
  );
  const bankSoalOptions: BankSoalItem[] = bankSoalData ?? [];

  const {
    data: ruangData,
    loading: loadingRuang,
    error: ruangError,
  } = useGetRuangUjian();
  const ruangOptions: RuangUjianRow[] = ruangData ?? [];

  const {
    data: sesiData,
    loading: loadingSesi,
    error: sesiError,
  } = useGetSesi();
  const sesiOptions: SesiRow[] = sesiData ?? [];

  const {
    data: guruData,
    loading: loadingGuru,
    error: guruError,
  } = useGetAllGuru();
  const guruOptions: DataGuru[] = guruData ?? [];

  const siswaPreviewEnabled =
    values.id_kelas > 0 &&
    (values.kelas_scope === "SEMUA" || values.id_nama_kelas > 0);
  const {
    data: siswaData,
    loading: loadingSiswa,
    error: siswaError,
  } = useGetListSiswa(
    {
      idKelas: values.id_kelas > 0 ? values.id_kelas : undefined,
      idNamaKelas:
        values.kelas_scope === "SPESIFIK" && values.id_nama_kelas > 0
          ? values.id_nama_kelas
          : undefined,
      limit: 200,
      offset: 0,
    },
    siswaPreviewEnabled,
  );
  const siswaPreview: DataAkunSiswa[] = siswaData ?? [];

  const loadErrorMessage =
    kelasError ||
    bankSoalError ||
    ruangError ||
    sesiError ||
    guruError ||
    siswaError;

  useEffect(() => {
    setField("id_nama_kelas", 0);
    setField("id_bank_soal", 0);
  }, [values.id_kelas]);

  useEffect(() => {
    if (values.kelas_scope === "SEMUA" && values.id_nama_kelas !== 0) {
      setField("id_nama_kelas", 0);
    }
  }, [values.kelas_scope, values.id_nama_kelas]);

  useEffect(() => {
    if (
      values.id_bank_soal !== 0 &&
      !bankSoalOptions.some((item) => item.id_bank_soal === values.id_bank_soal)
    ) {
      setField("id_bank_soal", 0);
    }
  }, [bankSoalOptions, values.id_bank_soal]);

  const durasiMenit = useMemo(
    () => calculateDuration(values.waktu_mulai, values.waktu_selesai),
    [values.waktu_mulai, values.waktu_selesai],
  );

  const validate = createValidator<BuatUjianFormValues>({
    nama_ujian: [requiredString("Nama ujian wajib diisi.")],
    deskripsi_ujian: [requiredString("Deskripsi ujian wajib diisi.")],
    id_kelas: [minNumber(1, "Tingkat kelas wajib dipilih.")],
    kelas_scope: [requiredValue("Cakupan kelas wajib dipilih.")],
    id_nama_kelas: [
      (value, currentValues) =>
        currentValues.kelas_scope === "SPESIFIK" && value === 0
          ? "Nama kelas wajib dipilih."
          : null,
    ],
    id_bank_soal: [minNumber(1, "Bank soal wajib dipilih.")],
    tanggal_ujian: [requiredString("Tanggal ujian wajib diisi.")],
    waktu_mulai: [requiredString("Waktu mulai wajib diisi."), timePatternRule],
    waktu_selesai: [
      requiredString("Waktu selesai wajib diisi."),
      timePatternRule,
      (_, currentValues) =>
        currentValues.waktu_mulai &&
        currentValues.waktu_selesai &&
        calculateDuration(currentValues.waktu_mulai, currentValues.waktu_selesai) <= 0
          ? "Waktu selesai harus setelah waktu mulai."
          : null,
    ],
    id_ruangan: [minNumber(1, "Ruang ujian wajib dipilih.")],
    id_pengawas: [minNumber(1, "Guru pengawas wajib dipilih.")],
    id_sesi: [minNumber(1, "Sesi ujian wajib dipilih.")],
    token: [
      requiredString("Token ujian wajib diisi."),
      (value) =>
        value.trim().length > 30
          ? "Token ujian maksimal 30 karakter."
          : null,
    ],
  });

  const errors = validate(values);
  const hasError = (name: keyof BuatUjianFormValues) =>
    Boolean(errors[name]) && Boolean(touched[name]);

  const onSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    setSubmitError(null);

    const nextTouched: Record<string, boolean> = {};
    Object.keys(values).forEach((key) => {
      nextTouched[key] = true;
    });
    setTouched(nextTouched);

    const currentErrors = validate(values);
    if (Object.keys(currentErrors).length > 0) {
      setSubmitError("Periksa kembali input yang masih kosong atau belum valid.");
      return;
    }

    if (!user?.id_pengguna || user.id_pengguna <= 0) {
      setSubmitError("Akun login tidak valid untuk membuat ujian.");
      return;
    }

    setSubmitting(true);
    try {
      const payload = buildCreatePenjadwalanUjianPayload({
        values,
        idGuru: user.id_pengguna,
      });
      await createJadwalUjian(payload);
      toast.success("Jadwal ujian berhasil dibuat.");
      const nextPath =
        user.role === "GURU"
          ? paths.dashboard.jadwal_ujian_guru
          : paths.dashboard.jadwal_ujian;
      navigate(nextPath);
    } catch (error) {
      if (error instanceof ApiError) {
        setSubmitError(error.message);
      } else {
        setSubmitError("Terjadi kesalahan saat menyimpan ujian.");
      }
    } finally {
      setSubmitting(false);
    }
  };

  return (
    <div className="min-h-screen w-full py-8">
      <div className="mx-auto w-full max-w-6xl px-4">
        <div className="mb-6 rounded-xl border border-slate-200 bg-white p-5 shadow-sm">
          <h1 className="text-base font-semibold text-slate-900">Buat Ujian</h1>
          <p className="mt-1 text-sm text-slate-500">
            Lengkapi detail ujian dan jadwalnya sesuai data server.
          </p>
        </div>

        <form onSubmit={onSubmit} className="space-y-6">
          <div className="rounded-xl border border-slate-200 bg-white p-5 shadow-sm">
            <div className="mb-4">
              <h2 className={sectionTitle}>Informasi Ujian</h2>
              <p className={helperText}>Isi nama dan deskripsi ujian.</p>
            </div>

            <div className="grid grid-cols-1 gap-4 md:grid-cols-2">
              <div className="md:col-span-2">
                <InputField
                  id="nama_ujian"
                  label="Nama Ujian"
                  value={values.nama_ujian}
                  onChange={(value) => setField("nama_ujian", value)}
                  onBlur={() => onBlur("nama_ujian")}
                  placeholder="Contoh: Ujian Tengah Semester Matematika"
                  inputClassName={hasError("nama_ujian") ? "border-rose-300 ring-rose-100" : ""}
                />
                {hasError("nama_ujian") && (
                  <p className="mt-1 text-xs text-rose-600">{errors.nama_ujian}</p>
                )}
              </div>

              <div className="md:col-span-2">
                <label htmlFor="deskripsi_ujian" className="text-xs font-medium text-slate-600">
                  Deskripsi Ujian
                </label>
                <textarea
                  id="deskripsi_ujian"
                  className={`min-h-[100px] w-full rounded-lg border border-slate-200 bg-white px-3 py-2 text-sm outline-none transition focus:border-[#397e50] focus:ring-2 focus:ring-[#397e50] ${
                    hasError("deskripsi_ujian") ? "border-rose-300 ring-rose-100" : ""
                  }`}
                  value={values.deskripsi_ujian}
                  onChange={(event) => setField("deskripsi_ujian", event.target.value)}
                  onBlur={() => onBlur("deskripsi_ujian")}
                  placeholder="Jelaskan cakupan materi ujian."
                />
                {hasError("deskripsi_ujian") && (
                  <p className="mt-1 text-xs text-rose-600">{errors.deskripsi_ujian}</p>
                )}
              </div>
            </div>
          </div>

          <div className="rounded-xl border border-slate-200 bg-white p-5 shadow-sm">
            <div className="mb-4">
              <h2 className={sectionTitle}>Kelas & Bank Soal</h2>
              <p className={helperText}>
                Gunakan data kelas penuh, lalu pilih bank soal sesuai tingkat kelas.
              </p>
            </div>

            <div className="grid grid-cols-1 gap-4 md:grid-cols-2 lg:grid-cols-3">
              <div>
                <label htmlFor="id_kelas" className="text-xs font-medium text-slate-600">
                  Tingkat Kelas
                </label>
                <select
                  id="id_kelas"
                  value={values.id_kelas}
                  onChange={(event) => setField("id_kelas", Number(event.target.value))}
                  onBlur={() => onBlur("id_kelas")}
                  disabled={loadingKelas}
                  className={`${selectBaseClass} ${hasError("id_kelas") ? "border-rose-300 ring-rose-100" : ""}`}
                >
                  <option value={0}>
                    {loadingKelas ? "Memuat tingkat kelas..." : "Pilih tingkat kelas"}
                  </option>
                  {tingkatKelasOptions.map((item) => (
                    <option key={item.id_tingkat_kelas} value={item.id_tingkat_kelas}>
                      Kelas {item.tingkat_kelas}
                    </option>
                  ))}
                </select>
                {hasError("id_kelas") && (
                  <p className="mt-1 text-xs text-rose-600">{errors.id_kelas}</p>
                )}
              </div>

              <div>
                <label htmlFor="kelas_scope" className="text-xs font-medium text-slate-600">
                  Cakupan Kelas
                </label>
                <select
                  id="kelas_scope"
                  value={values.kelas_scope}
                  onChange={(event) =>
                    setField(
                      "kelas_scope",
                      event.target.value === "SPESIFIK" ? "SPESIFIK" : "SEMUA",
                    )
                  }
                  onBlur={() => onBlur("kelas_scope")}
                  disabled={values.id_kelas === 0}
                  className={`${selectBaseClass} ${hasError("kelas_scope") ? "border-rose-300 ring-rose-100" : ""}`}
                >
                  <option value="SEMUA">Semua kelas di tingkat ini</option>
                  <option value="SPESIFIK">Spesifik nama kelas</option>
                </select>
                {hasError("kelas_scope") && (
                  <p className="mt-1 text-xs text-rose-600">{errors.kelas_scope}</p>
                )}
              </div>

              <div>
                <label htmlFor="id_nama_kelas" className="text-xs font-medium text-slate-600">
                  Nama Kelas
                </label>
                <select
                  id="id_nama_kelas"
                  value={values.id_nama_kelas}
                  onChange={(event) => setField("id_nama_kelas", Number(event.target.value))}
                  onBlur={() => onBlur("id_nama_kelas")}
                  disabled={values.id_kelas === 0 || values.kelas_scope !== "SPESIFIK"}
                  className={`${selectBaseClass} ${hasError("id_nama_kelas") ? "border-rose-300 ring-rose-100" : ""}`}
                >
                  <option value={0}>Pilih nama kelas</option>
                  {namaKelasOptions.map((item) => (
                    <option key={item.id_nama_kelas} value={item.id_nama_kelas}>
                      {item.nama_kelas}
                    </option>
                  ))}
                </select>
                {hasError("id_nama_kelas") && (
                  <p className="mt-1 text-xs text-rose-600">{errors.id_nama_kelas}</p>
                )}
              </div>

              <div className="md:col-span-2 lg:col-span-1">
                <label htmlFor="id_bank_soal" className="text-xs font-medium text-slate-600">
                  Bank Soal
                </label>
                <select
                  id="id_bank_soal"
                  value={values.id_bank_soal}
                  onChange={(event) => setField("id_bank_soal", Number(event.target.value))}
                  onBlur={() => onBlur("id_bank_soal")}
                  disabled={values.id_kelas === 0 || loadingBankSoal}
                  className={`${selectBaseClass} ${hasError("id_bank_soal") ? "border-rose-300 ring-rose-100" : ""}`}
                >
                  <option value={0}>
                    {loadingBankSoal ? "Memuat bank soal..." : "Pilih bank soal"}
                  </option>
                  {bankSoalOptions.map((item) => (
                    <option key={item.id_bank_soal} value={item.id_bank_soal}>
                      {item.nama_bank_soal}
                    </option>
                  ))}
                </select>
                {hasError("id_bank_soal") && (
                  <p className="mt-1 text-xs text-rose-600">{errors.id_bank_soal}</p>
                )}
              </div>
            </div>
          </div>

          <div className="rounded-xl border border-slate-200 bg-white p-5 shadow-sm">
            <div className="mb-4">
              <h2 className={sectionTitle}>Jadwal Ujian</h2>
              <p className={helperText}>
                Gunakan format waktu 24 jam `HH:mm` seperti 07:30 dan 13:45.
              </p>
            </div>

            <div className="grid grid-cols-1 gap-4 md:grid-cols-4">
              <div>
                <DatePicker
                  id="tanggal_ujian"
                  label="Tanggal Ujian"
                  value={values.tanggal_ujian}
                  onChange={(date) => setField("tanggal_ujian", date)}
                  onBlur={() => onBlur("tanggal_ujian")}
                  error={hasError("tanggal_ujian")}
                />
                {hasError("tanggal_ujian") && (
                  <p className="mt-1 text-xs text-rose-600">{errors.tanggal_ujian}</p>
                )}
              </div>

              <div>
                <TimePicker
                  id="waktu_mulai"
                  label="Waktu Mulai"
                  value={values.waktu_mulai}
                  onChange={(value) => setField("waktu_mulai", value)}
                  onBlur={() => onBlur("waktu_mulai")}
                  error={hasError("waktu_mulai")}
                />
                <p className="mt-1 text-[11px] text-slate-500">Waktu ujian memakai WIB (24 jam).</p>
                {hasError("waktu_mulai") && (
                  <p className="mt-1 text-xs text-rose-600">{errors.waktu_mulai}</p>
                )}
              </div>

              <div>
                <TimePicker
                  id="waktu_selesai"
                  label="Waktu Selesai"
                  value={values.waktu_selesai}
                  onChange={(value) => setField("waktu_selesai", value)}
                  onBlur={() => onBlur("waktu_selesai")}
                  error={hasError("waktu_selesai")}
                />
                <p className="mt-1 text-[11px] text-slate-500">Waktu ujian memakai WIB (24 jam).</p>
                {hasError("waktu_selesai") && (
                  <p className="mt-1 text-xs text-rose-600">{errors.waktu_selesai}</p>
                )}
              </div>

              <div>
                <InputField
                  id="durasi_menit"
                  type="number"
                  label="Durasi (menit)"
                  value={String(durasiMenit)}
                  onChange={() => undefined}
                  inputClassName="bg-slate-50 text-slate-500"
                  disabled
                />
              </div>
            </div>
          </div>

          <div className="rounded-xl border border-slate-200 bg-white p-5 shadow-sm">
            <div className="mb-4">
              <h2 className={sectionTitle}>Ruang, Sesi, Pengawas</h2>
              <p className={helperText}>
                Semua opsi diambil langsung dari endpoint server.
              </p>
            </div>

            <div className="grid grid-cols-1 gap-4 md:grid-cols-3">
              <div>
                <label htmlFor="id_ruangan" className="text-xs font-medium text-slate-600">
                  Ruang Ujian
                </label>
                <select
                  id="id_ruangan"
                  value={values.id_ruangan}
                  onChange={(event) => setField("id_ruangan", Number(event.target.value))}
                  onBlur={() => onBlur("id_ruangan")}
                  disabled={loadingRuang}
                  className={`${selectBaseClass} ${hasError("id_ruangan") ? "border-rose-300 ring-rose-100" : ""}`}
                >
                  <option value={0}>{loadingRuang ? "Memuat ruang..." : "Pilih ruang ujian"}</option>
                  {ruangOptions.map((item) => (
                    <option key={item.id_ruangan} value={item.id_ruangan}>
                      {item.nama_ruangan}
                    </option>
                  ))}
                </select>
                {hasError("id_ruangan") && (
                  <p className="mt-1 text-xs text-rose-600">{errors.id_ruangan}</p>
                )}
              </div>

              <div>
                <label htmlFor="id_sesi" className="text-xs font-medium text-slate-600">
                  Sesi Ujian
                </label>
                <select
                  id="id_sesi"
                  value={values.id_sesi}
                  onChange={(event) => setField("id_sesi", Number(event.target.value))}
                  onBlur={() => onBlur("id_sesi")}
                  disabled={loadingSesi}
                  className={`${selectBaseClass} ${hasError("id_sesi") ? "border-rose-300 ring-rose-100" : ""}`}
                >
                  <option value={0}>{loadingSesi ? "Memuat sesi..." : "Pilih sesi ujian"}</option>
                  {sesiOptions.map((item) => (
                    <option key={item.id_sesi} value={item.id_sesi}>
                      {item.kode_sesi} - {item.nama_sesi}
                    </option>
                  ))}
                </select>
                {hasError("id_sesi") && (
                  <p className="mt-1 text-xs text-rose-600">{errors.id_sesi}</p>
                )}
              </div>

              <div>
                <label htmlFor="id_pengawas" className="text-xs font-medium text-slate-600">
                  Guru Pengawas
                </label>
                <select
                  id="id_pengawas"
                  value={values.id_pengawas}
                  onChange={(event) => setField("id_pengawas", Number(event.target.value))}
                  onBlur={() => onBlur("id_pengawas")}
                  disabled={loadingGuru}
                  className={`${selectBaseClass} ${hasError("id_pengawas") ? "border-rose-300 ring-rose-100" : ""}`}
                >
                  <option value={0}>{loadingGuru ? "Memuat guru..." : "Pilih guru pengawas"}</option>
                  {guruOptions.map((item) => (
                    <option key={item.id_pengguna} value={item.id_pengguna}>
                      {item.nama_lengkap}
                    </option>
                  ))}
                </select>
                {hasError("id_pengawas") && (
                  <p className="mt-1 text-xs text-rose-600">{errors.id_pengawas}</p>
                )}
              </div>
            </div>
          </div>

          <div className="rounded-xl border border-slate-200 bg-white p-5 shadow-sm">
            <div className="mb-4">
              <h2 className={sectionTitle}>Keamanan & Token</h2>
              <p className={helperText}>Token bebas, maksimal 30 karakter.</p>
            </div>

            <div className="grid grid-cols-1 gap-4 md:grid-cols-2">
              <div>
                <label htmlFor="acak_soal" className="text-xs font-medium text-slate-600">
                  Acak Soal
                </label>
                <select
                  id="acak_soal"
                  value={values.acak_soal ? "ya" : "tidak"}
                  onChange={(event) => setField("acak_soal", event.target.value === "ya")}
                  className={selectBaseClass}
                >
                  <option value="ya">Ya, acak soal</option>
                  <option value="tidak">Tidak, urutan tetap</option>
                </select>
              </div>

              <div>
                <InputField
                  id="token"
                  label="Token Ujian"
                  value={values.token}
                  onChange={(value) => setField("token", value.slice(0, 30))}
                  onBlur={() => onBlur("token")}
                  placeholder="Contoh: UTS-MTK-01"
                  inputClassName={hasError("token") ? "border-rose-300 ring-rose-100" : ""}
                />
                <p className="mt-1 text-[11px] text-slate-500">
                  {values.token.length}/30 karakter
                </p>
                {hasError("token") && (
                  <p className="mt-1 text-xs text-rose-600">{errors.token}</p>
                )}
              </div>
            </div>
          </div>

          <div className="rounded-xl border border-slate-200 bg-white p-5 shadow-sm">
            <div className="mb-4 flex items-center justify-between">
              <div>
                <h2 className={sectionTitle}>Preview Daftar Siswa</h2>
                <p className={helperText}>
                  {values.kelas_scope === "SPESIFIK"
                    ? "Preview mengambil id_tingkat_kelas dan id_nama_kelas."
                    : "Preview mengambil id_tingkat_kelas saja."}
                </p>
              </div>
              <span className="rounded-full bg-slate-100 px-3 py-1 text-xs font-medium text-slate-600">
                {siswaPreview.length} siswa
              </span>
            </div>

            {!siswaPreviewEnabled && (
              <p className="text-sm text-slate-500">
                Pilih tingkat kelas, dan jika scope spesifik pilih nama kelas, untuk melihat siswa.
              </p>
            )}

            {siswaPreviewEnabled && loadingSiswa && (
              <p className="text-sm text-slate-500">Memuat data siswa...</p>
            )}

            {siswaPreviewEnabled && !loadingSiswa && siswaPreview.length === 0 && (
              <p className="text-sm text-slate-500">Belum ada siswa pada filter kelas ini.</p>
            )}

            {siswaPreview.length > 0 && (
              <div className="max-h-64 overflow-y-auto rounded-lg border border-slate-200">
                <table className="w-full text-left text-xs text-slate-600">
                  <thead className="bg-slate-50 text-[11px] uppercase text-slate-400">
                    <tr>
                      <th className="px-3 py-2 font-medium">Nama</th>
                      <th className="px-3 py-2 font-medium">Absen</th>
                      <th className="px-3 py-2 font-medium">Kelas</th>
                      <th className="px-3 py-2 font-medium">Status</th>
                    </tr>
                  </thead>
                  <tbody>
                    {siswaPreview.map((siswa) => (
                      <tr key={siswa.id_pengguna} className="border-t border-slate-100">
                        <td className="px-3 py-2 font-medium text-slate-700">{siswa.nama_lengkap}</td>
                        <td className="px-3 py-2">{siswa.no_absen}</td>
                        <td className="px-3 py-2">{siswa.nama_kelas}</td>
                        <td className="px-3 py-2 capitalize">{siswa.status_akun}</td>
                      </tr>
                    ))}
                  </tbody>
                </table>
              </div>
            )}
          </div>

          {(submitError || loadErrorMessage) && (
            <div className="rounded-lg border border-rose-200 bg-rose-50 px-4 py-3 text-sm text-rose-600">
              {submitError || loadErrorMessage}
            </div>
          )}

          <div className="flex flex-col gap-3 sm:flex-row sm:justify-end">
            <button
              type="button"
              className="inline-flex cursor-pointer items-center justify-center rounded-lg border border-slate-200 px-4 py-2 text-sm font-medium text-slate-600 transition hover:bg-slate-50 disabled:cursor-not-allowed"
              onClick={() => {
                setValues(initialValues);
                setTouched({});
                setSubmitError(null);
              }}
              disabled={submitting}
            >
              Reset
            </button>

            <button
              type="submit"
              className="inline-flex cursor-pointer items-center justify-center rounded-lg bg-[#397e50] px-4 py-2 text-sm font-semibold text-white shadow-sm transition hover:bg-[#2f6a43] disabled:cursor-not-allowed disabled:opacity-70"
              disabled={submitting}
            >
              {submitting ? "Menyimpan..." : "Simpan Ujian"}
            </button>
          </div>
        </form>
      </div>
    </div>
  );
};

export default BuatUjianForm;
