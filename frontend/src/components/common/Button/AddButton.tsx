import { Plus } from "lucide-react";

type AddButtonProps = {
  label?: string;
  onClick: () => void;
};

export const AddButton = ({ label = "Tambah", onClick }: AddButtonProps) => {
  return (
    <>
      <button
        className="
      flex items-center gap-2 cursor-pointer text-sm
      px-4 py-2 rounded-md bg-[#397e50] text-white font-bold 
      transition duration-200 hover:bg-white hover:text-black 
      border-2 border-transparent hover:border-[#397e50]"
        onClick={onClick}
      >
        <Plus />
        {label}
      </button>
    </>
  );
};
