import React, { useEffect, useMemo, useState } from "react";

import InputField from "@/components/common/Input/InputField";

import type {
  BankSoalOption,
  BuatUjianFormValues,
  GuruPengawasOption,
  SesiUjianOption,
  SiswaPreviewItem,
  TipeUjian,
} from "@/types/Ujian/BuatUjian";
import type { NamaKelas, TingkatKelas } from "@/types/DataMaster/Kelas";
import {
  getUjianSiswaPreview,
  submitBuatUjian,
} from "@/services/Api/features-api/Ujian/ujian.service";

import {
  getNamaKelas,
  getTingkatKelas,
  getTingkatKelasById,
} from "@/services/Api/features-api/DataMaster/kelas.service";
import {
  getRuangUjianOptions,
  getUjianBankSoalOptions,
  getUjianGuruPengawasOptions,
  getUjianSesiOptions,
} from "@/services/Api/features-api/GetOptions/options.service";

// helper
import { createSetField } from "@/helper/setField/setField";
import { calculateDuration } from "@/helper/CalculateDuration/calculateDuration";
import { getSubmitErrorMessage } from "@/helper/error/submitErrorMessage";
import {
  createValidator,
  minNumber,
  requiredString,
  requiredValue,
} from "@/helper/validate/validateForm";
import type { RuangUjianRow } from "@/types/DataMaster/RuangUjian";

const initialValues: BuatUjianFormValues = {
  nama_ujian: "",
  deskripsi_ujian: "",
  tipe_ujian: "PILIHAN_GANDA",
  kelas_id: 0,
  kelas_scope: "SEMUA",
  kelas_detail_id: 0,
  bank_soal_id: 0,
  jumlah_soal: 0,
  tanggal_ujian: "",
  waktu_mulai: "",
  waktu_selesai: "",
  durasi_menit: 0,
  ruang_ujian_id: 0,
  acak_soal: true,
  guru_pengawas_id: 0,
  sesi_id: 0,
  token_ujian: "",
};

const sectionTitle = "text-sm font-semibold text-slate-800";
const helperText = "text-xs text-slate-500";

const selectBaseClass =
  "w-full rounded-lg border border-slate-200 bg-white px-3 py-2 text-sm outline-none transition focus:border-[#397e50] focus:ring-2 focus:ring-[#397e50] disabled:bg-slate-50 disabled:text-slate-500";

