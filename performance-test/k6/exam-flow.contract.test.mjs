import assert from "node:assert/strict";
import { readFileSync } from "node:fs";
import { test } from "node:test";

const configSource = readFileSync(new URL("./config.js", import.meta.url), "utf8");
const examFlowSource = readFileSync(
  new URL("./exam-flow.test.js", import.meta.url),
  "utf8",
);
const students = JSON.parse(
  readFileSync(new URL("./data/students.json", import.meta.url), "utf8"),
);
const testTypes = ["baseline_10", "baseline_20", "load_25", "load_50", "load_100"];

test("config exposes all planned scenarios and thresholds", () => {
  for (const testType of testTypes) {
    assert.match(configSource, new RegExp(`${testType}:\\s*{`));
  }

  assert.match(configSource, /http_req_duration:\s*\["p\(95\)<4000"\]/);
  assert.match(configSource, /http_req_failed:\s*\["rate<0\.01"\]/);
  assert.match(configSource, /checks:\s*\["rate>0\.99"\]/);
  assert.match(configSource, /exam_flow_failed:\s*\["rate<0\.01"\]/);
});

test("each scenario logs in once per VU and runs one exam flow", () => {
  for (const testType of testTypes) {
    const scenarioBlock = configSource.slice(
      configSource.indexOf(`${testType}: {`),
      configSource.indexOf("},", configSource.indexOf(`${testType}: {`)) + 2,
    );

    assert.match(scenarioBlock, /executor:\s*"per-vu-iterations"/);
    assert.match(scenarioBlock, /iterations:\s*1/);
  }
});

test("exam flow uses actual backend routes and required steps", () => {
  for (const route of [
    "/auth/login",
    "/auth/me",
    "/siswa/ujian/list",
    "/siswa/ujian/attempt",
    "/siswa/ujian/attempt/active",
    "/siswa/soal-ujian/",
    "/siswa/ujian/jawaban",
    "/siswa/uijan/submit/",
    "/ujian/jawaban/hasil/",
  ]) {
    assert.ok(examFlowSource.includes(route), `expected ${route} route`);
  }

  for (const fnName of [
    "loginStudent",
    "getActiveExamSchedule",
    "startExamAttempt",
    "getExamQuestions",
    "saveAnswers",
    "submitExam",
  ]) {
    assert.match(examFlowSource, new RegExp(`function ${fnName}\\(`));
  }
});

test("exam flow logs backend responses to console", () => {
  assert.match(examFlowSource, /function logResponse\(label,\s*response\)/);

  for (const label of [
    "loginStudent",
    "getAuthMe",
    "getActiveExamSchedule",
    "startExamAttempt",
    "getActiveAttempt",
    "getExamQuestions",
    "saveAnswers",
    "submitExam",
    "validateResult",
  ]) {
    assert.match(
      examFlowSource,
      new RegExp(`logResponse\\("${label}",\\s*\\w+\\)`),
      `expected ${label} response logging`,
    );
  }
});

test("students fixture contains at least 100 unique students", () => {
  assert.ok(students.length >= 100);
  assert.equal(new Set(students.map((student) => student.username)).size, students.length);

  for (const student of students) {
    assert.equal(typeof student.username, "string");
    assert.equal(typeof student.password, "string");
    assert.ok(student.username.length > 0);
    assert.ok(student.password.length > 0);
  }
});
