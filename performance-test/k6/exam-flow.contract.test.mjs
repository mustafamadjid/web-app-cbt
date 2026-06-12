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
const testTypes = ["baseline_5", "load_25", "load_50", "load_100"];

test("config exposes duration-based rotating load scenarios and thresholds", () => {
  for (const testType of testTypes) {
    assert.match(configSource, new RegExp(`${testType}:\\s*{`));
  }

  assert.match(
    configSource,
    /baseline_5:\s*{[\s\S]*executor:\s*"constant-vus"[\s\S]*vus:\s*5[\s\S]*duration:\s*"5m"/,
  );
  assert.match(
    configSource,
    /load_25:\s*{[\s\S]*executor:\s*"constant-vus"[\s\S]*vus:\s*25[\s\S]*duration:\s*"10m"/,
  );
  assert.match(
    configSource,
    /load_50:\s*{[\s\S]*executor:\s*"constant-vus"[\s\S]*vus:\s*50[\s\S]*duration:\s*"10m"/,
  );
  assert.match(
    configSource,
    /load_100:\s*{[\s\S]*executor:\s*"constant-vus"[\s\S]*vus:\s*100[\s\S]*duration:\s*"10m"/,
  );
  assert.match(configSource, /__ENV\.TEST_TYPE\s*\|\|\s*"load_100"/);
  assert.match(configSource, /http_req_duration:\s*\["p\(95\)<4000"\]/);
  assert.match(configSource, /http_req_failed:\s*\["rate<0\.01"\]/);
  assert.match(configSource, /checks:\s*\["rate>0\.99"\]/);
  assert.match(configSource, /exam_flow_failed:\s*\["rate<0\.01"\]/);
});

test("each scenario runs constant VUs for the configured duration", () => {
  for (const testType of testTypes) {
    const scenarioBlock = configSource.slice(
      configSource.indexOf(`${testType}: {`),
      configSource.indexOf("},", configSource.indexOf(`${testType}: {`)) + 2,
    );

    assert.match(scenarioBlock, /executor:\s*"constant-vus"/);
    assert.match(scenarioBlock, /duration:\s*"(5m|10m)"/);
    assert.doesNotMatch(scenarioBlock, /iterations:\s*1/);
  }
});

test("exam flow rotates accounts by global iteration and aborts when accounts are exhausted", () => {
  assert.match(examFlowSource, /import exec from "k6\/execution"/);
  assert.match(examFlowSource, /exec\.scenario\.iterationInTest/);
  assert.doesNotMatch(examFlowSource, /const index = \(__VU - 1\) % students\.length/);
  assert.doesNotMatch(examFlowSource, /% students\.length/);
  assert.match(examFlowSource, /exec\.test\.abort/);
  assert.match(examFlowSource, /Not enough student accounts/);
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
