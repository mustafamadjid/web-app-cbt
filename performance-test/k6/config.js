export const DEFAULT_BASE_URL = "https://staging-srv.smafi.my.id/";
export const DEFAULT_STUDENTS_FILE = "./data/students.json";
export const DEFAULT_ANSWER_LIMIT = 10;

const TEST_SCENARIOS = {
  baseline_10: {
    executor: "per-vu-iterations",
    vus: 10,
    iterations: 1,
    maxDuration: "5m",
  },
  baseline_20: {
    executor: "per-vu-iterations",
    vus: 20,
    iterations: 1,
    maxDuration: "5m",
  },
  load_25: {
    executor: "per-vu-iterations",
    vus: 25,
    iterations: 1,
    maxDuration: "10m",
  },
  load_50: {
    executor: "per-vu-iterations",
    vus: 50,
    iterations: 1,
    maxDuration: "10m",
  },
  load_100: {
    executor: "per-vu-iterations",
    vus: 100,
    iterations: 1,
    maxDuration: "10m",
  },
};

export function getBaseUrl() {
  const raw = (__ENV.BASE_URL || DEFAULT_BASE_URL).trim();
  return raw.replace(/\/+$/, "");
}

export function getStudentsFile() {
  return (__ENV.STUDENTS_FILE || DEFAULT_STUDENTS_FILE).trim();
}

export function getTokenUjian() {
  return (__ENV.TOKEN_UJIAN || "").trim();
}

export function getAnswerLimit() {
  const raw = Number(__ENV.ANSWER_LIMIT || DEFAULT_ANSWER_LIMIT);
  if (!Number.isInteger(raw) || raw <= 0) return DEFAULT_ANSWER_LIMIT;
  return raw;
}

export function getOptions() {
  const testType = (__ENV.TEST_TYPE || "baseline_10").trim();
  const scenario = TEST_SCENARIOS[testType];

  if (!scenario) {
    throw new Error(
      `Unsupported TEST_TYPE "${testType}". Use one of: ${Object.keys(TEST_SCENARIOS).join(", ")}`,
    );
  }

  return {
    scenarios: {
      [testType]: scenario,
    },
    thresholds: {
      http_req_duration: ["p(95)<4000"],
      http_req_failed: ["rate<0.01"],
      checks: ["rate>0.99"],
      exam_flow_failed: ["rate<0.01"],
    },
  };
}
