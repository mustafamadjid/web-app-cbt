type InputFieldProps = {
  id: string;
  label: string;
  type?: React.HTMLInputTypeAttribute;
  value: string;
  onChange: (value: string) => void;

  autoComplete?: string;
  disabled?: boolean;

  required?: boolean;
};

export const InputField = ({
  id,
  label,
  type = "text",
  value,
  onChange,

  autoComplete,
  disabled,

  required = true,
}: InputFieldProps) => {
  return (
    <div className="relative">
      <label
        htmlFor={id}
        className="
              absolute left-[15px] top-1/2 -translate-y-1/2
              text-black pointer-events-none
              transition-all duration-150 ease-in-out
            "
      >
        {label}
      </label>
      <input
        id={id}
        type={type}
        value={value}
        required={required}
        autoComplete={autoComplete}
        disabled={disabled}
        onChange={(e: React.ChangeEvent<HTMLInputElement>) =>
          onChange(e.target.value)
        }
        className="
            w-full
              border-gray-100  rounded-xl
              text-[1rem] bg-white
              p-2 pt-4
              transition-[border-color] duration-150 ease-in-out
              focus:outline-none focus-visible:ring-2 focus-visible:ring-green-500
              disabled:opacity-60 disabled:cursor-not-allowed
            "
      />
    </div>
  );
};