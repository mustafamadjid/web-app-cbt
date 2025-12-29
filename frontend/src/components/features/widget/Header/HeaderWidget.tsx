import { User } from "lucide-react";



export const Header = () => {
  return (
    <>
      <header className="px-4 py-4 shrink-0">
        <div className="flex justify-between items-center-safe bg-white w-full h-full p-5 rounded-xl shadow-md">
          <h1 className="font-bold text-xl">Dashboard</h1>

          {/*Profil */}
          <div className="flex items-center gap-3">
            <div className="flex flex-col items-end">
              <h2 className="font-semibold text-md mb-1">
                Athaullah Mustafa Madjid
              </h2>

              <h3 className="bg-amber-100 px-3 py-1 rounded-md text-sm">
                Wali Kelas XI
              </h3>
            </div>

            <div className="relative">
              <div className="bg-amber-100 w-12 h-12 rounded-full flex items-center justify-center">
                {/* <img src="" alt="" /> */}
                <User />
              </div>
              <span className="absolute bottom-0 right-0 h-3 w-3 rounded-full ring-2 ring-white bg-green-500" />
            </div>
          </div>
        </div>
      </header>
    </>
  );
};
