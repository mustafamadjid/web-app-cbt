import InputField from "@/components/common/Input/InputField";

import FormSection from "./FormSection";
import type { BuatUjianFormController } from "./useBuatUjianForm";

type InformasiUjianSectionProps = Pick<
  BuatUjianFormController,
  "errors" | "hasError" | "onBlur" | "setField" | "values"
>;

const InformasiUjianSection = ({
  errors,
  hasError,
  onBlur,
  setField,
  values,
}: InformasiUjianSectionProps) => (
  <FormSection title="Informasi Ujian" description="Isi nama dan deskripsi ujian.">
    <div className="grid grid-cols-1 gap-4 md:grid-cols-2">
      <div className="md:col-span-2">
        <InputField
          id="nama_ujian"
          label="Nama Ujian"
          value={values.nama_ujian}
          onChange={(value) => setField("nama_ujian", value)}
          onBlur={() => onBlur("nama_ujian")}
          placeholder="Contoh: Ujian Tengah Semester Matematika"
          inputClassName={
            hasError("nama_ujian") ? "border-rose-300 ring-rose-100" : ""
          }
        />
        {hasError("nama_ujian") && (
          <p className="mt-1 text-xs text-rose-600">{errors.nama_ujian}</p>
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
          value={values.deskripsi_ujian}
          onChange={(event) => setField("deskripsi_ujian", event.target.value)}
          onBlur={() => onBlur("deskripsi_ujian")}
          placeholder="Jelaskan cakupan materi ujian."
        />
        {hasError("deskripsi_ujian") && (
          <p className="mt-1 text-xs text-rose-600">
            {errors.deskripsi_ujian}
          </p>
        )}
      </div>
    </div>
  </FormSection>
);

export default InformasiUjianSection;
