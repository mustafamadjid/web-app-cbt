import React, { useMemo, useState } from "react";

type StatusAkun = "aktif" | "nonaktif" | "dibekukan";
type JenisKelamin = "LAKI_LAKI" | "PEREMPUAN";

type BarisPengguna = {
  id: string;
  namaLengkap: string;
  email: string;
  username: string;
  nomorHp: string;
  jenisKelamin: JenisKelamin;
  statusAkun: StatusAkun;
  nip: string;
  jabatan: string;
  bidangStudi: string;
 urlGambarProfil: string;
};

const kelasTitikStatus: Record<StatusAkun, string> = {
  aktif: "bg-success",
  nonaktif: "bg-danger",
  dibekukan: "bg-neutral-tertiary-medium",
};

const labelStatus: Record<StatusAkun, string> = {
  aktif: "Aktif",
  nonaktif: "Nonaktif",
  dibekukan: "Dibekukan",
};

function samarkanNomorHp(nomorHp: string) {
  const digit = nomorHp.replace(/\s+/g, "");
  if (digit.length <= 4) return digit;
  const terlihat = digit.slice(-3);
  return `${digit.slice(0, 3)}****${terlihat}`;
}

function samarkanNip(nip: string) {
  const digit = nip.replace(/\s+/g, "");
  if (digit.length <= 6) return digit;
  const terlihat = digit.slice(-4);
  return `${digit.slice(0, 4)}****${terlihat}`;
}

