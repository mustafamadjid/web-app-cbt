export type KoreksiEssayItem = {
  id_jawaban: number;
  essay_is_benar: boolean;
};

export type SubmitKoreksiEssayRequest = {
  jawaban: KoreksiEssayItem[];
};
