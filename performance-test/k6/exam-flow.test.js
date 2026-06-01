import http from "k6/http";
import { check, fail, sleep } from "k6";
import { Rate, Trend } from "k6/metrics";
import { SharedArray } from "k6/data";
import {
  getAnswerLimit,
  getBaseUrl,
  getOptions,
  getStudentsFile,
  getTokenUjian,
} from "./config.js";

export const examFlowFailed = new Rate("exam_flow_failed");
export const examFlowDuration = new Trend("exam_flow_duration", true);
export const options = getOptions();

const BASE_URL = getBaseUrl();
const TOKEN_UJIAN = getTokenUjian();
const ANSWER_LIMIT = getAnswerLimit();
const MAX_LOG_BODY_LENGTH = 2000;

const students = new SharedArray("students", () => {
  const rows = JSON.parse(open(getStudentsFile()));
  if (!Array.isArray(rows) || rows.length === 0) {
    throw new Error("students.json must contain at least one student account");
  }
  return rows;
});

function jsonHeaders(extra = {}) {
  return {
    headers: {
      "Content-Type": "application/json",
      Accept: "application/json",
      ...extra,
    },
  };
}

function logResponse(label, response) {
  const body = String(response.body || "");
  const clippedBody =
    body.length > MAX_LOG_BODY_LENGTH
      ? `${body.slice(0, MAX_LOG_BODY_LENGTH)}... [truncated ${body.length - MAX_LOG_BODY_LENGTH} chars]`
      : body;

  console.log(`[${label}] status=${response.status} body=${clippedBody}`);
}

function parseJson(response, label) {
  try {
    return response.json();
  } catch (error) {
    throw new Error(`${label}: response is not valid JSON (${error.message})`);
  }
}

function unwrapData(response, label) {
  const body = parseJson(response, label);
  if (body && typeof body === "object" && "error" in body && body.error) {
    throw new Error(`${label}: ${body.error.code || "ERROR"} ${body.error.message || ""}`.trim());
  }
  return body && typeof body === "object" && "data" in body ? body.data : body;
}

function requireOk(response, label, expectedStatuses = [200]) {
  const ok = check(response, {
    [`${label}: status ${expectedStatuses.join(" or ")}`]: (res) =>
      expectedStatuses.includes(res.status),
    [`${label}: no 503`]: (res) => res.status !== 503,
  });

  if (!ok) {
    throw new Error(`${label}: unexpected status ${response.status}`);
  }
}

function withQuery(path, params) {
  const query = Object.entries(params)
    .filter(([, value]) => value !== undefined && value !== null && value !== "")
    .map(([key, value]) => `${encodeURIComponent(key)}=${encodeURIComponent(String(value))}`)
    .join("&");

  return query ? `${path}?${query}` : path;
}

function pickStudent() {
  const index = (__VU - 1) % students.length;
  const student = students[index];

  if (!student.username || !student.password) {
    fail(`Student fixture at index ${index} must contain username and password`);
  }

  return student;
}

function loginStudent(student) {
  const loginRes = http.post(
    `${BASE_URL}/auth/login`,
    JSON.stringify({
      username: student.username,
      password: student.password,
    }),
    jsonHeaders(),
  );
  logResponse("loginStudent", loginRes);
  requireOk(loginRes, "loginStudent");

  const meRes = http.get(`${BASE_URL}/auth/me`, jsonHeaders());
  logResponse("getAuthMe", meRes);
  requireOk(meRes, "getAuthMe");

  const user = unwrapData(meRes, "getAuthMe");
  check(user, {
    "getAuthMe: user id available": (value) =>
      Number.isInteger(value?.id_pengguna) && value.id_pengguna > 0,
    "getAuthMe: role is siswa": (value) => String(value?.role || "").toUpperCase() === "SISWA",
  });

  if (!Number.isInteger(user?.id_pengguna) || user.id_pengguna <= 0) {
    throw new Error("getAuthMe: id_pengguna is missing");
  }

  return user;
}

function getActiveExamSchedule() {
  const response = http.get(`${BASE_URL}/siswa/ujian/list`, jsonHeaders());
  logResponse("getActiveExamSchedule", response);
  requireOk(response, "getActiveExamSchedule");

  const schedules = unwrapData(response, "getActiveExamSchedule");
  check(schedules, {
    "getActiveExamSchedule: response is array": (value) => Array.isArray(value),
    "getActiveExamSchedule: active exam exists": (value) =>
      Array.isArray(value) &&
      value.some((item) => {
        const status = String(item.status_ujian || "").toUpperCase();
        return status === "BERLANGSUNG" || status === "MULAI" || status === "BELUM_DIMULAI";
      }),
  });

  const preferredId = Number(__ENV.ID_JADWAL_UJIAN || 0);
  const byEnv =
    preferredId > 0 ? schedules.find((item) => Number(item.id) === preferredId) : undefined;
  const schedule =
    byEnv ||
    schedules.find((item) => ["BERLANGSUNG", "MULAI"].includes(String(item.status_ujian || "").toUpperCase())) ||
    schedules.find((item) => String(item.status_ujian || "").toUpperCase() === "BELUM_DIMULAI") ||
    schedules[0];

  if (!schedule || !Number.isInteger(schedule.id) || schedule.id <= 0) {
    throw new Error("getActiveExamSchedule: no valid id_jadwal_ujian found");
  }

  return schedule;
}

