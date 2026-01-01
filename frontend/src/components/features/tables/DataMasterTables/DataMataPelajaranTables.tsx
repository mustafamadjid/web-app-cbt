import React, { useMemo, useState } from "react";

import type {
  KelasOption,
  MataPelajaranOption,
  MataPelajaranRow,
} from "../../../../types/DataMaster/MataPelajaran";

export const DataMataPelajaran: React.FC = () => {
  const [dropdownAksiTerbuka, setDropdownAksiTerbuka] = useState(false);
  const [kataKunci, setKataKunci] = useState("");
  const [idTerpilih, setIdTerpilih] = useState<Set<string>>(new Set());
  const [dropdownKelasTerbuka, setDropdownKelasTerbuka] = useState(false);
  const [dropdownMapelTerbuka, setDropdownMapelTerbuka] = useState(false);
  const [kelasTerpilih, setKelasTerpilih] = useState<string>("semua");
  const [mapelTerpilih, setMapelTerpilih] = useState<string>("semua");

  // tetap dipakai untuk menampilkan label kelas (bukan pilihan)
  const [opsiKelas] = useState<KelasOption[]>([
    { id: "kelas-10-ipa-1", label: "X IPA 1" },
    { id: "kelas-10-ipa-2", label: "X IPA 2" },
    { id: "kelas-10-ips-1", label: "X IPS 1" },
  ]);

  const [daftarMapel] = useState<MataPelajaranRow[]>([
    {
      id: "mapel-1",
      kelasId: "kelas-10-ipa-1",
      kodeMapel: "MAT-10-01",
      namaMapel: "Matematika",
      deskripsiMapel: "Aljabar dasar, geometri, dan statistika.",
    },
    {
      id: "mapel-2",
      kelasId: "kelas-10-ips-1",
      kodeMapel: "EKO-10-01",
      namaMapel: "Ekonomi",
      deskripsiMapel: "Dasar-dasar ekonomi mikro dan makro.",
    },
  ]);

  const [opsiMapel] = useState<MataPelajaranOption[]>([
    { id: "mapel-1", label: "Matematika" },
    { id: "mapel-2", label: "Ekonomi" },
  ]);

  const kelasById = useMemo(() => {
    return opsiKelas.reduce<Record<string, KelasOption>>((acc, opsi) => {
      acc[opsi.id] = opsi;
      return acc;
    }, {});
  }, [opsiKelas]);

  const mapelTersaring = useMemo(() => {
    const q = kataKunci.trim().toLowerCase();
    return daftarMapel.filter((mapel) => {
      const labelKelas = kelasById[mapel.kelasId]?.label ?? "";
      const cocokKata =
        !q ||
        mapel.kodeMapel.toLowerCase().includes(q) ||
        mapel.namaMapel.toLowerCase().includes(q) ||
        mapel.deskripsiMapel.toLowerCase().includes(q) ||
        labelKelas.toLowerCase().includes(q);
      const cocokKelas = kelasTerpilih === "semua" || mapel.kelasId === kelasTerpilih;
      const cocokMapel = mapelTerpilih === "semua" || mapel.id === mapelTerpilih;
      return cocokKata && cocokKelas && cocokMapel;
    });
  }, [kataKunci, daftarMapel, kelasById, kelasTerpilih, mapelTerpilih]);

  const semuaTerlihatTerpilih =
    mapelTersaring.length > 0 &&
    mapelTersaring.every((mapel) => idTerpilih.has(mapel.id));

  const togglePilihSemuaTerlihat = () => {
    setIdTerpilih((prev) => {
      const next = new Set(prev);
      if (semuaTerlihatTerpilih) {
        mapelTersaring.forEach((mapel) => next.delete(mapel.id));
      } else {
        mapelTersaring.forEach((mapel) => next.add(mapel.id));
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
      <div className="flex items-center justify-between flex-wrap gap-3 p-4">
        <div className="flex items-center flex-wrap gap-3">
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
                        jumlahTerpilih === 0
                          ? "Pilih minimal 1 mata pelajaran"
                          : ""
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
                        jumlahTerpilih === 0
                          ? "Pilih minimal 1 mata pelajaran"
                          : ""
                      }
                    >
                      Hapus
                    </button>
                  </li>
                </ul>
              </div>
            )}
          </div>

          <div className="flex items-center flex-wrap gap-2">
            <span className="text-xs font-medium text-body">Urutkan:</span>
            <div className="relative">
              <button
                type="button"
                onClick={() => setDropdownKelasTerbuka((prev) => !prev)}
                className="inline-flex items-center gap-2 text-sm cursor-pointer px-3 py-2 rounded-base border shadow-xs bg-neutral-secondary-medium text-body border-default-medium hover:bg-neutral-tertiary-medium hover:text-heading"
                aria-haspopup="menu"
                aria-expanded={dropdownKelasTerbuka}
              >
                Kelas:{" "}
                <span className="font-medium text-heading cursor-pointer">
                  {kelasTerpilih === "semua"
                    ? "Semua"
                    : kelasById[kelasTerpilih]?.label ?? "-"}
                </span>
                <svg
                  className="w-4 h-4"
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
              {dropdownKelasTerbuka && (
                <div
                  role="menu"
                  className="absolute mt-2 z-20 bg-neutral-primary-medium border border-default-medium rounded-base shadow-lg w-48"
                  onMouseLeave={() => setDropdownKelasTerbuka(false)}
                >
                  <ul className="p-2 text-sm text-body font-medium">
                    <li>
                      <button
                        type="button"
                        className="inline-flex items-center w-full p-2 hover:bg-neutral-tertiary-medium hover:text-heading rounded text-left"
                        onClick={() => {
                          setKelasTerpilih("semua");
                          setDropdownKelasTerbuka(false);
                        }}
                      >
                        Semua Kelas
                      </button>
                    </li>
                    {opsiKelas.map((opsi) => (
                      <li key={opsi.id}>
                        <button
                          type="button"
                          className="inline-flex items-center w-full p-2 hover:bg-neutral-tertiary-medium hover:text-heading rounded text-left"
                          onClick={() => {
                            setKelasTerpilih(opsi.id);
                            setDropdownKelasTerbuka(false);
                          }}
                        >
                          {opsi.label}
                        </button>
                      </li>
                    ))}
                  </ul>
                </div>
              )}
            </div>

            <div className="relative">
              <button
                type="button"
                onClick={() => setDropdownMapelTerbuka((prev) => !prev)}
                className="inline-flex cursor-pointer items-center gap-2 text-sm px-3 py-2 rounded-base border shadow-xs bg-neutral-secondary-medium text-body border-default-medium hover:bg-neutral-tertiary-medium hover:text-heading"
                aria-haspopup="menu"
                aria-expanded={dropdownMapelTerbuka}
              >
                Mapel:{" "}
                <span className="font-medium text-heading">
                  {mapelTerpilih === "semua"
                    ? "Semua"
                    : opsiMapel.find((opsi) => opsi.id === mapelTerpilih)
                        ?.label ?? "-"}
                </span>
                <svg
                  className="w-4 h-4"
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
              {dropdownMapelTerbuka && (
                <div
                  role="menu"
                  className="absolute mt-2 z-20 bg-neutral-primary-medium border border-default-medium rounded-base shadow-lg w-52"
                  onMouseLeave={() => setDropdownMapelTerbuka(false)}
                >
                  <ul className="p-2 text-sm text-body font-medium">
                    <li>
                      <button
                        type="button"
                        className="inline-flex items-center w-full p-2 hover:bg-neutral-tertiary-medium hover:text-heading rounded text-left"
                        onClick={() => {
                          setMapelTerpilih("semua");
                          setDropdownMapelTerbuka(false);
                        }}
                      >
                        Semua Mata Pelajaran
                      </button>
                    </li>
                    {opsiMapel.map((opsi) => (
                      <li key={opsi.id}>
                        <button
                          type="button"
                          className="inline-flex items-center w-full p-2 hover:bg-neutral-tertiary-medium hover:text-heading rounded text-left"
                          onClick={() => {
                            setMapelTerpilih(opsi.id);
                            setDropdownMapelTerbuka(false);
                          }}
                        >
                          {opsi.label}
                        </button>
                      </li>
                    ))}
                  </ul>
                </div>
              )}
            </div>
          </div>
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
            id="pencarian-mapel"
            value={kataKunci}
            onChange={(e) => setKataKunci(e.target.value)}
            className="block w-full ps-9 pe-3 py-2 bg-neutral-secondary-medium border border-default-medium text-heading text-sm rounded-base focus:ring-brand focus:border-brand shadow-xs placeholder:text-body"
            placeholder="Cari mata pelajaran..."
          />
        </div>
      </div>

      <div className="w-full overflow-x-auto">
        <table className="w-max min-w-full text-sm text-left rtl:text-right text-body">
          <thead className="text-sm text-body bg-neutral-secondary-medium border-b border-t border-default-medium">
            <tr className="whitespace-nowrap">
              <th scope="col" className="p-4">
                <div className="flex items-center">
                  <input
                    id="cek-semua-mapel"
                    type="checkbox"
                    checked={semuaTerlihatTerpilih}
                    onChange={togglePilihSemuaTerlihat}
                    className="w-4 h-4 border border-default-medium rounded-xs bg-neutral-secondary-medium focus:ring-2 focus:ring-brand-soft"
                  />
                  <label htmlFor="cek-semua-mapel" className="sr-only">
                    Pilih semua
                  </label>
                </div>
              </th>

              <th scope="col" className="px-6 py-3 font-medium">
                Nomor
              </th>
              <th scope="col" className="px-6 py-3 font-medium">
                Kelas
              </th>
              <th scope="col" className="px-6 py-3 font-medium">
                Kode Mapel
              </th>
              <th scope="col" className="px-6 py-3 font-medium">
                Nama Mapel
              </th>
              <th scope="col" className="px-6 py-3 font-medium">
                Deskripsi Mapel
              </th>
              <th scope="col" className="px-6 py-3 font-medium">
                Aksi
              </th>
            </tr>
          </thead>

          <tbody>
            {mapelTersaring.map((mapel, index) => {
              const labelKelas = kelasById[mapel.kelasId]?.label ?? "-";

              return (
                <tr
                  key={mapel.id}
                  className="bg-neutral-primary-soft border-b border-default hover:bg-neutral-secondary-medium whitespace-nowrap"
                >
                  <td className="w-4 p-4">
                    <div className="flex items-center">
                      <input
                        id={`cek-${mapel.id}`}
                        type="checkbox"
                        checked={idTerpilih.has(mapel.id)}
                        onChange={() => togglePilihBaris(mapel.id)}
                        className="w-4 h-4 border border-default-medium rounded-xs bg-neutral-secondary-medium focus:ring-2 focus:ring-brand-soft"
                      />
                      <label htmlFor={`cek-${mapel.id}`} className="sr-only">
                        Pilih baris
                      </label>
                    </div>
                  </td>

                  <td className="px-6 py-4">{index + 1}</td>

                  {/* KELAS: teks biasa (tanpa dropdown) */}
                  <td className="px-6 py-4 text-heading">{labelKelas}</td>

                  <td className="px-6 py-4 text-heading">
                    {mapel.kodeMapel}
                  </td>
                  <td className="px-6 py-4 text-heading font-medium">
                    {mapel.namaMapel}
                  </td>
                  <td className="px-6 py-4 text-body">
                    {mapel.deskripsiMapel}
                  </td>
                  <td className="px-6 py-4">
                    <button
                      type="button"
                      className="font-medium text-fg-brand hover:underline"
                      onClick={() => console.log("Ubah", mapel.id)}
                    >
                      Ubah
                    </button>
                  </td>
                </tr>
              );
            })}

            {mapelTersaring.length === 0 && (
              <tr className="bg-neutral-primary-soft">
                <td className="px-6 py-8 text-center text-body" colSpan={7}>
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
