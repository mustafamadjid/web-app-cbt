import React, { useMemo, useState } from "react";

type RuangUjianRow = {
  id: string;
  namaRuangan: string;
};

export const RuangUjianTables: React.FC = () => {
  const [dropdownAksiTerbuka, setDropdownAksiTerbuka] = useState(false);
  const [kataKunci, setKataKunci] = useState("");
  const [idTerpilih, setIdTerpilih] = useState<Set<string>>(new Set());

  // Dummy data (nanti ganti dari API)
  const [daftarRuangUjian] = useState<RuangUjianRow[]>([
    { id: "ruang-1", namaRuangan: "Ruang Ujian 01" },
    { id: "ruang-2", namaRuangan: "Ruang Ujian 02" },
    { id: "ruang-3", namaRuangan: "Lab Komputer" },
    { id: "ruang-4", namaRuangan: "Aula" },
  ]);

  const ruangTersaring = useMemo(() => {
    const q = kataKunci.trim().toLowerCase();
    if (!q) return daftarRuangUjian;

    return daftarRuangUjian.filter((r) =>
      r.namaRuangan.toLowerCase().includes(q)
    );
  }, [kataKunci, daftarRuangUjian]);

  const semuaTerlihatTerpilih =
    ruangTersaring.length > 0 &&
    ruangTersaring.every((r) => idTerpilih.has(r.id));

  const togglePilihSemuaTerlihat = () => {
    setIdTerpilih((prev) => {
      const next = new Set(prev);
      if (semuaTerlihatTerpilih) {
        ruangTersaring.forEach((r) => next.delete(r.id));
      } else {
        ruangTersaring.forEach((r) => next.add(r.id));
      }
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
      {/* Header actions */}
      <div className="flex items-center justify-between flex-wrap gap-3 p-4">
        <div className="flex items-center flex-wrap gap-3">
          {/* Dropdown Aksi */}
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
                      className="inline-flex items-center w-full p-2 hover:bg-neutral-tertiary-medium hover:text-heading rounded text-left disabled:opacity-60 disabled:cursor-not-allowed"
                      onClick={() => setDropdownAksiTerbuka(false)}
                      disabled={jumlahTerpilih === 0}
                      title={
                        jumlahTerpilih === 0 ? "Pilih minimal 1 ruangan" : ""
                      }
                    >
                      Arsipkan
                    </button>
                  </li>
                  <li>
                    <button
                      type="button"
                      className="inline-flex items-center w-full p-2 text-fg-danger hover:bg-neutral-tertiary-medium rounded text-left disabled:opacity-60 disabled:cursor-not-allowed"
                      onClick={() => setDropdownAksiTerbuka(false)}
                      disabled={jumlahTerpilih === 0}
                      title={
                        jumlahTerpilih === 0 ? "Pilih minimal 1 ruangan" : ""
                      }
                    >
                      Hapus
                    </button>
                  </li>
                </ul>
              </div>
            )}
          </div>

          {jumlahTerpilih > 0 && (
            <div className="text-xs text-body">
              Dipilih:{" "}
              <span className="font-medium text-heading">{jumlahTerpilih}</span>
            </div>
          )}
        </div>

        {/* Search */}
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
            id="pencarian-ruang-ujian"
            value={kataKunci}
            onChange={(e) => setKataKunci(e.target.value)}
            className="block w-full ps-9 pe-3 py-2 bg-neutral-secondary-medium border border-default-medium text-heading text-sm rounded-base focus:ring-brand focus:border-brand shadow-xs placeholder:text-body"
            placeholder="Cari ruangan..."
          />
        </div>
      </div>

      {/* Table */}
      <div className="w-full overflow-x-auto">
        <table className="w-max min-w-full text-sm text-left rtl:text-right text-body">
          <thead className="text-sm text-body bg-neutral-secondary-medium border-b border-t border-default-medium">
            <tr className="whitespace-nowrap">
              <th scope="col" className="p-4">
                <div className="flex items-center">
                  <input
                    id="cek-semua-ruang"
                    type="checkbox"
                    checked={semuaTerlihatTerpilih}
                    onChange={togglePilihSemuaTerlihat}
                    className="w-4 h-4 border border-default-medium rounded-xs bg-neutral-secondary-medium focus:ring-2 focus:ring-brand-soft"
                  />
                  <label htmlFor="cek-semua-ruang" className="sr-only">
                    Pilih semua
                  </label>
                </div>
              </th>

              <th scope="col" className="px-6 py-3 font-medium">
                Nomor
              </th>
              <th scope="col" className="px-6 py-3 font-medium">
                Nama Ruangan
              </th>
              <th scope="col" className="px-6 py-3 font-medium">
                Aksi
              </th>
            </tr>
          </thead>

          <tbody>
            {ruangTersaring.map((ruang, index) => (
              <tr
                key={ruang.id}
                className="bg-neutral-primary-soft border-b border-default hover:bg-neutral-secondary-medium whitespace-nowrap"
              >
                <td className="w-4 p-4">
                  <div className="flex items-center">
                    <input
                      id={`cek-${ruang.id}`}
                      type="checkbox"
                      checked={idTerpilih.has(ruang.id)}
                      onChange={() => togglePilihBaris(ruang.id)}
                      className="w-4 h-4 border border-default-medium rounded-xs bg-neutral-secondary-medium focus:ring-2 focus:ring-brand-soft"
                    />
                    <label htmlFor={`cek-${ruang.id}`} className="sr-only">
                      Pilih baris
                    </label>
                  </div>
                </td>

                <td className="px-6 py-4">{index + 1}</td>

                <td className="px-6 py-4 text-heading font-medium">
                  {ruang.namaRuangan}
                </td>

                <td className="px-6 py-4">
                  <button
                    type="button"
                    className="font-medium text-fg-brand hover:underline cursor-pointer"
                    onClick={() => console.log("Ubah", ruang.id)}
                  >
                    Ubah
                  </button>
                </td>
              </tr>
            ))}

            {ruangTersaring.length === 0 && (
              <tr className="bg-neutral-primary-soft">
                <td className="px-6 py-8 text-center text-body" colSpan={4}>
                  Tidak ada data.
                </td>
              </tr>
            )}
          </tbody>
        </table>
      </div>
    </div>
  );
};
