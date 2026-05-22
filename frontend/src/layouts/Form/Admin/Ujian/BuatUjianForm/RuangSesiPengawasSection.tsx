import { selectBaseClass } from "./constants";
import FormSection from "./FormSection";
import type { BuatUjianFormController } from "./useBuatUjianForm";

type RuangSesiPengawasSectionProps = Pick<
  BuatUjianFormController,
  | "errors"
  | "guruOptions"
  | "hasError"
  | "loadingGuru"
  | "loadingRuang"
  | "loadingSesi"
  | "onBlur"
  | "ruangOptions"
  | "sesiOptions"
  | "setField"
  | "submitting"
  | "values"
>;

const RuangSesiPengawasSection = ({
  errors,
  guruOptions,
  hasError,
  loadingGuru,
  loadingRuang,
  loadingSesi,
  onBlur,
  ruangOptions,
  sesiOptions,
  setField,
  submitting,
  values,
}: RuangSesiPengawasSectionProps) => (
  <FormSection
    title="Ruang, Sesi, Pengawas"
    description="Semua opsi diambil langsung dari endpoint server."
  >
    <div className="grid grid-cols-1 gap-4 md:grid-cols-3">
      <div>
        <label
          htmlFor="id_ruangan"
          className="text-xs font-medium text-slate-600"
        >
          Ruang Ujian
        </label>
        <select
          id="id_ruangan"
          value={values.id_ruangan}
          onChange={(event) => setField("id_ruangan", Number(event.target.value))}
          onBlur={() => onBlur("id_ruangan")}
          disabled={loadingRuang || submitting}
          className={`${selectBaseClass} ${hasError("id_ruangan") ? "border-rose-300 ring-rose-100" : ""}`}
        >
          <option value={0}>
            {loadingRuang ? "Memuat ruang..." : "Pilih ruang ujian"}
          </option>
          {ruangOptions.map((item) => (
            <option key={item.id_ruangan} value={item.id_ruangan}>
              {item.nama_ruangan}
            </option>
          ))}
        </select>
        {hasError("id_ruangan") && (
          <p className="mt-1 text-xs text-rose-600">{errors.id_ruangan}</p>
        )}
      </div>

      <div>
        <label htmlFor="id_sesi" className="text-xs font-medium text-slate-600">
          Sesi Ujian
        </label>
        <select
          id="id_sesi"
          value={values.id_sesi}
          onChange={(event) => setField("id_sesi", Number(event.target.value))}
          onBlur={() => onBlur("id_sesi")}
          disabled={loadingSesi || submitting}
          className={`${selectBaseClass} ${hasError("id_sesi") ? "border-rose-300 ring-rose-100" : ""}`}
        >
          <option value={0}>
            {loadingSesi ? "Memuat sesi..." : "Pilih sesi ujian"}
          </option>
          {sesiOptions.map((item) => (
            <option key={item.id_sesi} value={item.id_sesi}>
              {item.kode_sesi} - {item.nama_sesi}
            </option>
          ))}
        </select>
        {hasError("id_sesi") && (
          <p className="mt-1 text-xs text-rose-600">{errors.id_sesi}</p>
        )}
      </div>

      <div>
        <label
          htmlFor="id_pengawas"
          className="text-xs font-medium text-slate-600"
        >
          Guru Pengawas
        </label>
        <select
          id="id_pengawas"
          value={values.id_pengawas}
          onChange={(event) => setField("id_pengawas", Number(event.target.value))}
          onBlur={() => onBlur("id_pengawas")}
          disabled={loadingGuru || submitting}
          className={`${selectBaseClass} ${hasError("id_pengawas") ? "border-rose-300 ring-rose-100" : ""}`}
        >
          <option value={0}>
            {loadingGuru ? "Memuat guru..." : "Pilih guru pengawas"}
          </option>
          {guruOptions.map((item) => (
            <option key={item.id_pengguna} value={item.id_pengguna}>
              {item.nama_lengkap}
            </option>
          ))}
        </select>
        {hasError("id_pengawas") && (
          <p className="mt-1 text-xs text-rose-600">{errors.id_pengawas}</p>
        )}
      </div>
    </div>
  </FormSection>
);

export default RuangSesiPengawasSection;
