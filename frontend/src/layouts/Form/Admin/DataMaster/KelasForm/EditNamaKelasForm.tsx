import InputField from "@/components/common/Input/InputField";
import type { KelasFormValues } from "@/types/DataMaster/Kelas";

type EditNamaKelasFormProps = {
  value: KelasFormValues["nama_kelas"];
  error?: string;
  disabled?: boolean;
  onChange: (value: KelasFormValues["nama_kelas"]) => void;
  onBlur: () => void;
};

const EditNamaKelasForm = ({
  value,
  error,
  disabled = false,
  onChange,
  onBlur,
}: EditNamaKelasFormProps) => {
  return (
    <div>
      <InputField
        id="nama_kelas"
        label="Nama Kelas"
        value={value}
        onChange={onChange}
        onBlur={onBlur}
        placeholder="Contoh: X IPA 1"
        inputClassName={error ? "border-rose-300 ring-rose-100" : ""}
        disabled={disabled}
        required
      />
      {error && <p className="mt-1 text-xs text-rose-500">{error}</p>}
    </div>
  );
};

export default EditNamaKelasForm;
