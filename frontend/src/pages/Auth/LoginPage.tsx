import { useNavigate } from "react-router";
import LoginForm from "../../layouts/Form/Auth/LoginForm";

const LoginPage = () => {
  const navigate = useNavigate();

  return (
    <div className="flex min-h-screen w-full bg-white">
      {/* KIRI: Area Form */}
      <div className="flex w-full flex-col justify-center bg-white px-6 py-12 lg:w-1/2 lg:px-20 xl:px-24">
        <div className="mx-auto w-full max-w-md">
          {/* Header / Logo */}
          <div className="mb-10 flex items-center gap-3">
            <div className="flex h-12 w-12 items-center justify-center rounded-xl">
              <img
                src="/Images/LoginPageImg/logo-fi.png"
                alt="Logo"
                className="h-8 w-auto object-contain"
              />
            </div>
            <span className="text-lg font-bold text-slate-800">
              SMA IT Fitrah Insani
            </span>
          </div>

          {/* Titles */}
          <div className="mb-8">
            <h1 className="text-3xl font-bold tracking-tight text-slate-900 sm:text-4xl">
              Selamat Datang
            </h1>
            <p className="mt-3 text-base text-slate-500">
              Silakan masukkan kredensial akun Anda untuk mengakses dashboard.
            </p>
          </div>

          {/* Form Container */}
          <div className="space-y-6">
            <LoginForm onSuccess={() => navigate("/dashboard")} />
          </div>

          {/* Footer Copyright (Opsional) */}
          <div className="mt-10 text-center text-xs text-slate-400 lg:text-left">
            &copy; {new Date().getFullYear()} SMA IT Fitrah Insani. All rights
            reserved.
          </div>
        </div>
      </div>

      {/* KANAN: Area Gambar (Hanya muncul di layar besar) */}
      <div className="hidden lg:block lg:w-1/2 lg:p-6">
        <div className="relative h-full w-full overflow-hidden rounded-[2rem] border-4 border-white shadow-2xl">
          {/* Gambar Background */}
          <img
            className="absolute inset-0 h-full w-full object-cover"
            src="/Images/LoginPageImg/pict-1.jpg"
            alt="School Activity"
          />

          {/* Overlay tanpa gradasi */}
          <div className="absolute inset-0 bg-black/35" />

          {/* Konten Dekoratif di atas Gambar */}
          <div className="absolute bottom-0 left-0 p-12 text-white">
            <blockquote className="max-w-md border-l-4 border-white/40 pl-4">
              <p className="text-lg font-medium leading-relaxed">
                "Pendidikan adalah senjata paling mematikan di dunia, karena
                dengan pendidikan, Anda dapat mengubah dunia."
              </p>
              <footer className="mt-4 text-sm font-semibold opacity-80">
                - Nelson Mandela
              </footer>
            </blockquote>
          </div>
        </div>
      </div>
    </div>
  );
};

export default LoginPage;
