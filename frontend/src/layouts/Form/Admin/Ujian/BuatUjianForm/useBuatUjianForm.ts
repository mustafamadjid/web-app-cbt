import { useEffect, useMemo, useRef, useState } from "react";

import { calculateDuration } from "@/helper/CalculateDuration/calculateDuration";
import { createSetField } from "@/helper/setField/setField";
import { ApiError } from "@/services/Api/api";
import { getUserFriendlyErrorMessage } from "@/services/Api/errorMessage";
import { useGetBankSoal } from "@/services/Api/features-api/BankSoal/banksoal.service";
import { useGetDataKelasFull } from "@/services/Api/features-api/DataMaster/kelas.service";
import { useGetMapel } from "@/services/Api/features-api/DataMaster/mapel.service";
import { useGetRuangUjian } from "@/services/Api/features-api/DataMaster/ruang-ujian.service";
import { useGetSesi } from "@/services/Api/features-api/DataMaster/sesi.service";
import { useGetAllGuru } from "@/services/Api/features-api/KelolaAkun/akunguru.service";
import { useGetListSiswa } from "@/services/Api/features-api/KelolaAkun/akunsiswa.service";
import type { BankSoalItem } from "@/types/BankSoal/BankSoal";
import type {
  FullDataKelas,
  NamaKelas,
  TingkatKelas,
} from "@/types/DataMaster/Kelas";
import type { MataPelajaranRow } from "@/types/DataMaster/MataPelajaran";
import type { RuangUjianRow } from "@/types/DataMaster/RuangUjian";
import type { SesiRow } from "@/types/DataMaster/Sesi";
import type { DataGuru } from "@/types/KelolaAkun/AkunGuru";
import type { DataAkunSiswa } from "@/types/KelolaAkun/AkunSiswa";
import {
  EMPTY_BUAT_UJIAN_FORM_VALUES,
  type BuatUjianFormValues,
} from "@/types/Ujian/BuatUjian";

import { validateBuatUjianForm } from "./constants";

type UseBuatUjianFormParams = {
  onSubmit: (values: BuatUjianFormValues) => Promise<void>;
};

