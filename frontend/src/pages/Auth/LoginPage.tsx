import { useNavigate } from "react-router";
import { LoginForm } from "../../layouts/Form/Auth/LoginForm";

export const LoginPage = () => {
  const navigate = useNavigate();

  return (
    <>
      <section className="bg-gray-50 min-h-screen flex items-center justify-center">
        {/* Login Contrainer */}
        <div className="bg-gray-100 flex items-center rounded-2xl max-w-5xl shadow-lg  p-5">
          {/* Login Form */}
          <div className="flex flex-col gap-10 sm:w-1/2 px-16 mb-20">
            <div className="flex flex-col gap-8">
              <div className="flex text-sm gap-3 items-center">
                <div className="w-[45px]">
                  <img src="/Images/LoginPageImg/logo-fi.png" alt="" />
                </div>
                <h1 className="text-md font-semibold">SMA IT Fitrah Insani</h1>
              </div>

              <div>
                <h2 className="text-2xl font-bold">Halo,</h2>
                <h2 className="text-2xl font-bold">Selamat Datang</h2>{" "}
                <p className="text-sm opacity-80">
                  Silakan masuk menggunakan akun anda{" "}
                </p>{" "}
              </div>
            </div>

            <LoginForm onSuccess={() => navigate("/dashboard")} />
          </div>
          {/* Image */}
          <div className="sm:block hidden w-1/2 p-3">
            <img
              className="rounded-2xl "
              src="/Images/LoginPageImg/pict-1.jpg"
              alt=""
            />
          </div>
        </div>
      </section>
    </>
  );
};
