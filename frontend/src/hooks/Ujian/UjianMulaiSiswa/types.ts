import type { SoalPreviewItem } from "@/types/Ujian/SoalPreview";
import type {
  JawabanUjianSiswaItem,
  JawabanUjianSiswaResponse,
  SaveJawabanUjianSiswaRequest,
} from "@/types/Ujian/ujianSiswa";

export type UseUjianMulaiSiswaControllerParams = {
  resetKey: number;
  attemptId: number | null;
  soalPreview: SoalPreviewItem[];
  jawabanUjianData: JawabanUjianSiswaResponse | null;
  loadingJawabanUjian: boolean;
  refetchJawabanUjian: () => Promise<JawabanUjianSiswaResponse | null>;
  executeSaveJawaban: (
    payload: SaveJawabanUjianSiswaRequest,
  ) => Promise<boolean>;
  savingJawaban: boolean;
  saveJawabanError: string | null;
  resetSaveJawabanState: () => void;
  executeSubmitUjian: (attemptId: number) => Promise<boolean>;
  submittingAttempt: boolean;
  submitAttemptError: string | null;
  resetSubmitAttemptState: () => void;
  expiringAttempt: boolean;
  sessionExitError: string | null;
  clearSessionExitError: () => void;
  hasActiveExamSession: boolean;
  allowNavigation: () => void;
  handleBrowserViolation: () => Promise<boolean>;
  closeExamWindowOrLeave: () => void;
  clearAttemptCache: () => void;
  clearSoalCache: () => void;
  navigateToResults: () => void;
};

export type UseUjianMulaiSiswaControllerResult = {
  currentIndex: number;
  currentSoal: SoalPreviewItem | null;
  selectedOptions: Record<number, number>;
  essayAnswers: Record<number, string>;
  previewJawabanItems: JawabanUjianSiswaItem[];
  preparingPreview: boolean;
  isSubmitPreviewMode: boolean;
  isManualSubmitBusy: boolean;
  saveCurrentJawaban: (soalId: number | null) => Promise<boolean>;
  handleSelectOption: (soalId: number, optionId: number) => void;
  handleEssayAnswerChange: (soalId: number, value: string) => void;
  handleNavigateQuestion: (nextIndex: number) => void;
  handleOpenSubmitPreview: () => void;
  handleBackToQuestionMode: () => void;
  handleFinalSubmit: () => void;
};
