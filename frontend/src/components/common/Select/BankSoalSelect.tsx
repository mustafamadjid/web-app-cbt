import type { BankSoalItem } from "@/types/BankSoal/BankSoal";

type BankSoalSelectProps = {
  id: string;
  label: string;
  value: number;
  options: BankSoalItem[];
  onChange: (value: number) => void;
  onBlur?: () => void;
  loading?: boolean;
  disabled?: boolean;
  error?: boolean;
  placeholder?: string;
  loadingPlaceholder?: string;
};

const selectBaseClass =
  "w-full rounded-lg border border-slate-200 bg-white px-3 py-2 text-sm outline-none transition focus:border-[#397e50] focus:ring-2 focus:ring-[#397e50] disabled:bg-slate-50 disabled:text-slate-500";

const BankSoalSelect = ({
  id,
  label,
  value,
  options,
  onChange,
  onBlur,
  loading = false,
  disabled = false,
  error = false,
  placeholder = "Pilih bank soal",
  loadingPlaceholder = "Memuat bank soal...",
}: BankSoalSelectProps) => {
  return (
    <div>
      <label htmlFor={id} className="text-xs font-medium text-slate-600">
        {label}
      </label>
      <select
        id={id}
        value={value}
        onChange={(event) => onChange(Number(event.target.value))}
        onBlur={onBlur}
        disabled={disabled || loading}
        className={`${selectBaseClass} ${error ? "border-rose-300 ring-rose-100" : ""}`}
      >
        <option value={0}>{loading ? loadingPlaceholder : placeholder}</option>
        {options.map((item) => (
          <option key={item.id_bank_soal} value={item.id_bank_soal}>
            {item.nama_bank_soal}
          </option>
        ))}
      </select>
    </div>
  );
};

export default BankSoalSelect;
