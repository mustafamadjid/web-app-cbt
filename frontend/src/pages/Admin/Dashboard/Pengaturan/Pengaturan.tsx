import { useMemo, useState } from "react";
import { PengaturanProfilForm } from "@/layouts/Form/Admin/PengaturanProfil/PengaturanProfilForm";

// import BackupRestoreForm from "..."; // contoh

type TabKey = "profil" | "backup";

type TabItem = {
  key: TabKey;
  label: string;
};

export const PengaturanProfil = () => {
  const [activeTab, setActiveTab] = useState<TabKey>("profil");

  const tabs: TabItem[] = useMemo(
    () => [
      { key: "profil", label: "Profil Sekolah" },
      { key: "backup", label: "Backup dan Restore" },
    ],
    []
  );

  return (
   
      <div className="py-10 flex flex-col">
        {/* Tabs */}
        <div className="px-8 flex gap-8">
          {tabs.map((t) => (
            <TabButton
              key={t.key}
              label={t.label}
              isActive={activeTab === t.key}
              onClick={() => setActiveTab(t.key)}
            />
          ))}
        </div>

        {/* Content */}
        <div className="px-8 py-4">
          {activeTab === "profil" && <PengaturanProfilForm />}
          {activeTab === "backup" && (
            <div className="mt-4">
              {/* ganti dengan komponen asli */}
              <p>Halaman Backup dan Restore</p>
            </div>
          )}
        </div>
      </div>
   
  );
};

type TabButtonProps = {
  label: string;
  isActive: boolean;
  onClick: () => void;
};

function TabButton({ label, isActive, onClick }: TabButtonProps) {
  return (
    <>
    <button
      type="button"
      onClick={onClick}
      className={[
        "relative pb-2 text-lg font-semibold transition-colors cursor-pointer",
        isActive ? "text-green-800" : "text-[#397e50] hover:text-green-800",
      ].join(" ")}
    >
      {label}
      <span
        className={[
          "absolute left-0 -bottom-0.5 h-0.5 bg-green-800 transition-all duration-200",
          isActive ? "w-full" : "w-0",
        ].join(" ")}
      />
    </button>
    </>
  );
}
