import InputField from "@/components/common/Input/InputField";

import { selectBaseClass } from "./constants";
import FormSection from "./FormSection";
import type { BuatUjianFormController } from "./useBuatUjianForm";

type KeamananTokenSectionProps = Pick<
  BuatUjianFormController,
  "errors" | "hasError" | "onBlur" | "setField" | "submitting" | "values"
>;

const KeamananTokenSection = ({
  errors,
  hasError,
  onBlur,
  setField,
  submitting,
  values,
}: KeamananTokenSectionProps) => (
  <FormSection title="Keamanan & Token" description="Token maksimal 30 karakter.">
    <div className="grid grid-cols-1 gap-4 md:grid-cols-2">
      <div>
        <label htmlFor="acak_soal" className="text-xs font-medium text-slate-600">
          Acak Soal
        </label>
        <select
          id="acak_soal"
          value={values.acak_soal ? "ya" : "tidak"}
          onChange={(event) =>
            setField("acak_soal", event.target.value === "ya")
          }
          disabled={submitting}
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
          inputClassName={
            hasError("token") ? "border-rose-300 ring-rose-100" : ""
          }
        />
        <p className="mt-1 text-[11px] text-slate-500">
          {values.token.length}/30 karakter
        </p>
        {hasError("token") && (
          <p className="mt-1 text-xs text-rose-600">{errors.token}</p>
        )}
      </div>
    </div>
  </FormSection>
);

export default KeamananTokenSection;
