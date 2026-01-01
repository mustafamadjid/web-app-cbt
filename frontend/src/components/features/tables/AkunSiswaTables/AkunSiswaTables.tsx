import React, { useMemo, useState } from "react";

type StatusAkun = "aktif" | "nonaktif" | "dibekukan";
type JenisKelamin = "LAKI_LAKI" | "PEREMPUAN";

type BarisSiswa = {
  id: string;
  namaLengkap: string;
  username: string;
  email?: string;
  nomorHp?: string;
  jenisKelamin: JenisKelamin;
  statusAkun: StatusAkun;
  noAbsen: number;
  angkatan: number;
  tempatLahir: string;
  tanggalLahir: string;
  kelas: string;
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

function samarkanEmail(email: string) {
  const e = email.trim();
  const [user, domain] = e.split("@");
  if (!domain) return e;
  if (user.length <= 2) return `**@${domain}`;
  return `${user.slice(0, 2)}****@${domain}`;
}

function formatTanggalIndo(yyyyMmDd: string) {
  if (!yyyyMmDd) return "-";
  const [y, m, d] = yyyyMmDd.split("-");
  if (!y || !m || !d) return yyyyMmDd;
  return `${d}/${m}/${y}`;
}

export const AkunSiswaTables: React.FC = () => {
  const [dropdownAksiTerbuka, setDropdownAksiTerbuka] = useState(false);
  const [kataKunci, setKataKunci] = useState("");
  const [idTerpilih, setIdTerpilih] = useState<Set<string>>(new Set());
  const [samarkanDataSensitif, setSamarkanDataSensitif] = useState(true);

  const [daftarSiswa] = useState<BarisSiswa[]>([
    {
      id: "s1",
      namaLengkap: "Siti Aminah",
      username: "siti.aminah",
      email: "siti.aminah@gmail.com",
      nomorHp: "081234567890",
      jenisKelamin: "PEREMPUAN",
      statusAkun: "aktif",
      noAbsen: 12,
      angkatan: 2025,
      tempatLahir: "Bandung",
      tanggalLahir: "2008-01-31",
      kelas: "X IPA 1",
      urlGambarProfil: "/docs/images/people/profile-picture-1.jpg",
    },
    {
      id: "s2",
      namaLengkap: "Raka Pratama",
      username: "raka.pratama",
      email: "",
      nomorHp: "",
      jenisKelamin: "LAKI_LAKI",
      statusAkun: "nonaktif",
      noAbsen: 7,
      angkatan: 2025,
      tempatLahir: "Jakarta",
      tanggalLahir: "2008-08-12",
      kelas: "X IPS 1",
      urlGambarProfil: "/docs/images/people/profile-picture-3.jpg",
    },
  ]);

  const siswaTersaring = useMemo(() => {
    const q = kataKunci.trim().toLowerCase();
    if (!q) return daftarSiswa;

    return daftarSiswa.filter((s) => {
      const email = (s.email ?? "").toLowerCase();
      const hp = (s.nomorHp ?? "").toLowerCase();
      return (
        s.namaLengkap.toLowerCase().includes(q) ||
        s.username.toLowerCase().includes(q) ||
        email.includes(q) ||
        hp.includes(q) ||
        String(s.noAbsen).includes(q) ||
        String(s.angkatan).includes(q) ||
        s.jenisKelamin.toLowerCase().includes(q) ||
        s.tempatLahir.toLowerCase().includes(q) ||
        s.tanggalLahir.toLowerCase().includes(q) ||
        s.kelas.toLowerCase().includes(q) ||
        labelStatus[s.statusAkun].toLowerCase().includes(q)
      );
    });
  }, [kataKunci, daftarSiswa]);

  const semuaTerlihatTerpilih =
    siswaTersaring.length > 0 &&
    siswaTersaring.every((s) => idTerpilih.has(s.id));

  const togglePilihSemuaTerlihat = () => {
    setIdTerpilih((prev) => {
      const next = new Set(prev);
      if (semuaTerlihatTerpilih)
        siswaTersaring.forEach((s) => next.delete(s.id));
      else siswaTersaring.forEach((s) => next.add(s.id));
      return next;
    });
  };

  const togglePilihBaris = (id: string) => {
    setIdTerpilih((prev) => {
      const next = new Set(prev);
      if (next.has(id)) next.delete(id);
      else next.add(id);
      return next;
    });
  };

  const jumlahTerpilih = idTerpilih.size;

  return (
    <div className="bg-neutral-primary-soft shadow-xs rounded-base border border-default w-full min-w-0">
      {/* Topbar */}
      <div className="flex items-center justify-between flex-wrap gap-3 p-4">
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
                className="absolute mt-2 z-20 bg-neutral-primary-medium border border-default-medium rounded-base shadow-lg w-36"
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
                        jumlahTerpilih === 0 ? "Pilih minimal 1 siswa" : ""
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
                        jumlahTerpilih === 0 ? "Pilih minimal 1 siswa" : ""
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

        <div className="relative w-full sm:w-80 md:w-96">
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
            className="block w-full ps-9 pe-3 py-2 bg-neutral-secondary-medium border border-default-medium text-heading text-sm rounded-base focus:ring-green-800 focus:border-green-800 shadow-xs placeholder:text-body"
            placeholder="Cari siswa..."
          />
        </div>
      </div>

      {/* SCROLL WRAPPER: ini kunci responsive */}
      <div className="w-full overflow-x-auto">
        <table className="w-max min-w-full text-sm text-left rtl:text-right text-body">
          <thead className="text-sm text-body bg-neutral-secondary-medium border-b border-t border-default-medium">
            <tr className="whitespace-nowrap">
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
                Email
              </th>

              <th scope="col" className="px-6 py-3 font-medium">
                No Absen
              </th>
              <th scope="col" className="px-6 py-3 font-medium">
                Angkatan
              </th>
              <th scope="col" className="px-6 py-3 font-medium">
                Kelas
              </th>
              <th scope="col" className="px-6 py-3 font-medium">
                TTL
              </th>
              <th scope="col" className="px-6 py-3 font-medium">
                Status Akun
              </th>
              <th scope="col" className="px-6 py-3 font-medium">
                Aksi
              </th>
            </tr>
          </thead>

          <tbody>
            {siswaTersaring.map((s, indeks) => {
              const hpRaw = s.nomorHp ?? "";
              const emailRaw = s.email ?? "";

              const hpTampil = !hpRaw
                ? "-"
                : samarkanDataSensitif
                ? samarkanNomorHp(hpRaw)
                : hpRaw;

              const emailTampil = !emailRaw
                ? "-"
                : samarkanDataSensitif
                ? samarkanEmail(emailRaw)
                : emailRaw;

              return (
                <tr
                  key={s.id}
                  className="bg-neutral-primary-soft border-b border-default hover:bg-neutral-secondary-medium whitespace-nowrap"
                >
                  <td className="w-4 p-4">
                    <div className="flex items-center">
                      <input
                        id={`cek-${s.id}`}
                        type="checkbox"
                        checked={idTerpilih.has(s.id)}
                        onChange={() => togglePilihBaris(s.id)}
                        className="w-4 h-4 border border-default-medium rounded-xs bg-neutral-secondary-medium focus:ring-2 focus:ring-brand-soft"
                      />
                      <label htmlFor={`cek-${s.id}`} className="sr-only">
                        Pilih baris
                      </label>
                    </div>
                  </td>

                  <td className="px-6 py-4">{indeks + 1}</td>

                  <td className="px-6 py-4">
                    <div className="flex items-center text-heading">
                      <img
                        className="w-10 h-10 rounded-full"
                        src={s.urlGambarProfil}
                        alt={`${s.namaLengkap} avatar`}
                      />
                      <div className="ps-3">
                        <div className="text-base font-semibold">
                          {s.namaLengkap}
                        </div>
                        <div className="font-normal text-body">
                          {emailTampil}
                        </div>
                      </div>
                    </div>
                  </td>

                  <td className="px-6 py-4">{s.username}</td>
                  <td className="px-6 py-4">{s.jenisKelamin}</td>
                  <td className="px-6 py-4">{hpTampil}</td>
                  <td className="px-6 py-4">{emailTampil}</td>
                  <td className="px-6 py-4">{s.noAbsen}</td>
                  <td className="px-6 py-4">{s.angkatan}</td>
                  <td className="px-6 py-4">{s.kelas}</td>
                  <td className="px-6 py-4">
                    {s.tempatLahir}, {formatTanggalIndo(s.tanggalLahir)}
                  </td>

                  <td className="px-6 py-4">
                    <div className="flex items-center">
                      <div
                        className={`h-2.5 w-2.5 rounded-full ${
                          kelasTitikStatus[s.statusAkun]
                        } me-2`}
                      />
                      {labelStatus[s.statusAkun]}
                    </div>
                  </td>

                  <td className="px-6 py-4">
                    <button
                      type="button"
                      className="font-medium text-fg-brand hover:underline"
                      onClick={() => console.log("Edit", s.id)}
                    >
                      Ubah
                    </button>
                  </td>
                </tr>
              );
            })}

            {siswaTersaring.length === 0 && (
              <tr className="bg-neutral-primary-soft">
                <td className="px-6 py-8 text-center text-body" colSpan={13}>
                  Tidak ada data.
                </td>
              </tr>
            )}
          </tbody>
        </table>
      </div>

      <div className="px-4 py-3 text-xs text-slate-500 border-t border-default">
        Geser tabel ke kanan/kiri untuk melihat kolom lainnya.
      </div>
    </div>
  );
};
