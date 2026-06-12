export const DEFAULT_BASE_URL = "https://staging-srv.smafi.my.id/";
export const DEFAULT_STUDENTS_FILE = "./data/students.json";
export const DEFAULT_ANSWER_LIMIT = 10;

const TEST_SCENARIOS = {
  baseline_5: {
    executor: "constant-vus",
    vus: 5,
    duration: "5m",
    gracefulStop: "30s",
  },
  load_25: {
    executor: "constant-vus",
    vus: 25,
    duration: "10m",
    gracefulStop: "30s",
  },
  load_50: {
    executor: "constant-vus",
    vus: 50,
    duration: "10m",
    gracefulStop: "30s",
  },
  load_100: {
    executor: "constant-vus",
    vus: 100,
    duration: "10m",
    gracefulStop: "30s",
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
  const testType = (__ENV.TEST_TYPE || "load_100").trim();
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