function startExamAttempt(user, schedule) {
  if (!TOKEN_UJIAN) {
    throw new Error("TOKEN_UJIAN environment variable is required");
  }

  const response = http.post(
    `${BASE_URL}/siswa/ujian/attempt`,
    JSON.stringify({
      id_siswa: user.id_pengguna,
      id_jadwal_ujian: schedule.id,
      token_ujian: TOKEN_UJIAN,
      waktu_mulai: new Date().toISOString(),
    }),
    jsonHeaders(),
  );

  logResponse("startExamAttempt", response);
  requireOk(response, "startExamAttempt", [200, 201, 409]);

  if (response.status === 409) {
    const body = parseJson(response, "startExamAttempt");
    const code = body?.error?.code;
    if (code !== "DOUBLE_ATTEMPT_NOT_ALLOWED") {
      throw new Error(`startExamAttempt: conflict ${code || "UNKNOWN"}`);
    }
  }

  const activeAttemptPath = withQuery("/siswa/ujian/attempt/active", {
    id_jadwal_ujian: schedule.id,
  });
  const activeResponse = http.get(`${BASE_URL}${activeAttemptPath}`, jsonHeaders());
  logResponse("getActiveAttempt", activeResponse);
  requireOk(activeResponse, "getActiveAttempt");

  const activeAttempt = unwrapData(activeResponse, "getActiveAttempt");
  check(activeAttempt, {
    "getActiveAttempt: id_attempt available": (value) =>
      Number.isInteger(value?.id_attempt) && value.id_attempt > 0,
  });

  if (!Number.isInteger(activeAttempt?.id_attempt) || activeAttempt.id_attempt <= 0) {
    throw new Error("getActiveAttempt: id_attempt is missing");
  }

  return activeAttempt.id_attempt;
}

function getExamQuestions(schedule) {
  const response = http.get(`${BASE_URL}/siswa/soal-ujian/${schedule.id}`, jsonHeaders());
  logResponse("getExamQuestions", response);
  requireOk(response, "getExamQuestions");

  const questions = unwrapData(response, "getExamQuestions");
  check(questions, {
    "getExamQuestions: response is array": (value) => Array.isArray(value),
    "getExamQuestions: questions available": (value) => Array.isArray(value) && value.length > 0,
    "getExamQuestions: each question has id": (value) =>
      Array.isArray(value) && value.every((item) => Number.isInteger(item.id_soal) && item.id_soal > 0),
  });

  if (!Array.isArray(questions) || questions.length === 0) {
    throw new Error("getExamQuestions: no questions returned");
  }

  return questions;
}

function buildAnswer(question, attemptId) {
  const now = new Date().toISOString();
  const options = Array.isArray(question.opsi_jawaban) ? question.opsi_jawaban : [];

  if (options.length > 0) {
    const option = options[(__VU + __ITER + question.no_urut_soal) % options.length];
    return {
      id_attempt: attemptId,
      jawaban: [
        {
          id_soal: question.id_soal,
          id_pilihan: option.id_pilihan_ganda,
          jawaban_essay: null,
          waktu_jawab: now,
        },
      ],
    };
  }

  return {
    id_attempt: attemptId,
    jawaban: [
      {
        id_soal: question.id_soal,
        id_pilihan: null,
        jawaban_essay: `Jawaban essay performa VU ${__VU} iterasi ${__ITER}`,
        waktu_jawab: now,
      },
    ],
  };
}

function saveAnswers(attemptId, questions) {
  const selectedQuestions = questions.slice(0, Math.min(ANSWER_LIMIT, questions.length));

  for (const question of selectedQuestions) {
    const response = http.post(
      `${BASE_URL}/siswa/ujian/jawaban`,
      JSON.stringify(buildAnswer(question, attemptId)),
      jsonHeaders(),
    );
    logResponse("saveAnswers", response);
    requireOk(response, `saveAnswers question ${question.id_soal}`, [200, 201]);
    sleep(1 + Math.floor(Math.random() * 3));
  }
}

function submitExam(attemptId) {
  const response = http.patch(`${BASE_URL}/siswa/uijan/submit/${attemptId}`, null, jsonHeaders());
  logResponse("submitExam", response);
  requireOk(response, "submitExam", [200, 201, 409]);

  if (response.status === 409) {
    const body = parseJson(response, "submitExam");
    const code = body?.error?.code;
    if (code !== "DOUBLE_SUBMIT_NOT_ALLOWED") {
      throw new Error(`submitExam: conflict ${code || "UNKNOWN"}`);
    }
  }
}

function validateResult(attemptId) {
  const response = http.get(`${BASE_URL}/ujian/jawaban/hasil/${attemptId}`, jsonHeaders());
  logResponse("validateResult", response);
  requireOk(response, "validateResult", [200, 404]);
}

export default function () {
  const startedAt = Date.now();

  try {
    const student = pickStudent();
    const user = loginStudent(student);
    sleep(1);

    const schedule = getActiveExamSchedule();
    sleep(1);

    const attemptId = startExamAttempt(user, schedule);
    sleep(1);

    const questions = getExamQuestions(schedule);
    saveAnswers(attemptId, questions);
    sleep(1);

    submitExam(attemptId);
    validateResult(attemptId);

    examFlowFailed.add(0);
  } catch (error) {
    examFlowFailed.add(1);
    throw error;
  } finally {
    examFlowDuration.add(Date.now() - startedAt);
  }
}
