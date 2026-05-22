import BankSoalSelect from "@/components/common/Select/BankSoalSelect";

import { selectBaseClass } from "./constants";
import FormSection from "./FormSection";
import type { BuatUjianFormController } from "./useBuatUjianForm";

type KelasBankSoalSectionProps = Pick<
  BuatUjianFormController,
  | "bankSoalOptions"
  | "bankSoalPlaceholder"
  | "errors"
  | "hasError"
  | "loadingBankSoal"
  | "loadingKelas"
  | "loadingMapel"
  | "mapelOptions"
  | "namaKelasOptions"
  | "onBlur"
  | "selectedMapelId"
  | "setField"
  | "setSelectedMapelId"
  | "submitting"
  | "tingkatKelasOptions"
  | "values"
>;

const KelasBankSoalSection = ({
  bankSoalOptions,
  bankSoalPlaceholder,
  errors,
  hasError,
  loadingBankSoal,
  loadingKelas,
  loadingMapel,
  mapelOptions,
  namaKelasOptions,
  onBlur,
  selectedMapelId,
  setField,
  setSelectedMapelId,
  submitting,
  tingkatKelasOptions,
  values,
}: KelasBankSoalSectionProps) => (
  <FormSection title="Kelas & Bank Soal">
    <div className="grid grid-cols-1 gap-4 md:grid-cols-2 lg:grid-cols-4">
      <div>
        <label htmlFor="id_kelas" className="text-xs font-medium text-slate-600">
          Tingkat Kelas
        </label>
        <select
          id="id_kelas"
          value={values.id_kelas}
          onChange={(event) => setField("id_kelas", Number(event.target.value))}
          onBlur={() => onBlur("id_kelas")}
          disabled={loadingKelas || submitting}
          className={`${selectBaseClass} ${hasError("id_kelas") ? "border-rose-300 ring-rose-100" : ""}`}
        >
          <option value={0}>
            {loadingKelas ? "Memuat tingkat kelas..." : "Pilih tingkat kelas"}
          </option>
          {tingkatKelasOptions.map((item) => (
            <option key={item.id_tingkat_kelas} value={item.id_tingkat_kelas}>
              Kelas {item.tingkat_kelas}
            </option>
          ))}
        </select>
        {hasError("id_kelas") && (
          <p className="mt-1 text-xs text-rose-600">{errors.id_kelas}</p>
        )}
      </div>

      <div>
        <label
          htmlFor="kelas_scope"
          className="text-xs font-medium text-slate-600"
        >
          Cakupan Kelas
        </label>
        <select
          id="kelas_scope"
          value={values.kelas_scope}
          onChange={(event) =>
            setField(
              "kelas_scope",
              event.target.value === "SPESIFIK" ? "SPESIFIK" : "SEMUA",
            )
          }
          onBlur={() => onBlur("kelas_scope")}
          disabled={values.id_kelas === 0 || submitting}
          className={`${selectBaseClass} ${hasError("kelas_scope") ? "border-rose-300 ring-rose-100" : ""}`}
        >
          <option value="SEMUA">Semua kelas di tingkat ini</option>
          <option value="SPESIFIK">Spesifik nama kelas</option>
        </select>
        {hasError("kelas_scope") && (
          <p className="mt-1 text-xs text-rose-600">{errors.kelas_scope}</p>
        )}
      </div>

      <div>
        <label
          htmlFor="id_nama_kelas"
          className="text-xs font-medium text-slate-600"
        >
          Nama Kelas
        </label>
        <select
          id="id_nama_kelas"
          value={values.id_nama_kelas}
          onChange={(event) =>
            setField("id_nama_kelas", Number(event.target.value))
          }
          onBlur={() => onBlur("id_nama_kelas")}
          disabled={
            values.id_kelas === 0 ||
            values.kelas_scope !== "SPESIFIK" ||
            submitting
          }
          className={`${selectBaseClass} ${hasError("id_nama_kelas") ? "border-rose-300 ring-rose-100" : ""}`}
        >
          <option value={0}>Pilih nama kelas</option>
          {namaKelasOptions.map((item) => (
            <option key={item.id_nama_kelas} value={item.id_nama_kelas}>
              {item.nama_kelas}
            </option>
          ))}
        </select>
        {hasError("id_nama_kelas") && (
          <p className="mt-1 text-xs text-rose-600">{errors.id_nama_kelas}</p>
        )}
      </div>

      <div className="md:col-span-2 lg:col-span-1">
        <label
          htmlFor="id_mapel_filter"
          className="text-xs font-medium text-slate-600"
        >
          Mapel
        </label>
        <select
          id="id_mapel_filter"
          value={selectedMapelId}
          onChange={(event) => {
            setSelectedMapelId(Number(event.target.value));
            setField("id_bank_soal", 0);
          }}
          disabled={values.id_kelas === 0 || loadingMapel || submitting}
          className={selectBaseClass}
        >
          <option value={0}>
            {values.id_kelas === 0
              ? "Pilih tingkat kelas dulu"
              : loadingMapel
                ? "Memuat mapel..."
                : "Pilih mapel"}
          </option>
          {mapelOptions.map((item) => (
            <option key={item.id} value={item.id}>
              {item.namaMapel}
            </option>
          ))}
        </select>
      </div>

      <div className="md:col-span-2 lg:col-span-2">
        <BankSoalSelect
          id="id_bank_soal"
          label="Bank Soal"
          value={values.id_bank_soal}
          options={bankSoalOptions}
          onChange={(selectedId) => setField("id_bank_soal", selectedId)}
          onBlur={() => onBlur("id_bank_soal")}
          disabled={values.id_kelas === 0 || selectedMapelId === 0 || submitting}
          loading={loadingBankSoal}
          placeholder={bankSoalPlaceholder}
          error={hasError("id_bank_soal")}
        />
        {hasError("id_bank_soal") && (
          <p className="mt-1 text-xs text-rose-600">{errors.id_bank_soal}</p>
        )}
      </div>
    </div>
  </FormSection>
);

export default KelasBankSoalSection;
