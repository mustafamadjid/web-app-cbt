import React, { useMemo, useState } from "react";

export type KelasRow = {
  id: string;
  tingkat_kelas: number;
  nama_kelas: string;
};

type SortKey = "tingkat_kelas" | "nama_kelas";
type SortDir = "asc" | "desc";

export const DataKelasTables: React.FC = () => {
  const [kataKunci, setKataKunci] = useState("");

  // dropdown
  const [dropdownAksiTerbuka, setDropdownAksiTerbuka] = useState(false);
  const [dropdownSortTerbuka, setDropdownSortTerbuka] = useState(false);

  // selection
  const [idTerpilih, setIdTerpilih] = useState<Set<string>>(new Set());

  // sorting
  const [sortKey, setSortKey] = useState<SortKey>("tingkat_kelas");
  const [sortDir, setSortDir] = useState<SortDir>("asc");

  // dummy data
  const [daftarKelas] = useState<KelasRow[]>([
    { id: "kelas-1", tingkat_kelas: 12, nama_kelas: "XII IPA 1" },
    { id: "kelas-2", tingkat_kelas: 10, nama_kelas: "X IPA 2" },
    { id: "kelas-3", tingkat_kelas: 11, nama_kelas: "XI IPS 1" },
    { id: "kelas-4", tingkat_kelas: 10, nama_kelas: "X IPA 1" },
  ]);

  const kelasTersaringDanTerurut = useMemo(() => {
    const q = kataKunci.trim().toLowerCase();

    const filtered = daftarKelas.filter((k) => {
      if (!q) return true;
      return (
        String(k.tingkat_kelas).includes(q) ||
        k.nama_kelas.toLowerCase().includes(q)
      );
    });

    const sorted = [...filtered].sort((a, b) => {
      let cmp = 0;
      if (sortKey === "tingkat_kelas") {
        cmp = a.tingkat_kelas - b.tingkat_kelas;
      } else {
        cmp = a.nama_kelas.localeCompare(b.nama_kelas, "id", {
          sensitivity: "base",
          numeric: true,
        });
      }
      return sortDir === "asc" ? cmp : -cmp;
    });

    return sorted;
  }, [daftarKelas, kataKunci, sortKey, sortDir]);

  const semuaTerlihatTerpilih =
    kelasTersaringDanTerurut.length > 0 &&
    kelasTersaringDanTerurut.every((k) => idTerpilih.has(k.id));

  const togglePilihSemuaTerlihat = () => {
    setIdTerpilih((prev) => {
      const next = new Set(prev);
      if (semuaTerlihatTerpilih) {
        kelasTersaringDanTerurut.forEach((k) => next.delete(k.id));
      } else {
        kelasTersaringDanTerurut.forEach((k) => next.add(k.id));
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

  const toggleSort = (key: SortKey) => {
    setSortKey((prevKey) => {
      if (prevKey === key) {
        setSortDir((prevDir) => (prevDir === "asc" ? "desc" : "asc"));
        return prevKey;
      }
      setSortDir("asc");
      return key;
    });
  };

  const labelSort = useMemo(() => {
    const keyLabel =
      sortKey === "tingkat_kelas" ? "Tingkat Kelas" : "Nama Kelas";
    const dirLabel = sortDir === "asc" ? "ASC" : "DESC";
    return `${keyLabel} • ${dirLabel}`;
  }, [sortKey, sortDir]);

  return (
    <div className="bg-neutral-primary-soft shadow-xs rounded-base border border-default w-full min-w-0">
      <div className="flex items-center justify-between flex-wrap gap-3 p-4">
        <div className="flex items-center flex-wrap gap-3">
          {/* Aksi */}
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
                      onClick={() => {
                        // contoh action
                        console.log("Arsipkan:", Array.from(idTerpilih));
                        setDropdownAksiTerbuka(false);
                      }}
                      disabled={jumlahTerpilih === 0}
                      title={
                        jumlahTerpilih === 0 ? "Pilih minimal 1 kelas" : ""
                      }
                    >
                      Arsipkan
                    </button>
                  </li>
                  <li>
                    <button
                      type="button"
                      className="inline-flex items-center w-full p-2 text-fg-danger hover:bg-neutral-tertiary-medium rounded text-left"
                      onClick={() => {
                        console.log("Hapus:", Array.from(idTerpilih));
                        setDropdownAksiTerbuka(false);
                      }}
                      disabled={jumlahTerpilih === 0}
                      title={
                        jumlahTerpilih === 0 ? "Pilih minimal 1 kelas" : ""
                      }
                    >
                      Hapus
                    </button>
                  </li>
                </ul>
              </div>
            )}
          </div>

          {/* Sorting */}
          <div className="relative">
            <button
              type="button"
              onClick={() => setDropdownSortTerbuka((v) => !v)}
              className="inline-flex items-center gap-2 text-sm px-3 py-2 rounded-base border shadow-xs bg-neutral-secondary-medium text-body border-default-medium hover:bg-neutral-tertiary-medium hover:text-heading"
              aria-haspopup="menu"
              aria-expanded={dropdownSortTerbuka}
            >
              Urutkan:{" "}
              <span className="font-medium text-heading">{labelSort}</span>
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

            {dropdownSortTerbuka && (
              <div
                role="menu"
                className="absolute mt-2 z-20 bg-neutral-primary-medium border border-default-medium rounded-base shadow-lg w-64"
                onMouseLeave={() => setDropdownSortTerbuka(false)}
              >
                <ul className="p-2 text-sm text-body font-medium">
                  <li>
                    <button
                      type="button"
                      className="inline-flex items-center justify-between w-full p-2 hover:bg-neutral-tertiary-medium hover:text-heading rounded text-left"
                      onClick={() => toggleSort("tingkat_kelas")}
                    >
                      <span>Tingkat Kelas</span>
                      <span className="text-xs text-body">
                        {sortKey === "tingkat_kelas"
                          ? sortDir === "asc"
                            ? "ASC"
                            : "DESC"
                          : ""}
                      </span>
                    </button>
                  </li>
                  <li>
                    <button
                      type="button"
                      className="inline-flex items-center justify-between w-full p-2 hover:bg-neutral-tertiary-medium hover:text-heading rounded text-left"
                      onClick={() => toggleSort("nama_kelas")}
                    >
                      <span>Nama Kelas</span>
                      <span className="text-xs text-body">
                        {sortKey === "nama_kelas"
                          ? sortDir === "asc"
                            ? "ASC"
                            : "DESC"
                          : ""}
                      </span>
                    </button>
                  </li>

                  <li className="border-t border-default-medium my-2" />

                  <li>
                    <button
                      type="button"
                      className="inline-flex items-center w-full p-2 hover:bg-neutral-tertiary-medium hover:text-heading rounded text-left"
                      onClick={() => setSortDir("asc")}
                    >
                      Ascending
                    </button>
                  </li>
                  <li>
                    <button
                      type="button"
                      className="inline-flex items-center w-full p-2 hover:bg-neutral-tertiary-medium hover:text-heading rounded text-left"
                      onClick={() => setSortDir("desc")}
                    >
                      Descending
                    </button>
                  </li>
                </ul>
              </div>
            )}
          </div>
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
            id="pencarian-kelas"
            value={kataKunci}
            onChange={(e) => setKataKunci(e.target.value)}
            className="block w-full ps-9 pe-3 py-2 bg-neutral-secondary-medium border border-default-medium text-heading text-sm rounded-base focus:ring-brand focus:border-brand shadow-xs placeholder:text-body"
            placeholder="Cari kelas..."
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
                    id="cek-semua-kelas"
                    type="checkbox"
                    checked={semuaTerlihatTerpilih}
                    onChange={togglePilihSemuaTerlihat}
                    className="w-4 h-4 border border-default-medium rounded-xs bg-neutral-secondary-medium focus:ring-2 focus:ring-brand-soft"
                  />
                  <label htmlFor="cek-semua-kelas" className="sr-only">
                    Pilih semua
                  </label>
                </div>
              </th>

              <th className="px-6 py-3 font-medium">Nomor</th>
              <th className="px-6 py-3 font-medium">Tingkat Kelas</th>
              <th className="px-6 py-3 font-medium">Nama Kelas</th>
              <th className="px-6 py-3 font-medium">Aksi</th>
            </tr>
          </thead>

          <tbody>
            {kelasTersaringDanTerurut.map((kelas, index) => (
              <tr
                key={kelas.id}
                className="bg-neutral-primary-soft border-b border-default hover:bg-neutral-secondary-medium whitespace-nowrap"
              >
                <td className="w-4 p-4">
                  <div className="flex items-center">
                    <input
                      id={`cek-${kelas.id}`}
                      type="checkbox"
                      checked={idTerpilih.has(kelas.id)}
                      onChange={() => togglePilihBaris(kelas.id)}
                      className="w-4 h-4 border border-default-medium rounded-xs bg-neutral-secondary-medium focus:ring-2 focus:ring-brand-soft"
                    />
                    <label htmlFor={`cek-${kelas.id}`} className="sr-only">
                      Pilih baris
                    </label>
                  </div>
                </td>

                <td className="px-6 py-4">{index + 1}</td>
                <td className="px-6 py-4 text-heading">
                  {kelas.tingkat_kelas}
                </td>
                <td className="px-6 py-4 text-heading font-medium">
                  {kelas.nama_kelas}
                </td>

                <td className="px-6 py-4">
                  <button
                    type="button"
                    className="font-medium text-fg-brand hover:underline cursor-pointer"
                    onClick={() => console.log("Ubah", kelas.id)}
                  >
                    Ubah
                  </button>
                </td>
              </tr>
            ))}

            {kelasTersaringDanTerurut.length === 0 && (
              <tr className="bg-neutral-primary-soft">
                <td className="px-6 py-8 text-center text-body" colSpan={5}>
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
