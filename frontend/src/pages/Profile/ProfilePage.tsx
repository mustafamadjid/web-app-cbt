import { useEffect, useMemo, useState } from "react";
import { User } from "lucide-react";
import { useAuth } from "@/contexts/AuthContext";
import {
  getProfileByUserId,
  type ProfileData,
} from "@/services/Api/features-api/Profile/profile.service";

type LoadState = "loading" | "success" | "error";

type ProfileField = {
  label: string;
  value: string;
};

const roleLabels: Record<string, string> = {
  ADMIN: "Administrator",
  GURU: "Guru",
  SISWA: "Siswa",
};

const formatLabel = (value?: string | number | null) => {
  if (value === null || value === undefined || value === "") return "-";
  return String(value)
    .toLowerCase()
    .replace(/_/g, " ")
    .replace(/\b\w/g, (match) => match.toUpperCase());
};

const getInitials = (name?: string) => {
  if (!name) return "U";
  const parts = name.trim().split(/\s+/);
  const first = parts[0]?.[0] ?? "";
  const last = parts.length > 1 ? parts[parts.length - 1]?.[0] ?? "" : "";
  return `${first}${last}`.toUpperCase();
};

const ProfilePage = () => {
  const { user } = useAuth();
  const [profile, setProfile] = useState<ProfileData | null>(null);
  const [state, setState] = useState<LoadState>("loading");
  const [errorMessage, setErrorMessage] = useState<string | null>(null);

  useEffect(() => {
    if (!user?.id) {
      setState("error");
      setErrorMessage("Data profil belum tersedia.");
      return;
    }

    let isActive = true;
    setState("loading");
    void getProfileByUserId(user.id)
      .then((data) => {
        if (!isActive) return;
        setProfile(data);
        setState("success");
        setErrorMessage(null);
      })
      .catch((error: Error) => {
        if (!isActive) return;
        setProfile(null);
        setState("error");
        setErrorMessage(error.message || "Gagal memuat profil.");
      });

    return () => {
      isActive = false;
    };
  }, [user?.id]);

  const displayName = profile?.namaLengkap ?? "-";
  const displayRole = profile?.role ? roleLabels[profile.role] ?? profile.role : "-";
  const avatarUrl = profile?.urlGambarProfil ?? null;

  const personalFields = useMemo<ProfileField[]>(
    () => [
      { label: "Nama Lengkap", value: formatLabel(profile?.namaLengkap) },
      { label: "Username", value: formatLabel(profile?.username) },
      { label: "Email", value: formatLabel(profile?.email) },
      { label: "Nomor HP", value: formatLabel(profile?.noHp) },
      { label: "Jenis Kelamin", value: formatLabel(profile?.jenisKelamin) },
      { label: "Status Akun", value: formatLabel(profile?.statusAkun) },
    ],
    [profile]
  );

  const roleFields = useMemo<ProfileField[]>(() => {
    if (!profile) return [];

    if ("nip" in profile) {
      return [
        { label: "NIP", value: formatLabel(profile.nip) },
        { label: "Jabatan", value: formatLabel(profile.jabatan) },
        { label: "Bidang Studi", value: formatLabel(profile.bidangStudi) },
      ];
    }

    return [
      { label: "Nomor Absen", value: formatLabel(profile.noAbsen) },
      { label: "Angkatan", value: formatLabel(profile.angkatan) },
      { label: "Tempat Lahir", value: formatLabel(profile.tempatLahir) },
      { label: "Tanggal Lahir", value: formatLabel(profile.tanggalLahir) },
      { label: "Tingkat Kelas", value: formatLabel(profile.id_tingkat_kelas) },
      { label: "Nama Kelas", value: formatLabel(profile.id_nama_kelas) },
    ];
  }, [profile]);

  return (
    <div className="px-6 py-8 text-slate-800">
      <div className="mx-auto flex w-full max-w-5xl flex-col gap-6">
        <div className="rounded-2xl border border-slate-200 bg-white p-6 shadow-sm transition-all duration-300 hover:-translate-y-0.5 hover:shadow-md motion-safe:animate-[profileFade_0.45s_ease-out]">
          <div className="flex flex-col items-start justify-between gap-6 sm:flex-row sm:items-center">
            <div>
              <h1 className="text-2xl font-bold text-slate-900">
                Informasi Profil
              </h1>
              <p className="mt-1 text-sm text-slate-500">
                Data profil ditampilkan berdasarkan peran akun.
              </p>
            </div>
            <span className="inline-flex items-center rounded-full border border-[#397e50]/30 bg-[#397e50]/10 px-3 py-1 text-xs font-semibold uppercase tracking-wide text-[#397e50]">
              {displayRole}
            </span>
          </div>

          <div className="mt-6 flex flex-col gap-4 sm:flex-row sm:items-center">
            <div className="relative h-20 w-20 overflow-hidden rounded-full border border-slate-200 shadow-sm">
              {avatarUrl ? (
                <img
                  src={avatarUrl}
                  alt={displayName}
                  className="h-full w-full object-cover"
                />
              ) : (
                <div className="flex h-full w-full items-center justify-center bg-[#397e50]/10 text-[#397e50]">
                  <span className="text-lg font-semibold">
                    {getInitials(displayName)}
                  </span>
                </div>
              )}
              {!avatarUrl && (
                <div className="absolute -bottom-2 -right-2 rounded-full bg-white p-1 text-[#397e50] shadow-sm">
                  <User className="h-4 w-4" />
                </div>
              )}
            </div>
            <div>
              <h2 className="text-lg font-semibold text-slate-900">
                {displayName}
              </h2>
              <p className="text-sm text-slate-500">{displayRole}</p>
            </div>
          </div>
        </div>

        <div className="rounded-2xl border border-slate-200 bg-white p-6 shadow-sm transition-all duration-300 hover:-translate-y-0.5 hover:shadow-md motion-safe:animate-[profileFade_0.55s_ease-out]">
          <h3 className="text-lg font-semibold text-slate-900">
            Detail Personal
          </h3>
          <div className="mt-4 grid gap-4 sm:grid-cols-2">
            {personalFields.map((field) => (
              <ProfileFieldItem key={field.label} {...field} />
            ))}
          </div>
        </div>

        <div className="rounded-2xl border border-slate-200 bg-white p-6 shadow-sm transition-all duration-300 hover:-translate-y-0.5 hover:shadow-md motion-safe:animate-[profileFade_0.65s_ease-out]">
          <h3 className="text-lg font-semibold text-slate-900">
            Informasi Role
          </h3>
          <div className="mt-4 grid gap-4 sm:grid-cols-2">
            {roleFields.map((field) => (
              <ProfileFieldItem key={field.label} {...field} />
            ))}
          </div>
        </div>

        {state === "loading" && (
          <div className="rounded-2xl border border-dashed border-slate-200 bg-white/70 p-6 text-sm text-slate-500 shadow-sm">
            Memuat data profil...
          </div>
        )}

        {state === "error" && (
          <div className="rounded-2xl border border-red-200 bg-red-50 p-6 text-sm text-red-700 shadow-sm">
            {errorMessage ?? "Gagal memuat profil."}
          </div>
        )}
      </div>
    </div>
  );
};

const ProfileFieldItem = ({ label, value }: ProfileField) => (
  <div className="space-y-2">
    <p className="text-xs font-semibold uppercase tracking-wide text-slate-500">
      {label}
    </p>
    <div className="rounded-xl border border-slate-200 bg-white px-4 py-3 text-sm font-medium text-slate-800 shadow-sm transition-colors duration-200 hover:border-[#397e50]/50">
      {value}
    </div>
  </div>
);

export default ProfilePage;
