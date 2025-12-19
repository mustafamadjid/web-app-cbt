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

export const LoginInputField = ({
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
          peer w-[50%]
          border-[1.5px] border-green-800 rounded-md
          text-[1rem] bg-transparent
          px-3 py-2 pt-4
          transition-[border-color] duration-150 ease-in-out
          focus:outline-none focus:border-green-500
          valid:border-green-500
          disabled:opacity-60 disabled:cursor-not-allowed
        "
      />

      <label
        htmlFor={id}
        className="
          absolute left-[15px] top-1/2 -translate-y-1/2
          text-green-800 pointer-events-none
          transition-all duration-150 ease-in-out
          peer-focus:top-0 peer-focus:-translate-y-1/2 peer-focus:scale-[0.8]
          peer-focus:bg-white peer-focus:px-[0.2em] peer-focus:text-green-500
          peer-valid:top-0 peer-valid:-translate-y-1/2 peer-valid:scale-[0.8]
          peer-valid:bg-white peer-valid:px-[0.2em] peer-valid:text-green-500
        "
      >
        {label}
      </label>
    </div>
  );
};
