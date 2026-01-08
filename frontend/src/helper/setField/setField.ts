import type React from "react";

export function createSetField<T extends Record<string, unknown>>(
  setValues: React.Dispatch<React.SetStateAction<T>>
) {
  return function setField<K extends keyof T>(key: K, value: T[K]) {
    setValues((prev) => {
      if (Object.is(prev[key], value)) return prev;
      return { ...prev, [key]: value };
    });
  };
}
