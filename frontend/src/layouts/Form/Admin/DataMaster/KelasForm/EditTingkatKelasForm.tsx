import type { KelasFormValues } from "@/types/DataMaster/Kelas";

type EditTingkatKelasFormProps = {
  value: KelasFormValues["tingkat_kelas"];
  error?: string;
  disabled?: boolean;
  onChange: (value: KelasFormValues["tingkat_kelas"]) => void;
  onBlur: () => void;
};

const EditTingkatKelasForm = ({
  value,
  error,
  disabled = false,
  onChange,
  onBlur,
}: EditTingkatKelasFormProps) => {
  return (
    <div>
      <label
        htmlFor="tingkat_kelas"
        className="text-xs font-medium text-slate-600"
      >
        Tingkat Kelas
      </label>

      <input
        id="tingkat_kelas"
        type="number"
        inputMode="numeric"
        className={`w-full rounded-lg border border-slate-200 bg-white px-3 py-2 text-sm outline-none transition focus:border-[#397e50] focus:ring-2 focus:ring-[#397e50] disabled:bg-slate-50 disabled:text-slate-500 ${
          error ? "border-rose-300 ring-rose-100" : ""
        }`}
        placeholder="Contoh: 10"
        value={value}
        onChange={(e) => {
          const raw = e.target.value;
          onChange(raw === "" ? "" : Number(raw));
        }}
        onBlur={onBlur}
        min={1}
        step={1}
        disabled={disabled}
        required
      />

      {error && <p className="mt-1 text-xs text-rose-500">{error}</p>}
    </div>
  );
};

export default EditTingkatKelasForm;
