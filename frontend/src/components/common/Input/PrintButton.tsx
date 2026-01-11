import { Printer } from "lucide-react";

type PrintButtonProps = {
  label: string;
  onClick?: () => void;
  variant?: "primary" | "outline";
  className?: string;
};

const variantStyles: Record<NonNullable<PrintButtonProps["variant"]>, string> = {
  primary:
    "bg-[#397e50] text-white hover:bg-[#2f6842] focus-visible:ring-[#397e50]/30",
  outline:
    "border border-[#397e50] text-[#397e50] hover:bg-[#397e50]/10 focus-visible:ring-[#397e50]/20",
};

export const PrintButton = ({
  label,
  onClick,
  variant = "primary",
  className = "",
}: PrintButtonProps) => {
  return (
    <button
      type="button"
      onClick={onClick}
      className={`cursor-pointer inline-flex items-center justify-center gap-2 rounded-xl px-4 py-2 text-sm font-semibold transition focus-visible:outline-none focus-visible:ring-4 ${variantStyles[variant]} ${className}`}
    >
      <Printer size={16} />
      {label}
    </button>
  );
};
