type InputFieldProps = {
  id?: string;
  type?: string;
  value?: string;
  label?: string;
  onChange: (value: string) => void;
  onBlur?: () => void;
  inputClassName?: string;
  labelClassName?: string;
  placeholder?: string;
  autoComplete?: string;
  disabled?: boolean;
  required?: boolean;
};

export const InputField = ({
  id,
  type,
  value,
  onChange,
  onBlur,
  inputClassName,
  label,
  labelClassName,
  placeholder,
  autoComplete = "off",
  disabled = false,
  required = true,
}: InputFieldProps) => {
  return (
    <>
      <div>
        <label
          htmlFor={id}
          className={`
            text-xs font-medium text-slate-600 
            ${labelClassName}
            `}
        >
          {label}
        </label>
        <input
          type={type}
          id={id}
          value={value}
          onChange={(e: React.ChangeEvent<HTMLInputElement>) =>
            onChange(e.target.value)
          }
          onBlur={onBlur}
          className={`
            ${inputClassName} w-full 
            rounded-lg border border-slate-200 bg-white px-3 py-2 
            text-sm outline-none transition focus:border-[#397e50] 
            focus:ring-2 focus:ring-[#397e50] disabled:bg-slate-50 
            disabled:text-slate-500`}
          placeholder={placeholder}
          autoComplete={autoComplete}
          disabled={disabled}
          required={required}
        />
      </div>
    </>
  );
};