const BuatUjianForm = () => {
  const [values, setValues] = useState<BuatUjianFormValues>(initialValues);
  const [touched, setTouched] = useState<Record<string, boolean>>({});
  const [submitError, setSubmitError] = useState<string | null>(null);
  const [loadingBankSoal, setLoadingBankSoal] = useState(false);
  const [loadingSiswa, setLoadingSiswa] = useState(false);
  const [loadingKelasDetail, setLoadingKelasDetail] = useState(false);

  const [kelasOptions, setKelasOptions] = useState<TingkatKelas[]>([]);
  const [kelasDetailOptions, setKelasDetailOptions] = useState<NamaKelas[]>([]);
  const [bankSoalOptions, setBankSoalOptions] = useState<BankSoalOption[]>([]);
  const [ruangOptions, setRuangOptions] = useState<RuangUjianRow[]>([]);
  const [guruOptions, setGuruOptions] = useState<GuruPengawasOption[]>([]);
  const [sesiOptions, setSesiOptions] = useState<SesiUjianOption[]>([]);
  const [siswaPreview, setSiswaPreview] = useState<SiswaPreviewItem[]>([]);

  const setField = createSetField(setValues);

  const onBlur = (name: keyof BuatUjianFormValues) => {
    setTouched((prev) => ({ ...prev, [name]: true }));
  };

  const selectedKelasId = useMemo(() => {
    const tingkatKelas =
      values.kelas_id === 0 ? undefined : getTingkatKelasById(values.kelas_id);
    return tingkatKelas ?? undefined;
  }, [values.kelas_id]);

  const kelasDetailById = useMemo(() => {
    return new Map(
      kelasDetailOptions.map((kelas) => [kelas.id_nama_kelas, kelas]),
    );
  }, [kelasDetailOptions]);

  const selectedKelasDetail =
    values.kelas_detail_id === 0
      ? undefined
      : kelasDetailById.get(values.kelas_detail_id);

  const bankSoalById = useMemo(() => {
    const map = new Map(bankSoalOptions.map((x) => [x.id, x]));
    return map;
  }, [bankSoalOptions]);

  const selectedBankSoal =
    values.bank_soal_id === 0
      ? undefined
      : bankSoalById.get(values.bank_soal_id);

  useEffect(() => {
    let active = true;

    const loadOptions = async () => {
      try {
        const [kelas, ruang, guru, sesi] = await Promise.all([
          getTingkatKelas(),
          getRuangUjianOptions(),
          getUjianGuruPengawasOptions(),
          getUjianSesiOptions(),
        ]);

        if (!active) return;
        setKelasOptions(kelas);
        setRuangOptions(ruang);
        setGuruOptions(guru);
        setSesiOptions(sesi);
      } catch (error) {
        if (!active) return;
        setSubmitError("Gagal memuat data pendukung ujian.");
      }
    };

    loadOptions();
    return () => {
      active = false;
    };
  }, []);

  useEffect(() => {
    let active = true;
    const loadKelasDetail = async () => {
      setLoadingKelasDetail(true);
      try {
        const data = await getNamaKelas({
          tingkatKelas: values.kelas_id === 0 ? undefined : values.kelas_id,
        });
        if (!active) return;
        setKelasDetailOptions(data);
      } finally {
        if (active) setLoadingKelasDetail(false);
      }
    };

    if (values.kelas_id !== 0) {
      loadKelasDetail();
    } else {
      setKelasDetailOptions([]);
    }

    setField("kelas_detail_id", 0);

    return () => {
      active = false;
    };
  }, [values.kelas_id]);

  useEffect(() => {
    if (values.kelas_scope === "SEMUA") {
      setField("kelas_detail_id", 0);
    }
  }, [values.kelas_scope]);

  useEffect(() => {
    let active = true;
    const loadBankSoal = async () => {
      setLoadingBankSoal(true);
      try {
        const data = await getUjianBankSoalOptions({
          tingkatKelasId: values.kelas_id === 0 ? undefined : values.kelas_id,
        });
        if (!active) return;
        setBankSoalOptions(data);
      } finally {
        if (active) setLoadingBankSoal(false);
      }
    };

    if (selectedKelasId != null) {
      loadBankSoal();
    } else {
      setBankSoalOptions([]);
      setField("bank_soal_id", 0);
      setField("jumlah_soal", 0);
    }

    return () => {
      active = false;
    };
  }, [selectedKelasId]);

  useEffect(() => {
    let active = true;
    const loadSiswa = async () => {
      setLoadingSiswa(true);
      try {
        const data = await getUjianSiswaPreview({
          tingkatKelasId: values.kelas_id === 0 ? undefined : values.kelas_id,
        });
        if (!active) return;
        // harusnya data dari server sudah di sorting
        setSiswaPreview(data);
      } finally {
        if (active) setLoadingSiswa(false);
      }
    };

    if (selectedKelasId != null) {
      loadSiswa();
    } else {
      setSiswaPreview([]);
    }

    return () => {
      active = false;
    };
  }, [selectedKelasId]);

  useEffect(() => {
    if (selectedBankSoal) {
      setField("jumlah_soal", selectedBankSoal.total_soal);
    } else {
      setField("jumlah_soal", 0);
    }
  }, [selectedBankSoal]);

  useEffect(() => {
    const duration = calculateDuration(
      values.waktu_mulai,
      values.waktu_selesai,
    );
    setField("durasi_menit", duration);
  }, [values.waktu_mulai, values.waktu_selesai]);

  const siswaPreviewFiltered = useMemo(() => {
    if (values.kelas_scope !== "SPESIFIK") {
      return siswaPreview;
    }

    if (!selectedKelasDetail) {
      return [];
    }

    return siswaPreview.filter(
      (siswa) => siswa.kelas === selectedKelasDetail.nama_kelas,
    );
  }, [siswaPreview, values.kelas_scope, selectedKelasDetail]);

  const validate = createValidator<BuatUjianFormValues>({
    nama_ujian: [requiredString("Nama ujian wajib diisi.")],
    deskripsi_ujian: [requiredString("Deskripsi ujian wajib diisi.")],
    tipe_ujian: [requiredValue("Tipe ujian wajib dipilih.")],
    kelas_id: [minNumber(1, "Tingkat kelas wajib dipilih.")],
    kelas_scope: [requiredValue("Cakupan kelas wajib dipilih.")],
    kelas_detail_id: [
      (value, values) =>
        values.kelas_scope === "SPESIFIK" && value === 0
          ? "Nama kelas wajib dipilih."
          : null,
    ],
    bank_soal_id: [minNumber(1, "Bank soal wajib dipilih.")],
    tanggal_ujian: [requiredString("Tanggal ujian wajib diisi.")],
    waktu_mulai: [requiredString("Waktu mulai wajib diisi.")],
    waktu_selesai: [
      requiredString("Waktu selesai wajib diisi."),
      (_, values) =>
        values.waktu_mulai && values.waktu_selesai && values.durasi_menit <= 0
          ? "Waktu selesai harus setelah waktu mulai."
          : null,
    ],
    ruang_ujian_id: [minNumber(1, "Ruang ujian wajib dipilih.")],
    guru_pengawas_id: [minNumber(1, "Guru pengawas wajib dipilih.")],
    sesi_id: [minNumber(1, "Sesi ujian wajib dipilih.")],
    token_ujian: [requiredString("Token ujian wajib diisi.")],
  });

  const errors = validate(values);
  const hasError = (name: keyof BuatUjianFormValues) =>
    !!errors[name] && !!touched[name];

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
      setSubmitError(
        "Periksa kembali input yang masih kosong atau belum valid.",
      );
      return;
    }

    try {
      await submitBuatUjian(values);
      alert("Ujian berhasil dibuat.");
    } catch (error) {
      setSubmitError(
        getSubmitErrorMessage(error, {
          defaultMessage: "Terjadi kesalahan saat menyimpan ujian.",
        }),
      );
    }
  };

  return (
    <div className="min-h-screen w-full py-8">
      <div className="mx-auto w-full max-w-6xl px-4">
        <div className="mb-6 rounded-xl border border-slate-200 bg-white p-5 shadow-sm">
          <h1 className="text-base font-semibold text-slate-900">Buat Ujian</h1>
          <p className="mt-1 text-sm text-slate-500">
            Lengkapi detail ujian, jadwal, serta bank soal yang digunakan.
          </p>
        </div>

        <form onSubmit={onSubmit} className="space-y-6">
          <div className="rounded-xl border border-slate-200 bg-white p-5 shadow-sm">
            <div className="mb-4">
              <h2 className={sectionTitle}>Informasi Ujian</h2>
              <p className={helperText}>
                Isi nama ujian, deskripsi, dan tipe soal yang digunakan.
              </p>
            </div>

            <div className="grid grid-cols-1 gap-4 md:grid-cols-2">
              <div className="md:col-span-2">
                <InputField
                  id="nama_ujian"
                  label="Nama Ujian"
                  value={values.nama_ujian}
                  onChange={(v) => setField("nama_ujian", v)}
                  onBlur={() => onBlur("nama_ujian")}
                  placeholder="Contoh: Ujian Tengah Semester Matematika"
                  inputClassName={
                    hasError("nama_ujian")
                      ? "border-rose-300 ring-rose-100"
                      : ""
                  }
                  required
                />
                {hasError("nama_ujian") && (
                  <p className="mt-1 text-xs text-rose-600">
                    {errors.nama_ujian}
                  </p>
                )}
              </div>

              <div className="md:col-span-2">
                <label
                  htmlFor="deskripsi_ujian"
                  className="text-xs font-medium text-slate-600"
                >
                  Deskripsi Ujian
                </label>
                <textarea
                  id="deskripsi_ujian"
                  className={`min-h-[100px] w-full rounded-lg border border-slate-200 bg-white px-3 py-2 text-sm outline-none transition focus:border-[#397e50] focus:ring-2 focus:ring-[#397e50] ${
                    hasError("deskripsi_ujian")
                      ? "border-rose-300 ring-rose-100"
                      : ""
                  }`}
                  placeholder="Jelaskan ujian ini untuk apa dan cakupan materinya."
                  value={values.deskripsi_ujian}
                  onChange={(e) => setField("deskripsi_ujian", e.target.value)}
                  onBlur={() => onBlur("deskripsi_ujian")}
                  required
                />
                {hasError("deskripsi_ujian") && (
                  <p className="mt-1 text-xs text-rose-600">
                    {errors.deskripsi_ujian}
                  </p>
                )}
              </div>

              <div>
                <label
                  htmlFor="tipe_ujian"
                  className="text-xs font-medium text-slate-600"
                >
                  Tipe Ujian
                </label>
                <select
                  id="tipe_ujian"
                  className={`${selectBaseClass} ${
                    hasError("tipe_ujian")
                      ? "border-rose-300 ring-rose-100"
                      : ""
                  }`}
                  value={values.tipe_ujian}
                  onChange={(e) =>
                    setField("tipe_ujian", e.target.value as TipeUjian)
                  }
                  onBlur={() => onBlur("tipe_ujian")}
                >
                  <option value="PILIHAN_GANDA">Pilihan Ganda</option>
                  <option value="ESSAY">Essay</option>
                  <option value="CAMPURAN">Pilihan Ganda & Essay</option>
                </select>
                {hasError("tipe_ujian") && (
                  <p className="mt-1 text-xs text-rose-600">
                    {errors.tipe_ujian}
                  </p>
                )}
              </div>
            </div>
          </div>

          <div className="rounded-xl border border-slate-200 bg-white p-5 shadow-sm">
            <div className="mb-4">
              <h2 className={sectionTitle}>Tingkat Kelas & Bank Soal</h2>
              <p className={helperText}>
                Tentukan tingkat, cakupan kelas, dan bank soal ujian.
              </p>
            </div>

            <div className="grid grid-cols-1 gap-4 md:grid-cols-2 lg:grid-cols-3">
              <div>
                <label
                  htmlFor="kelas_id"
                  className="text-xs font-medium text-slate-600"
                >
                  Tingkat Kelas
                </label>
                <select
                  id="kelas_id"
                  className={`${selectBaseClass} ${
                    hasError("kelas_id") ? "border-rose-300 ring-rose-100" : ""
                  }`}
                  value={values.kelas_id}
                  onChange={(e) => setField("kelas_id", Number(e.target.value))}
                  onBlur={() => onBlur("kelas_id")}
                  required
                >
                  <option value={0}>Pilih tingkat kelas</option>
                  {kelasOptions.map((tingkat) => (
                    <option
                      key={tingkat.id_tingkat_kelas}
                      value={tingkat.id_tingkat_kelas}
                    >
                      Kelas {tingkat.tingkat_kelas}
                    </option>
                  ))}
                </select>
                {hasError("kelas_id") && (
                  <p className="mt-1 text-xs text-rose-600">
                    {errors.kelas_id}
                  </p>
                )}
              </div>

              <div>
                <label
                  htmlFor="kelas_scope"
                  className="text-xs font-medium text-slate-600"
                >
                  Cakupan Kelas
                </label>
                <select
                  id="kelas_scope"
                  className={`${selectBaseClass} ${
                    hasError("kelas_scope")
                      ? "border-rose-300 ring-rose-100"
                      : ""
                  }`}
                  value={values.kelas_scope}
                  onChange={(e) =>
                    setField(
                      "kelas_scope",
                      e.target.value === "SPESIFIK" ? "SPESIFIK" : "SEMUA",
                    )
                  }
                  onBlur={() => onBlur("kelas_scope")}
                  disabled={values.kelas_id === 0}
                >
                  <option value="SEMUA">Semua kelas di tingkat ini</option>
                  <option value="SPESIFIK">Spesifik nama kelas</option>
                </select>
                {hasError("kelas_scope") && (
                  <p className="mt-1 text-xs text-rose-600">
                    {errors.kelas_scope}
                  </p>
                )}
              </div>

              <div>
                <label
                  htmlFor="kelas_detail_id"
                  className="text-xs font-medium text-slate-600"
                >
                  Nama Kelas
                </label>
                <select
                  id="kelas_detail_id"
                  className={`${selectBaseClass} ${
                    hasError("kelas_detail_id")
                      ? "border-rose-300 ring-rose-100"
                      : ""
                  }`}
                  value={values.kelas_detail_id}
                  onChange={(e) =>
                    setField("kelas_detail_id", Number(e.target.value))
                  }
                  onBlur={() => onBlur("kelas_detail_id")}
                  disabled={
                    values.kelas_id === 0 ||
                    values.kelas_scope !== "SPESIFIK" ||
                    loadingKelasDetail
                  }
                >
                  <option value={0}>
                    {loadingKelasDetail
                      ? "Memuat nama kelas..."
                      : "Pilih nama kelas"}
                  </option>
                  {kelasDetailOptions.map((kelas) => (
                    <option
                      key={kelas.id_nama_kelas}
                      value={kelas.id_nama_kelas}
                    >
                      {kelas.nama_kelas}
                    </option>
                  ))}
                </select>
                {hasError("kelas_detail_id") && (
                  <p className="mt-1 text-xs text-rose-600">
                    {errors.kelas_detail_id}
                  </p>
                )}
              </div>

              <div>
                <label
                  htmlFor="bank_soal_id"
                  className="text-xs font-medium text-slate-600"
                >
                  Bank Soal
                </label>
                <select
                  id="bank_soal_id"
                  className={`${selectBaseClass} ${
                    hasError("bank_soal_id")
                      ? "border-rose-300 ring-rose-100"
                      : ""
                  }`}
                  value={values.bank_soal_id}
                  onChange={(e) =>
                    setField("bank_soal_id", Number(e.target.value))
                  }
                  onBlur={() => onBlur("bank_soal_id")}
                  disabled={values.kelas_id === 0 || loadingBankSoal}
                >
                  <option value={0}>
                    {loadingBankSoal
                      ? "Memuat bank soal..."
                      : "Pilih bank soal"}
                  </option>
                  {bankSoalOptions.map((bank) => (
                    <option key={bank.id} value={bank.id}>
                      {bank.nama}
                    </option>
                  ))}
                </select>
                {hasError("bank_soal_id") && (
                  <p className="mt-1 text-xs text-rose-600">
                    {errors.bank_soal_id}
                  </p>
                )}
              </div>

              <div className="lg:col-start-3">
                <InputField
                  id="jumlah_soal"
                  type="number"
                  label="Jumlah Soal"
                  value={String(values.jumlah_soal)}
                  onChange={(v) => setField("jumlah_soal", Number(v))}
                  onBlur={() => onBlur("jumlah_soal")}
                  inputClassName="bg-slate-50 text-slate-500"
                  disabled
                  required
                />
              </div>
            </div>
          </div>

          <div className="rounded-xl border border-slate-200 bg-white p-5 shadow-sm">
            <div className="mb-4">
              <h2 className={sectionTitle}>Jadwal Ujian</h2>
              <p className={helperText}>
                Tentukan tanggal, waktu mulai, dan waktu berakhir ujian.
              </p>
            </div>

            <div className="grid grid-cols-1 gap-4 md:grid-cols-4">
              <div>
                <InputField
                  id="tanggal_ujian"
                  type="date"
                  label="Tanggal Ujian"
                  value={values.tanggal_ujian}
                  onChange={(v) => setField("tanggal_ujian", v)}
                  onBlur={() => onBlur("tanggal_ujian")}
                  inputClassName={
                    hasError("tanggal_ujian")
                      ? "border-rose-300 ring-rose-100"
                      : ""
                  }
                  required
                />
                {hasError("tanggal_ujian") && (
                  <p className="mt-1 text-xs text-rose-600">
                    {errors.tanggal_ujian}
                  </p>
                )}
              </div>

              <div>
                <InputField
                  id="waktu_mulai"
                  type="time/datetime-local"
                  label="Waktu Mulai"
                  value={values.waktu_mulai}
                  onChange={(v) => setField("waktu_mulai", v)}
                  onBlur={() => onBlur("waktu_mulai")}
                  inputClassName={
                    hasError("waktu_mulai")
                      ? "border-rose-300 ring-rose-100"
                      : ""
                  }
                  required
                />
                {hasError("waktu_mulai") && (
                  <p className="mt-1 text-xs text-rose-600">
                    {errors.waktu_mulai}
                  </p>
                )}
              </div>

              <div>
                <InputField
                  id="waktu_selesai"
                  type="time/datetime-local"
                  label="Waktu Selesai"
                  value={values.waktu_selesai}
                  onChange={(v) => setField("waktu_selesai", v)}
                  onBlur={() => onBlur("waktu_selesai")}
                  inputClassName={
                    hasError("waktu_selesai")
                      ? "border-rose-300 ring-rose-100"
                      : ""
                  }
                  required
                />
                {hasError("waktu_selesai") && (
                  <p className="mt-1 text-xs text-rose-600">
                    {errors.waktu_selesai}
                  </p>
                )}
              </div>

              <div>
                <InputField
                  id="durasi_menit"
                  type="number"
                  label="Durasi (menit)"
                  value={String(values.durasi_menit)}
                  onChange={(v) => setField("durasi_menit", Number(v))}
                  onBlur={() => onBlur("durasi_menit")}
                  inputClassName="bg-slate-50 text-slate-500"
                  disabled
                  required
                />
              </div>
            </div>
          </div>

          <div className="rounded-xl border border-slate-200 bg-white p-5 shadow-sm">
            <div className="mb-4">
              <h2 className={sectionTitle}>Ruang, Sesi, Pengawas</h2>
              <p className={helperText}>
                Atur ruang ujian, sesi, dan guru pengawas.
              </p>
            </div>

            <div className="grid grid-cols-1 gap-4 md:grid-cols-3">
              <div>
                <label
                  htmlFor="ruang_ujian_id"
                  className="text-xs font-medium text-slate-600"
                >
                  Ruang Ujian
                </label>
                <select
                  id="ruang_ujian_id"
                  className={`${selectBaseClass} ${
                    hasError("ruang_ujian_id")
                      ? "border-rose-300 ring-rose-100"
                      : ""
                  }`}
                  value={values.ruang_ujian_id}
                  onChange={(e) =>
                    setField("ruang_ujian_id", Number(e.target.value))
                  }
                  onBlur={() => onBlur("ruang_ujian_id")}
                  required
                >
                  <option value={0}>Pilih ruang ujian</option>
                  {ruangOptions.map((ruang) => (
                    <option key={ruang.id} value={ruang.id}>
                      {ruang.namaRuangan}
                    </option>
                  ))}
                </select>
                {hasError("ruang_ujian_id") && (
                  <p className="mt-1 text-xs text-rose-600">
                    {errors.ruang_ujian_id}
                  </p>
                )}
              </div>

              <div>
                <label
                  htmlFor="sesi_id"
                  className="text-xs font-medium text-slate-600"
                >
                  Sesi Ujian
                </label>
                <select
                  id="sesi_id"
                  className={`${selectBaseClass} ${
                    hasError("sesi_id") ? "border-rose-300 ring-rose-100" : ""
                  }`}
                  value={values.sesi_id}
                  onChange={(e) => setField("sesi_id", Number(e.target.value))}
                  onBlur={() => onBlur("sesi_id")}
                  required
                >
                  <option value={0}>Pilih sesi ujian</option>
                  {sesiOptions.map((sesi) => (
                    <option key={sesi.id} value={sesi.id}>
                      {sesi.kode} - {sesi.nama}
                    </option>
                  ))}
                </select>
                {hasError("sesi_id") && (
                  <p className="mt-1 text-xs text-rose-600">{errors.sesi_id}</p>
                )}
              </div>

              <div>
                <label
                  htmlFor="guru_pengawas_id"
                  className="text-xs font-medium text-slate-600"
                >
                  Guru Pengawas
                </label>
                <select
                  id="guru_pengawas_id"
                  className={`${selectBaseClass} ${
                    hasError("guru_pengawas_id")
                      ? "border-rose-300 ring-rose-100"
                      : ""
                  }`}
                  value={values.guru_pengawas_id}
                  onChange={(e) =>
                    setField("guru_pengawas_id", Number(e.target.value))
                  }
                  onBlur={() => onBlur("guru_pengawas_id")}
                  required
                >
                  <option value={0}>Pilih guru pengawas</option>
                  {guruOptions.map((guru) => (
                    <option key={guru.id} value={guru.id}>
                      {guru.nama} {guru.mapel ? `- ${guru.mapel}` : ""}
                    </option>
                  ))}
                </select>
                {hasError("guru_pengawas_id") && (
                  <p className="mt-1 text-xs text-rose-600">
                    {errors.guru_pengawas_id}
                  </p>
                )}
              </div>
            </div>
          </div>

          <div className="rounded-xl border border-slate-200 bg-white p-5 shadow-sm">
            <div className="mb-4">
              <h2 className={sectionTitle}>Keamanan & Token</h2>
              <p className={helperText}>
                Atur pengacakan soal dan token ujian.
              </p>
            </div>

            <div className="grid grid-cols-1 gap-4 md:grid-cols-2">
              <div>
                <label
                  htmlFor="acak_soal"
                  className="text-xs font-medium text-slate-600"
                >
                  Acak Soal
                </label>
                <select
                  id="acak_soal"
                  className={selectBaseClass}
                  value={values.acak_soal ? "ya" : "tidak"}
                  onChange={(e) =>
                    setField("acak_soal", e.target.value === "ya")
                  }
                  onBlur={() => onBlur("acak_soal")}
                >
                  <option value="ya">Ya, acak soal</option>
                  <option value="tidak">Tidak, urutan tetap</option>
                </select>
              </div>

              <div>
                <InputField
                  id="token_ujian"
                  label="Token Ujian"
                  value={values.token_ujian}
                  onChange={(v) => setField("token_ujian", v)}
                  onBlur={() => onBlur("token_ujian")}
                  placeholder="Masukkan token ujian"
                  inputClassName={
                    hasError("token_ujian")
                      ? "border-rose-300 ring-rose-100"
                      : ""
                  }
                  required
                />
                {hasError("token_ujian") && (
                  <p className="mt-1 text-xs text-rose-600">
                    {errors.token_ujian}
                  </p>
                )}
              </div>
            </div>
          </div>

          <div className="grid grid-cols-1 gap-6 lg:grid-cols-2">
            <div className="rounded-xl border border-slate-200 bg-white p-5 shadow-sm">
              <div className="mb-4 flex items-center justify-between">
                <div>
                  <h2 className={sectionTitle}>Preview Daftar Siswa</h2>
                  <p className={helperText}>
                    {values.kelas_scope === "SPESIFIK"
                      ? "Siswa terdaftar berdasarkan tingkat dan nama kelas yang dipilih."
                      : "Siswa terdaftar berdasarkan tingkat kelas yang dipilih."}
                  </p>
                </div>
                <span className="rounded-full bg-slate-100 px-3 py-1 text-xs font-medium text-slate-600">
                  {siswaPreviewFiltered.length} siswa
                </span>
              </div>

              {values.kelas_id === 0 && (
                <p className="text-sm text-slate-500">
                  Pilih tingkat kelas untuk melihat daftar siswa.
                </p>
              )}

              {values.kelas_id !== 0 &&
                values.kelas_scope === "SPESIFIK" &&
                values.kelas_detail_id === 0 && (
                  <p className="text-sm text-slate-500">
                    Pilih nama kelas untuk melihat daftar siswa.
                  </p>
                )}

              {values.kelas_id !== 0 &&
                loadingSiswa &&
                (values.kelas_scope !== "SPESIFIK" ||
                  values.kelas_detail_id !== 0) && (
                  <p className="text-sm text-slate-500">Memuat data siswa...</p>
                )}

              {values.kelas_id !== 0 &&
                !loadingSiswa &&
                (values.kelas_scope !== "SPESIFIK" ||
                  values.kelas_detail_id !== 0) &&
                siswaPreviewFiltered.length === 0 && (
                  <p className="text-sm text-slate-500">
                    {values.kelas_scope === "SPESIFIK"
                      ? "Belum ada siswa pada kelas ini."
                      : "Belum ada siswa pada tingkat kelas ini."}
                  </p>
                )}

              {siswaPreviewFiltered.length > 0 && (
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
                      {siswaPreviewFiltered.map((siswa) => (
                        <tr
                          key={siswa.id}
                          className="border-t border-slate-100"
                        >
                          <td className="px-3 py-2 font-medium text-slate-700">
                            {siswa.nama}
                          </td>
                          <td className="px-3 py-2">{siswa.no_absen}</td>
                          <td className="px-3 py-2">{siswa.kelas}</td>
                          <td className="px-3 py-2 capitalize">
                            {siswa.status_akun}
                          </td>
                        </tr>
                      ))}
                    </tbody>
                  </table>
                </div>
              )}
            </div>

            <div className="rounded-xl border border-slate-200 bg-white p-5 shadow-sm">
              <div className="mb-4">
                <h2 className={sectionTitle}>Preview Bank Soal</h2>
                <p className={helperText}>
                  Detail bank soal yang dipilih akan tampil di sini.
                </p>
              </div>

              {!selectedBankSoal && (
                <p className="text-sm text-slate-500">
                  Pilih bank soal untuk melihat detailnya.
                </p>
              )}

              {selectedBankSoal && (
                <div className="space-y-3 text-sm text-slate-600">
                  <div>
                    <p className="text-xs text-slate-400">Nama Bank Soal</p>
                    <p className="font-semibold text-slate-800">
                      {selectedBankSoal.nama}
                    </p>
                  </div>
                  <div className="grid grid-cols-2 gap-3">
                    <div>
                      <p className="text-xs text-slate-400">Mata Pelajaran</p>
                      <p className="font-medium text-slate-700">
                        {selectedBankSoal.mata_pelajaran ?? "-"}
                      </p>
                    </div>
                    <div>
                      <p className="text-xs text-slate-400">Materi</p>
                      <p className="font-medium text-slate-700">
                        {selectedBankSoal.materi ?? "-"}
                      </p>
                    </div>
                    <div>
                      <p className="text-xs text-slate-400">Jumlah PG</p>
                      <p className="font-medium text-slate-700">
                        {selectedBankSoal.jumlah_soal_pg ?? 0}
                      </p>
                    </div>
                    <div>
                      <p className="text-xs text-slate-400">Jumlah Essay</p>
                      <p className="font-medium text-slate-700">
                        {selectedBankSoal.jumlah_soal_essay ?? 0}
                      </p>
                    </div>
                  </div>
                  <div className="rounded-lg bg-slate-50 px-3 py-2 text-sm font-semibold text-slate-700">
                    Total Soal: {selectedBankSoal.total_soal}
                  </div>
                </div>
              )}
            </div>
          </div>

          {submitError && (
            <div className="rounded-lg border border-rose-200 bg-rose-50 px-4 py-3 text-sm text-rose-600">
              {submitError}
            </div>
          )}

          <div className="flex flex-col gap-3 sm:flex-row sm:justify-end">
            <button
              type="button"
              className="inline-flex items-center justify-center rounded-lg border border-slate-200 px-4 py-2 text-sm font-medium text-slate-600 transition hover:bg-slate-50"
              onClick={() => {
                setValues(initialValues);
                setTouched({});
                setSubmitError(null);
                setBankSoalOptions([]);
                setKelasDetailOptions([]);
                setSiswaPreview([]);
              }}
            >
              Reset
            </button>

            <button
              type="submit"
              className="inline-flex items-center justify-center rounded-lg bg-[#397e50] px-4 py-2 text-sm font-semibold text-white shadow-sm transition hover:bg-[#2f6a43]"
            >
              Simpan Ujian
            </button>
          </div>
        </form>
      </div>
    </div>
  );
};

export default BuatUjianForm;
