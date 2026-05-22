import DatePicker from "@/components/common/DateInput/DatePicker";
import TimePicker from "@/components/common/DateInput/TimePicker";
import InputField from "@/components/common/Input/InputField";

import FormSection from "./FormSection";
import type { BuatUjianFormController } from "./useBuatUjianForm";

type JadwalUjianSectionProps = Pick<
  BuatUjianFormController,
  "durasiMenit" | "errors" | "hasError" | "onBlur" | "setField" | "values"
>;

const JadwalUjianSection = ({
  durasiMenit,
  errors,
  hasError,
  onBlur,
  setField,
  values,
}: JadwalUjianSectionProps) => (
  <FormSection title="Jadwal Ujian" description="">
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
        <p className="mt-1 text-[11px] text-slate-500">
          Waktu ujian memakai WIB (24 jam).
        </p>
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
        <p className="mt-1 text-[11px] text-slate-500">
          Waktu ujian memakai WIB (24 jam).
        </p>
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
  </FormSection>
);

export default JadwalUjianSection;