export const AkunGuruTables: React.FC = () => {
  const [dropdownAksiTerbuka, setDropdownAksiTerbuka] = useState(false);
  const [kataKunci, setKataKunci] = useState("");
  const [idTerpilih, setIdTerpilih] = useState<Set<string>>(new Set());

  // Privasi: default ON (data sensitif disamarkan)
  const [samarkanDataSensitif, setSamarkanDataSensitif] = useState(true);

  const [daftarPengguna] = useState<BarisPengguna[]>([
    {
      id: "u1",
      namaLengkap: "Neil Sims",
      email: "neil.sims@flowbite.com",
      username: "neilsims",
      nomorHp: "081234567890",
      jenisKelamin: "LAKI_LAKI",
      statusAkun: "aktif",
      nip: "198701012010121001",
      jabatan: "React Developer",
      bidangStudi: "Informatika",
     urlGambarProfil: "/docs/images/people/profile-picture-1.jpg",
    },
    {
      id: "u2",
      namaLengkap: "Bonnie Green",
      email: "bonnie@flowbite.com",
      username: "bonnieg",
      nomorHp: "082233445566",
      jenisKelamin: "PEREMPUAN",
      statusAkun: "aktif",
      nip: "199002022011112002",
      jabatan: "Desainer",
      bidangStudi: "DKV",
     urlGambarProfil: "/docs/images/people/profile-picture-3.jpg",
    },
  ]);

  const penggunaTersaring = useMemo(() => {
    const q = kataKunci.trim().toLowerCase();
    if (!q) return daftarPengguna;

    return daftarPengguna.filter((p) => {
      return (
        p.namaLengkap.toLowerCase().includes(q) ||
        p.email.toLowerCase().includes(q) ||
        p.username.toLowerCase().includes(q) ||
        p.nomorHp.toLowerCase().includes(q) ||
        p.nip.toLowerCase().includes(q) ||
        p.jenisKelamin.toLowerCase().includes(q) ||
        p.jabatan.toLowerCase().includes(q) ||
        p.bidangStudi.toLowerCase().includes(q) ||
        labelStatus[p.statusAkun].toLowerCase().includes(q)
      );
    });
  }, [kataKunci, daftarPengguna]);

  const semuaTerlihatTerpilih =
    penggunaTersaring.length > 0 &&
    penggunaTersaring.every((p) => idTerpilih.has(p.id));

  const togglePilihSemuaTerlihat = () => {
    setIdTerpilih((sebelumnya) => {
      const berikutnya = new Set(sebelumnya);
      if (semuaTerlihatTerpilih) {
        penggunaTersaring.forEach((p) => berikutnya.delete(p.id));
      } else {
        penggunaTersaring.forEach((p) => berikutnya.add(p.id));
      }
      return berikutnya;
    });
  };

  const togglePilihBaris = (id: string) => {
    setIdTerpilih((sebelumnya) => {
      const berikutnya = new Set(sebelumnya);
      if (berikutnya.has(id)) berikutnya.delete(id);
      else berikutnya.add(id);
      return berikutnya;
    });
  };

  const jumlahTerpilih = idTerpilih.size;

  return (
    <div className="relative overflow-x-auto bg-neutral-primary-soft shadow-xs rounded-base border border-default">
      <div className="flex items-center justify-between flex-column md:flex-row flex-wrap space-y-4 md:space-y-0 p-4">
        <div className="flex items-center gap-3">
          <div className="relative">
            <button
              className="inline-flex items-center justify-center text-body bg-neutral-secondary-medium box-border border border-default-medium hover:bg-neutral-tertiary-medium hover:text-heading focus:ring-4 focus:ring-neutral-tertiary shadow-xs font-medium leading-5 rounded-base text-sm px-3 py-2 focus:outline-none"
              type="button"
              aria-haspopup="menu"
              aria-expanded={dropdownAksiTerbuka}
              onClick={() => setDropdownAksiTerbuka((v) => !v)}
            >
              Aksi
              <svg
                className="w-4 h-4 ms-1.5 -me-0.5"
                aria-hidden="true"
                xmlns="http://www.w3.org/2000/svg"
                width="24"
                height="24"
                fill="none"
                viewBox="0 0 24 24"
              >
                <path
                  stroke="currentColor"
                  strokeLinecap="round"
                  strokeLinejoin="round"
                  strokeWidth="2"
                  d="m19 9-7 7-7-7"
                />
              </svg>
            </button>

            {dropdownAksiTerbuka && (
              <div
                role="menu"
                className="absolute mt-2 z-10 bg-neutral-primary-medium border border-default-medium rounded-base shadow-lg w-36"
                onMouseLeave={() => setDropdownAksiTerbuka(false)}
              >
                <ul className="p-2 text-sm text-body font-medium">
                  <li>
                    <button
                      type="button"
                      className="inline-flex items-center w-full p-2 hover:bg-neutral-tertiary-medium hover:text-heading rounded text-left"
                      onClick={() => setDropdownAksiTerbuka(false)}
                      disabled={jumlahTerpilih === 0}
                      title={
                        jumlahTerpilih === 0 ? "Pilih minimal 1 pengguna" : ""
                      }
                    >
                      Arsipkan
                    </button>
                  </li>
                  <li>
                    <button
                      type="button"
                      className="inline-flex items-center w-full p-2 text-fg-danger hover:bg-neutral-tertiary-medium rounded text-left"
                      onClick={() => setDropdownAksiTerbuka(false)}
                      disabled={jumlahTerpilih === 0}
                      title={
                        jumlahTerpilih === 0 ? "Pilih minimal 1 pengguna" : ""
                      }
                    >
                      Hapus
                    </button>
                  </li>
                </ul>
              </div>
            )}
          </div>

          <label className="inline-flex items-center gap-2 text-sm text-body select-none">
            <input
              type="checkbox"
              checked={samarkanDataSensitif}
              onChange={(e) => setSamarkanDataSensitif(e.target.checked)}
              className="w-4 h-4 border border-default-medium rounded-xs bg-neutral-secondary-medium focus:ring-2 focus:ring-brand-soft"
            />
            Samarkan data sensitif
          </label>
        </div>

        <label htmlFor="pencarian" className="sr-only">
          Pencarian
        </label>
        <div className="relative">
          <div className="absolute inset-y-0 start-0 flex items-center ps-3 pointer-events-none">
            <svg
              className="w-4 h-4 text-body"
              aria-hidden="true"
              xmlns="http://www.w3.org/2000/svg"
              width="24"
              height="24"
              fill="none"
              viewBox="0 0 24 24"
            >
              <path
                stroke="currentColor"
                strokeLinecap="round"
                strokeWidth="2"
                d="m21 21-3.5-3.5M17 10a7 7 0 1 1-14 0 7 7 0 0 1 14 0Z"
              />
            </svg>
          </div>
          <input
            type="text"
            id="pencarian"
            value={kataKunci}
            onChange={(e) => setKataKunci(e.target.value)}
            className="block w-full max-w-96 ps-9 pe-3 py-2 bg-neutral-secondary-medium border border-default-medium text-heading text-sm rounded-base focus:ring-brand focus:border-brand shadow-xs placeholder:text-body"
            placeholder="Cari pengguna..."
          />
        </div>
      </div>

      <table className="w-full text-sm text-left rtl:text-right text-body">
        <thead className="text-sm text-body bg-neutral-secondary-medium border-b border-t border-default-medium">
          <tr>
            <th scope="col" className="p-4">
              <div className="flex items-center">
                <input
                  id="cek-semua"
                  type="checkbox"
                  checked={semuaTerlihatTerpilih}
                  onChange={togglePilihSemuaTerlihat}
                  className="w-4 h-4 border border-default-medium rounded-xs bg-neutral-secondary-medium focus:ring-2 focus:ring-brand-soft"
                />
                <label htmlFor="cek-semua" className="sr-only">
                  Pilih semua
                </label>
              </div>
            </th>

            <th scope="col" className="px-6 py-3 font-medium">
              Nomor
            </th>
            <th scope="col" className="px-6 py-3 font-medium">
              Nama Lengkap
            </th>
            <th scope="col" className="px-6 py-3 font-medium">
              Nama Pengguna
            </th>
            <th scope="col" className="px-6 py-3 font-medium">
              Jenis Kelamin
            </th>
            <th scope="col" className="px-6 py-3 font-medium">
              Nomor HP
            </th>
            <th scope="col" className="px-6 py-3 font-medium">
              Status Akun
            </th>
            <th scope="col" className="px-6 py-3 font-medium">
              NIP
            </th>
            <th scope="col" className="px-6 py-3 font-medium">
              Jabatan
            </th>
            <th scope="col" className="px-6 py-3 font-medium">
              Bidang Studi
            </th>
            <th scope="col" className="px-6 py-3 font-medium">
              Aksi
            </th>
          </tr>
        </thead>

        <tbody>
          {penggunaTersaring.map((p, indeks) => (
            <tr
              key={p.id}
              className="bg-neutral-primary-soft border-b border-default hover:bg-neutral-secondary-medium"
            >
              <td className="w-4 p-4">
                <div className="flex items-center">
                  <input
                    id={`cek-${p.id}`}
                    type="checkbox"
                    checked={idTerpilih.has(p.id)}
                    onChange={() => togglePilihBaris(p.id)}
                    className="w-4 h-4 border border-default-medium rounded-xs bg-neutral-secondary-medium focus:ring-2 focus:ring-brand-soft"
                  />
                  <label htmlFor={`cek-${p.id}`} className="sr-only">
                    Pilih baris
                  </label>
                </div>
              </td>

              <td className="px-6 py-4">{indeks + 1}</td>

              <td className="px-6 py-4">
                <div className="flex items-center text-heading whitespace-nowrap">
                  <img
                    className="w-10 h-10 rounded-full"
                    src={p.urlGambarProfil}
                    alt={`${p.namaLengkap} avatar`}
                  />
                  <div className="ps-3">
                    <div className="text-base font-semibold">
                      {p.namaLengkap}
                    </div>
                    <div className="font-normal text-body">{p.email}</div>
                  </div>
                </div>
              </td>

              <td className="px-6 py-4">{p.username}</td>

              <td className="px-6 py-4">{p.jenisKelamin}</td>

              <td className="px-6 py-4">
                {samarkanDataSensitif ? samarkanNomorHp(p.nomorHp) : p.nomorHp}
              </td>

              <td className="px-6 py-4">
                <div className="flex items-center">
                  <div
                    className={`h-2.5 w-2.5 rounded-full ${
                      kelasTitikStatus[p.statusAkun]
                    } me-2`}
                  />
                  {labelStatus[p.statusAkun]}
                </div>
              </td>

              <td className="px-6 py-4">
                {samarkanDataSensitif ? samarkanNip(p.nip) : p.nip}
              </td>

              <td className="px-6 py-4">{p.jabatan}</td>

              <td className="px-6 py-4">{p.bidangStudi}</td>

              {/* Tambahkan Link/navigate yang membawa Id sebagai route param */}
              <td className="px-6 py-4">
                <button
                  type="button"
                  className="font-medium text-fg-brand hover:underline"
                  onClick={() => console.log("Edit", p.id)}
                >
                  Ubah
                </button>
              </td>
            </tr>
          ))}

          {penggunaTersaring.length === 0 && (
            <tr className="bg-neutral-primary-soft">
              <td className="px-6 py-8 text-center text-body" colSpan={10}>
                Tidak ada data.
              </td>
            </tr>
          )}
        </tbody>
      </table>
      <div className="px-4 py-3 text-xs text-slate-500 border-t border-default">
        Geser tabel ke kanan/kiri untuk melihat kolom lainnya.
      </div>
    </div>
  );
};
