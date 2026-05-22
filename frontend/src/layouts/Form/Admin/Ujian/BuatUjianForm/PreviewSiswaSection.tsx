import FormSection from "./FormSection";
import type { BuatUjianFormController } from "./useBuatUjianForm";

type PreviewSiswaSectionProps = Pick<
  BuatUjianFormController,
  "loadingSiswa" | "siswaPreview" | "siswaPreviewEnabled" | "values"
>;

const PreviewSiswaSection = ({
  loadingSiswa,
  siswaPreview,
  siswaPreviewEnabled,
  values,
}: PreviewSiswaSectionProps) => (
  <FormSection
    title="Preview Daftar Siswa"
    description={
      values.kelas_scope === "SPESIFIK"
        ? "Preview mengambil id_tingkat_kelas dan id_nama_kelas."
        : "."
    }
    action={
      <span className="rounded-full bg-slate-100 px-3 py-1 text-xs font-medium text-slate-600">
        {siswaPreview.length} siswa
      </span>
    }
  >
    {!siswaPreviewEnabled && (
      <p className="text-sm text-slate-500">
        Pilih tingkat kelas, dan jika scope spesifik pilih nama kelas, untuk
        melihat siswa.
      </p>
    )}

    {siswaPreviewEnabled && loadingSiswa && (
      <p className="text-sm text-slate-500">Memuat data siswa...</p>
    )}

    {siswaPreviewEnabled && !loadingSiswa && siswaPreview.length === 0 && (
      <p className="text-sm text-slate-500">
        Belum ada siswa pada filter kelas ini.
      </p>
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
                <td className="px-3 py-2 font-medium text-slate-700">
                  {siswa.nama_lengkap}
                </td>
                <td className="px-3 py-2">{siswa.no_absen}</td>
                <td className="px-3 py-2">{siswa.nama_kelas}</td>
                <td className="px-3 py-2 capitalize">{siswa.status_akun}</td>
              </tr>
            ))}
          </tbody>
        </table>
      </div>
    )}
  </FormSection>
);

export default PreviewSiswaSection;
