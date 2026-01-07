type Primitive = string | number | boolean | null | undefined;
type FormDataValue = Primitive | File | Blob;

type BuildFormDataOptions = {
  /**
   * Kalau true: skip null/undefined (default).
   * Kalau false: null/undefined tetap di-append sebagai string ("").
   */
  skipNullish?: boolean;

  /**
   * Transform value per-key sebelum append.
   * Cocok untuk enum mapping, trimming, lowercasing, dsb.
   */
  transform?: (
    key: string,
    value: unknown
  ) => FormDataValue | FormDataValue[] | null | undefined;

  /**
   * Override nama key di payload (opsional).
   * Cocok kalau API minta snake_case, sementara UI camelCase.
   */
  keyMap?: Record<string, string>;
};

export function buildFormData<T extends Record<string, unknown>>(
  values: T,
  opts: BuildFormDataOptions = {}
): FormData {
  const fd = new FormData();
  const { skipNullish = true, transform, keyMap } = opts;

  for (const [rawKey, rawValue] of Object.entries(values)) {
    const key = keyMap?.[rawKey] ?? rawKey;

    let value: any = rawValue;
    if (transform) value = transform(rawKey, rawValue);

    if (value == null) {
      if (!skipNullish) fd.append(key, "");
      continue;
    }

    // Support array (misal multiple files / tags)
    if (Array.isArray(value)) {
      for (const item of value) {
        if (item == null) continue;
        fd.append(key, toFormDataAcceptable(item));
      }
      continue;
    }

    fd.append(key, toFormDataAcceptable(value));
  }

  return fd;
}

function toFormDataAcceptable(value: unknown): string | Blob {
  if (value instanceof Blob) return value; // File juga termasuk Blob
  if (typeof value === "string") return value;
  if (typeof value === "number") return String(value);
  if (typeof value === "boolean") return value ? "true" : "false";
  if (value instanceof Date) return value.toISOString();

 
  throw new Error(
    `buildFormData: Unsupported value type: ${Object.prototype.toString.call(
      value
    )}. ` + `Flatten/mapping dulu atau pakai transform() untuk stringify.`
  );
}
