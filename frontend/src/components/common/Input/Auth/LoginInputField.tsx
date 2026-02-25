import React from "react";

type LoginInputFieldProps = {
  id: string;
  label: string;
  type?: React.HTMLInputTypeAttribute;
  value: string;
  onChange: (value: string) => void;

  autoComplete?: string;
  disabled?: boolean;

  required?: boolean;
};

const LoginInputField = ({
  id,
  label,
  type = "text",
  value,
  onChange,

  autoComplete,
  disabled,

  required = true,
}: LoginInputFieldProps) => {
  return (
    <div className="relative">
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
          peer w-full
          border-gray-300  rounded-xl
          text-[1rem] bg-white
          p-2 pt-4
          transition-[border-color] duration-150 ease-in-out
          focus:outline-none focus-visible:ring-2 focus-visible:ring-[#397e50]
          disabled:opacity-60 disabled:cursor-not-allowed
        "
      />

      <label
        htmlFor={id}
        className="
          absolute left-[15px] top-1/2 -translate-y-1/2
          text-black pointer-events-none
          transition-all duration-150 ease-in-out
          peer-focus:top-0 peer-focus:-translate-y-1/2 peer-focus:scale-[0.8]
          peer-focus:bg-white peer-focus:rounded-2xl peer-focus:px-[0.5em] peer-focus:text-black
          peer-valid:top-0 peer-valid:-translate-y-1/2 peer-valid:scale-[0.8]
          peer-valid:bg-white peer-valid:px-[0.5em] peer-valid:text-black
        "
      >
        {label}
      </label>
    </div>
  );
};

export default LoginInputField;
