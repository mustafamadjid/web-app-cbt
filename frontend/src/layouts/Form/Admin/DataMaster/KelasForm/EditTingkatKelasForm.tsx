import { useMemo } from "react";

import { useGetDataKelasFull } from "@/services/Api/features-api/DataMaster/kelas.service";
import type { KelasFormValues, TingkatKelas } from "@/types/DataMaster/Kelas";

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
  const { data: kelasData } = useGetDataKelasFull();
  const opsiTingkat: TingkatKelas[] = useMemo(
    () => kelasData?.item_tingkat_kelas ?? [],
    [kelasData?.item_tingkat_kelas],
  );

  return (
    <div>
      <label
        htmlFor="tingkat_kelas"
        className="text-xs font-medium text-slate-600"
      >
        Tingkat Kelas
      </label>

      <select
        id="tingkat_kelas"
        className={`w-full rounded-lg border border-slate-200 bg-white px-3 py-2 text-sm outline-none transition focus:border-[#397e50] focus:ring-2 focus:ring-[#397e50] disabled:bg-slate-50 disabled:text-slate-500 ${
          error ? "border-rose-300 ring-rose-100" : ""
        }`}
        value={value}
        onChange={(e) => {
          const selectedValue = e.target.value;
          onChange(selectedValue === "" ? "" : Number(selectedValue));
        }}
        onBlur={onBlur}
        disabled={disabled}
        required
      >
        <option value="">Pilih tingkat kelas</option>
        {opsiTingkat.map((item) => (
          <option key={item.id_tingkat_kelas} value={item.tingkat_kelas}>
            {item.tingkat_kelas}
          </option>
        ))}
      </select>

      {error && <p className="mt-1 text-xs text-rose-500">{error}</p>}
    </div>
  );
};

export default EditTingkatKelasForm;
