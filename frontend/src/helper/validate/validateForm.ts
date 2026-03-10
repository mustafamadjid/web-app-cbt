export type ValidationRule<Value, Values> = (
  value: Value,
  values: Values
) => string | null;

type ValidationSchema<Values extends Record<string, unknown>> = {
  [Key in keyof Values]?: Array<ValidationRule<Values[Key], Values>>;
};

const isEmptyValue = (value: unknown) =>
  value === "" || value === null || value === undefined;

export const createValidator = <Values extends Record<string, unknown>>(
  schema: ValidationSchema<Values>
) => {
  return (values: Values) => {
    const errors: Partial<Record<keyof Values, string>> = {};

    (Object.keys(schema) as Array<keyof Values>).forEach((field) => {
      const rules = schema[field];
      if (!rules || rules.length === 0) return;

      const value = values[field];
      for (const rule of rules) {
        const message = rule(value, values);
        if (message) {
          errors[field] = message;
          break;
        }
      }
    });

    return errors;
  };
};

export const requiredValue = <Value>(
  message: string
): ValidationRule<Value, Record<string, unknown>> => {
  return (value) => (isEmptyValue(value) ? message : null);
};

export const requiredString = (
  message: string
): ValidationRule<string, Record<string, unknown>> => {
  return (value) => (!value || value.trim() === "" ? message : null);
};

export const emailFormat = (
  message: string
): ValidationRule<string, Record<string, unknown>> => {
  return (value) => {
    if (isEmptyValue(value)) return null;
    return /^\S+@\S+\.\S+$/.test(value) ? null : message;
  };
};

export const minLength = (
  min: number,
  message: string
): ValidationRule<string, Record<string, unknown>> => {
  return (value) => {
    if (isEmptyValue(value)) return null;
    return value.length < min ? message : null;
  };
};

export const maxLength = (
  max: number,
  message: string
): ValidationRule<string, Record<string, unknown>> => {
  return (value) => {
    if (isEmptyValue(value)) return null;
    return value.length > max ? message : null;
  };
};

export const integerNumber = (
  message: string
): ValidationRule<number | "", Record<string, unknown>> => {
  return (value) => {
    if (isEmptyValue(value)) return null;
    return Number.isInteger(value) ? null : message;
  };
};

export const minNumber = (
  min: number,
  message: string
): ValidationRule<number | "", Record<string, unknown>> => {
  return (value) => {
    if (isEmptyValue(value)) return null;
    return Number(value) < min ? message : null;
  };
};

export const matchesPattern = (
  pattern: RegExp,
  message: string,
  options?: { trim?: boolean }
): ValidationRule<string, Record<string, unknown>> => {
  return (value) => {
    if (isEmptyValue(value)) return null;
    const target = options?.trim ? value.trim() : value;
    return pattern.test(target) ? null : message;
  };
};

export const fileMaxSize = (
  maxBytes: number,
  message: string
): ValidationRule<File | null, Record<string, unknown>> => {
  return (value) =>
    value && value.size > maxBytes ? message : null;
};

export const fileTypeStartsWith = (
  prefix: string,
  message: string
): ValidationRule<File | null, Record<string, unknown>> => {
  return (value) =>
    value && !value.type.startsWith(prefix) ? message : null;
};

export const fileDocxOnly = (
  message: string
): ValidationRule<File | null, Record<string, unknown>> => {
  const DOCX_MIME =
    "application/vnd.openxmlformats-officedocument.wordprocessingml.document";

  return (value) => {
    if (!value) return null;

    const nameOk = /\.docx$/i.test(value.name);
    const typeOk = value.type === DOCX_MIME;

    // Banyak kasus: type kosong, jadi fallback ke ekstensi
    return typeOk || nameOk ? null : message;
  };
};
