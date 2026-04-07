import type { SoalPreviewItem } from "@/types/Ujian/SoalPreview";
import type {
  JawabanUjianSiswaResponse,
  SaveJawabanUjianSiswaRequest,
} from "@/types/Ujian/ujianSiswa";

export type JawabanDraft = {
  id_soal: number;
  id_pilihan: number | null;
  jawaban_essay: string;
  waktu_jawab: string | null;
  isDirty: boolean;
};

export type JawabanDraftMap = Record<number, JawabanDraft>;

export const createJawabanDraft = (idSoal: number): JawabanDraft => ({
  id_soal: idSoal,
  id_pilihan: null,
  jawaban_essay: "",
  waktu_jawab: null,
  isDirty: false,
});

export const buildSaveJawabanPayload = (
  idAttempt: number,
  jawaban: JawabanDraft,
): SaveJawabanUjianSiswaRequest => {
  const normalizedEssay = jawaban.jawaban_essay.trim();

  return {
    id_attempt: idAttempt,
    jawaban: [
      {
        id_soal: jawaban.id_soal,
        id_pilihan: jawaban.id_pilihan,
        jawaban_essay: jawaban.id_pilihan === null ? normalizedEssay || null : null,
        waktu_jawab: jawaban.waktu_jawab ?? new Date().toISOString(),
      },
    ],
  };
};

export const buildSelectedOptionsMap = (
  draftAnswers: JawabanDraftMap,
): Record<number, number> => {
  const mapped: Record<number, number> = {};

  for (const jawaban of Object.values(draftAnswers)) {
    if (jawaban.id_pilihan !== null) {
      mapped[jawaban.id_soal] = jawaban.id_pilihan;
    }
  }

  return mapped;
};

export const buildEssayAnswersMap = (
  draftAnswers: JawabanDraftMap,
): Record<number, string> => {
  const mapped: Record<number, string> = {};

  for (const jawaban of Object.values(draftAnswers)) {
    if (jawaban.jawaban_essay !== "") {
      mapped[jawaban.id_soal] = jawaban.jawaban_essay;
    }
  }

  return mapped;
};

type MergeServerJawabanIntoDraftsParams = {
  draftAnswers: JawabanDraftMap;
  jawabanUjianData: JawabanUjianSiswaResponse;
  soalPreview: SoalPreviewItem[];
};

export const mergeServerJawabanIntoDrafts = ({
  draftAnswers,
  jawabanUjianData,
  soalPreview,
}: MergeServerJawabanIntoDraftsParams): JawabanDraftMap => {
  const next = { ...draftAnswers };
  const serverJawabanBySoal = new Map(
    jawabanUjianData.jawaban.map((item) => [item.id_soal, item]),
  );

  for (const soal of soalPreview) {
    const existing = draftAnswers[soal.id];
    if (existing?.isDirty) {
      continue;
    }

    const serverJawaban = serverJawabanBySoal.get(soal.id);
    if (!serverJawaban) {
      delete next[soal.id];
      continue;
    }

    next[soal.id] = {
      id_soal: serverJawaban.id_soal,
      id_pilihan: serverJawaban.id_pilihan,
      jawaban_essay: serverJawaban.jawaban_essay ?? "",
      waktu_jawab: serverJawaban.waktu_jawab,
      isDirty: false,
    };
  }

  return next;
};

type ApplySavedJawabanToDraftsParams = {
  draftAnswers: JawabanDraftMap;
  soalId: number;
  payload: SaveJawabanUjianSiswaRequest;
};

export const applySavedJawabanToDrafts = ({
  draftAnswers,
  soalId,
  payload,
}: ApplySavedJawabanToDraftsParams): JawabanDraftMap => {
  const latestJawaban = draftAnswers[soalId];
  if (!latestJawaban) {
    return draftAnswers;
  }

  const [savedItem] = payload.jawaban;
  const shouldClearAnswer =
    savedItem.id_pilihan === null && savedItem.jawaban_essay === null;

  if (shouldClearAnswer) {
    const rest = { ...draftAnswers };
    delete rest[soalId];
    return rest;
  }

  return {
    ...draftAnswers,
    [soalId]: {
      ...latestJawaban,
      id_pilihan: savedItem.id_pilihan,
      jawaban_essay: savedItem.jawaban_essay ?? "",
      waktu_jawab: savedItem.waktu_jawab,
      isDirty: false,
    },
  };
};
