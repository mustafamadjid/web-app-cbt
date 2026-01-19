export type SiswaProfile = {
  id: number;
  nama: string;
  kelas: string;
  tingkat_kelas_id: number;
};

export type SiswaSemesterAverage = {
  semester: string;
  rata_rata: number;
  target?: number;
};

export type SiswaDashboardSummary = {
  ujian_selesai: number;
  total_ujian: number;
  rata_rata_semester: SiswaSemesterAverage[];
};
