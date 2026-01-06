type BuildJsonDataOptions = {
  keyMap?: Record<string, string>;
  nullishToEmptyString?: boolean; 
};

export function buildJsonData<T extends Record<string, unknown>>(
  values: T,
  opts: BuildJsonDataOptions = {}
): Record<string, unknown> {
  const json: Record<string, unknown> = {};
  const { keyMap, nullishToEmptyString = true } = opts;

  for (const [rawKey, rawValue] of Object.entries(values)) {
    const key = keyMap?.[rawKey] ?? rawKey;

    if (rawValue == null) {
      if (nullishToEmptyString) json[key] = "";
      else json[key] = rawValue; // null/undefined dipertahankan
      continue;
    }

    json[key] = rawValue;
  }

  return json;
}