export const useBuatUjianForm = ({ onSubmit }: UseBuatUjianFormParams) => {
  const [values, setValues] = useState<BuatUjianFormValues>(
    EMPTY_BUAT_UJIAN_FORM_VALUES,
  );
  const [selectedMapelId, setSelectedMapelId] = useState(0);
  const [touched, setTouched] = useState<Record<string, boolean>>({});
  const [submitError, setSubmitError] = useState<string | null>(null);
  const [submitting, setSubmitting] = useState(false);
  const skipNextKelasResetRef = useRef(true);

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
    data: mapelData,
    loading: loadingMapel,
    error: mapelError,
  } = useGetMapel(
    {
      tingkatKelas: values.id_kelas > 0 ? values.id_kelas : undefined,
      limit: 100,
      offset: 0,
    },
    values.id_kelas > 0,
  );
  const mapelOptions: MataPelajaranRow[] = useMemo(
    () => mapelData ?? [],
    [mapelData],
  );

  const bankSoalFetchEnabled = values.id_kelas > 0 && selectedMapelId > 0;
  const {
    data: bankSoalData,
    loading: loadingBankSoal,
    error: bankSoalError,
  } = useGetBankSoal(
    {
      id_kelas: values.id_kelas > 0 ? values.id_kelas : undefined,
      id_mapel: selectedMapelId > 0 ? selectedMapelId : undefined,
      limit: 100,
      offset: 0,
    },
    bankSoalFetchEnabled,
  );
  const bankSoalOptions: BankSoalItem[] = useMemo(
    () => bankSoalData ?? [],
    [bankSoalData],
  );

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
    mapelError ||
    bankSoalError ||
    ruangError ||
    sesiError ||
    guruError ||
    siswaError;

  useEffect(() => {
    if (skipNextKelasResetRef.current) {
      skipNextKelasResetRef.current = false;
      return;
    }

    setValues((prev) => ({
      ...prev,
      id_nama_kelas: 0,
      id_bank_soal: 0,
    }));
    setSelectedMapelId(0);
  }, [values.id_kelas]);

  useEffect(() => {
    if (values.kelas_scope === "SEMUA" && values.id_nama_kelas !== 0) {
      setField("id_nama_kelas", 0);
    }
  }, [setField, values.kelas_scope, values.id_nama_kelas]);

  useEffect(() => {
    if (loadingBankSoal || values.id_bank_soal === 0) {
      return;
    }

    if (
      !bankSoalOptions.some((item) => item.id_bank_soal === values.id_bank_soal)
    ) {
      setField("id_bank_soal", 0);
    }
  }, [bankSoalOptions, loadingBankSoal, setField, values.id_bank_soal]);

  useEffect(() => {
    if (loadingMapel || selectedMapelId === 0) {
      return;
    }

    if (!mapelOptions.some((item) => item.id === selectedMapelId)) {
      setSelectedMapelId(0);
      setField("id_bank_soal", 0);
    }
  }, [loadingMapel, mapelOptions, selectedMapelId, setField]);

  const durasiMenit = useMemo(
    () => calculateDuration(values.waktu_mulai, values.waktu_selesai),
    [values.waktu_mulai, values.waktu_selesai],
  );

  const bankSoalPlaceholder =
    values.id_kelas === 0
      ? "Pilih tingkat kelas terlebih dahulu"
      : selectedMapelId === 0
        ? "Pilih mapel terlebih dahulu"
        : "Pilih bank soal";

  const errors = validateBuatUjianForm(values);
  const hasError = (name: keyof BuatUjianFormValues) =>
    Boolean(errors[name]) && Boolean(touched[name]);

  const handleReset = () => {
    setValues(EMPTY_BUAT_UJIAN_FORM_VALUES);
    setSelectedMapelId(0);
    setTouched({});
    setSubmitError(null);
    skipNextKelasResetRef.current = true;
  };

  const handleSubmit = async (event: React.FormEvent) => {
    event.preventDefault();
    setSubmitError(null);

    const nextTouched: Record<string, boolean> = {};
    Object.keys(values).forEach((key) => {
      nextTouched[key] = true;
    });
    setTouched(nextTouched);

    const currentErrors = validateBuatUjianForm(values);
    if (Object.keys(currentErrors).length > 0) {
      setSubmitError(
        "Periksa kembali input yang masih kosong atau belum valid.",
      );
      return;
    }

    setSubmitting(true);
    try {
      await onSubmit(values);
    } catch (error) {
      if (error instanceof ApiError) {
        setSubmitError(
          getUserFriendlyErrorMessage(error, {
            action: "create",
            entity: "ujian",
          }),
        );
      } else if (error instanceof Error) {
        setSubmitError(
          getUserFriendlyErrorMessage(error, {
            action: "create",
            entity: "ujian",
          }),
        );
      } else {
        setSubmitError("Terjadi kesalahan saat menyimpan ujian.");
      }
    } finally {
      setSubmitting(false);
    }
  };

  return {
    values,
    errors,
    setField,
    onBlur,
    hasError,
    selectedMapelId,
    setSelectedMapelId,
    submitting,
    submitError,
    loadErrorMessage,
    loadingKelas,
    loadingMapel,
    loadingBankSoal,
    loadingRuang,
    loadingSesi,
    loadingGuru,
    loadingSiswa,
    tingkatKelasOptions,
    namaKelasOptions,
    mapelOptions,
    bankSoalOptions,
    ruangOptions,
    sesiOptions,
    guruOptions,
    siswaPreviewEnabled,
    siswaPreview,
    durasiMenit,
    bankSoalPlaceholder,
    handleReset,
    handleSubmit,
  };
};

export type BuatUjianFormController = ReturnType<typeof useBuatUjianForm>;
