# Duration-Based Rotating User Load Test Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Mengubah performance test k6 menjadi duration-based rotating user load test dengan baseline 5 VU selama 5 menit dan skenario load 25 VU, 50 VU, dan 100 VU selama 10 menit, serta memastikan setiap iterasi memakai akun siswa yang unik pada satu run.

**Architecture:** Alur ujian siswa tetap memakai `performance-test/k6/exam-flow.test.js`. Konfigurasi skenario dipusatkan di `performance-test/k6/config.js` dengan executor `constant-vus`. Rotasi akun dilakukan di script k6 menggunakan indeks iterasi global dari `k6/execution`, sehingga akun tidak dipilih dari `__VU` dan tidak dipakai ulang selama run yang sama.

**Tech Stack:** k6, JavaScript ES modules, `SharedArray`, `k6/execution`, Node.js built-in test runner untuk contract test.

---

## Requirements

- Jenis test: duration-based rotating user load test.
- Skenario baseline: `baseline_5` dengan 5 VU selama 5 menit.
- Skenario load: `load_25`, `load_50`, dan `load_100`.
- Durasi skenario load: 10 menit.
- Skenario flow ujian tetap sama seperti script existing:
  - login siswa.
  - ambil `/auth/me`.
  - ambil list ujian siswa.
  - mulai atau ambil attempt aktif.
  - ambil soal.
  - simpan jawaban.
  - submit ujian.
  - validasi hasil.
- Data akun siswa dibaca dari `performance-test/k6/data/students.json`.
- Setiap iterasi memakai akun siswa yang berbeda.
- Akun yang sudah digunakan tidak boleh dipakai ulang pada run yang sama.
- Jika jumlah akun tidak cukup, test harus berhenti dengan error jelas.
- Command utama harus dapat menjalankan 100 VU selama 10 menit.

## File Structure

- Modify: `performance-test/k6/config.js`
  - Tambahkan skenario `baseline_5` berbasis `constant-vus` dengan `vus: 5` dan `duration: "5m"`.
  - Ubah skenario `load_25`, `load_50`, dan `load_100` menjadi `constant-vus` dengan `duration: "10m"`.
  - Default `TEST_TYPE` diarahkan ke `load_100` agar script utama menjalankan 100 VU selama 10 menit.

- Modify: `performance-test/k6/exam-flow.test.js`
  - Import `exec` dari `k6/execution`.
  - Ganti `pickStudent()` agar memakai `exec.scenario.iterationInTest`.
  - Tambahkan validasi schema akun pada saat load `students.json`.
  - Tambahkan abort jelas jika indeks iterasi global melebihi jumlah akun.

- Modify: `performance-test/k6/exam-flow.contract.test.mjs`
  - Update contract scenario dari `per-vu-iterations` menjadi `constant-vus`.
  - Tambahkan assertion bahwa script memakai `k6/execution` dan `iterationInTest`.
  - Tambahkan assertion bahwa error akun habis jelas dan tidak fallback ke akun lama.

- Modify: `performance-test/README.md`
  - Update deskripsi jenis test.
  - Update command run untuk baseline 5 VU, 25 VU, 50 VU, dan 100 VU.
  - Jelaskan kebutuhan jumlah akun minimal sama dengan jumlah total iterasi selama 10 menit, bukan hanya jumlah VU.

## Design Notes

- Gunakan `exec.scenario.iterationInTest` sebagai indeks global iterasi. Nilai ini unik per iterasi di dalam satu skenario, sehingga cocok untuk mapping satu iterasi ke satu akun.
- Jangan gunakan `__VU`, karena `__VU` hanya unik per virtual user dan akan mengulang akun yang sama pada iterasi berikutnya.
- Jangan gunakan `__ITER`, karena `__ITER` hanya unik per VU dan akan bentrok antar VU.
- Pada `constant-vus`, jumlah total iterasi tidak diketahui sebelum run karena bergantung pada durasi satu flow. Karena itu validasi akun dilakukan saat setiap iterasi mengambil akun.
- Saat akun habis, panggil `exec.test.abort(message)` supaya run berhenti, bukan hanya satu iterasi gagal.

## Implementation Tasks

### Task 1: Update k6 Scenario Config

**Files:**
- Modify: `performance-test/k6/config.js`
- Test: `performance-test/k6/exam-flow.contract.test.mjs`

- [ ] **Step 1: Write failing contract expectations**

Update `testTypes` agar hanya mencakup skenario duration-based:

```javascript
const testTypes = ["baseline_5", "load_25", "load_50", "load_100"];
```

Update test scenario menjadi:

```javascript
test("config exposes duration-based rotating load scenarios", () => {
  for (const testType of testTypes) {
    assert.match(configSource, new RegExp(`${testType}:\\s*{`));
  }

  assert.match(configSource, /baseline_5:\s*{[\s\S]*executor:\s*"constant-vus"[\s\S]*vus:\s*5[\s\S]*duration:\s*"5m"/);
  assert.match(configSource, /load_25:\s*{[\s\S]*executor:\s*"constant-vus"[\s\S]*vus:\s*25[\s\S]*duration:\s*"10m"/);
  assert.match(configSource, /load_50:\s*{[\s\S]*executor:\s*"constant-vus"[\s\S]*vus:\s*50[\s\S]*duration:\s*"10m"/);
  assert.match(configSource, /load_100:\s*{[\s\S]*executor:\s*"constant-vus"[\s\S]*vus:\s*100[\s\S]*duration:\s*"10m"/);
  assert.match(configSource, /__ENV\.TEST_TYPE\s*\|\|\s*"load_100"/);
});
```

- [ ] **Step 2: Run contract test and verify failure**

Run:

```bash
cd performance-test/k6
node --test exam-flow.contract.test.mjs
```

Expected: FAIL karena `config.js` belum memiliki `baseline_5`, masih memakai `per-vu-iterations`, `iterations: 1`, dan default `baseline_10`.

- [ ] **Step 3: Update `TEST_SCENARIOS`**

Replace scenario config with:

```javascript
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
```

Update default test type:

```javascript
const testType = (__ENV.TEST_TYPE || "load_100").trim();
```

- [ ] **Step 4: Run contract test and verify scenario config passes**

Run:

```bash
cd performance-test/k6
node --test exam-flow.contract.test.mjs
```

Expected: scenario config assertions pass. Other tests may still fail until Task 2 updates account rotation.

- [ ] **Step 5: Commit**

```bash
git add performance-test/k6/config.js performance-test/k6/exam-flow.contract.test.mjs
git commit -m "test(k6): switch load scenarios to duration based"
```

### Task 2: Add Rotating User Account Mechanism

**Files:**
- Modify: `performance-test/k6/exam-flow.test.js`
- Test: `performance-test/k6/exam-flow.contract.test.mjs`

- [ ] **Step 1: Write failing contract for rotation**

Add a contract test:

```javascript
test("exam flow rotates accounts by global iteration and aborts when accounts are exhausted", () => {
  assert.match(examFlowSource, /import exec from "k6\/execution"/);
  assert.match(examFlowSource, /exec\.scenario\.iterationInTest/);
  assert.doesNotMatch(examFlowSource, /const index = \(__VU - 1\) % students\.length/);
  assert.doesNotMatch(examFlowSource, /% students\.length/);
  assert.match(examFlowSource, /exec\.test\.abort/);
  assert.match(examFlowSource, /Not enough student accounts/);
});
```

- [ ] **Step 2: Run contract test and verify failure**

Run:

```bash
cd performance-test/k6
node --test exam-flow.contract.test.mjs
```

Expected: FAIL karena `pickStudent()` masih memakai `(__VU - 1) % students.length`.

- [ ] **Step 3: Import k6 execution API**

At the top of `performance-test/k6/exam-flow.test.js`, add:

```javascript
import exec from "k6/execution";
```

- [ ] **Step 4: Validate `students.json` rows when loaded**

Update SharedArray loader:

```javascript
const students = new SharedArray("students", () => {
  const rows = JSON.parse(open(getStudentsFile()));
  if (!Array.isArray(rows) || rows.length === 0) {
    throw new Error("students.json must contain at least one student account");
  }

  const usernames = new Set();
  rows.forEach((student, index) => {
    if (!student || typeof student.username !== "string" || student.username.trim() === "") {
      throw new Error(`students.json row ${index} must contain a non-empty username`);
    }
    if (typeof student.password !== "string" || student.password.trim() === "") {
      throw new Error(`students.json row ${index} must contain a non-empty password`);
    }
    if (usernames.has(student.username)) {
      throw new Error(`students.json contains duplicate username: ${student.username}`);
    }
    usernames.add(student.username);
  });

  return rows;
});
```

- [ ] **Step 5: Replace `pickStudent()`**

Replace the current implementation with:

```javascript
function pickStudent() {
  const index = exec.scenario.iterationInTest;

  if (index >= students.length) {
    exec.test.abort(
      `Not enough student accounts for this run: iteration ${index + 1} requires a new account, but students.json only contains ${students.length} accounts`,
    );
  }

  return students[index];
}
```

- [ ] **Step 6: Run contract test and verify rotation passes**

Run:

```bash
cd performance-test/k6
node --test exam-flow.contract.test.mjs
```

Expected: PASS for rotation assertions.

- [ ] **Step 7: Commit**

```bash
git add performance-test/k6/exam-flow.test.js performance-test/k6/exam-flow.contract.test.mjs
git commit -m "test(k6): rotate siswa account per iteration"
```

### Task 3: Update Documentation and Run Commands

**Files:**
- Modify: `performance-test/README.md`

- [ ] **Step 1: Update test type description**

Replace the execution model section with:

```markdown
Semua skenario memakai executor `constant-vus`.
`baseline_5` menjalankan 5 VU selama 5 menit.
`load_25`, `load_50`, dan `load_100` menjalankan VU sesuai nama skenario selama 10 menit.
Setiap iterasi mengambil akun siswa yang belum pernah dipakai pada run tersebut.
Jika akun di `k6/data/students.json` habis, test akan berhenti dengan error `Not enough student accounts`.
```

- [ ] **Step 2: Update environment variable table**

Use this row for `TEST_TYPE`:

```markdown
| `TEST_TYPE` | Tidak | `load_100` | `baseline_5`, `load_25`, `load_50`, atau `load_100` |
```

- [ ] **Step 3: Add commands for all scenarios**

Document this as the baseline command:

```bash
k6 run k6/exam-flow.test.js \
  -e BASE_URL=https://staging-srv.smafi.my.id/ \
  -e TOKEN_UJIAN=ABC123 \
  -e TEST_TYPE=baseline_5 \
  --summary-export results/baseline-5.json
```

Document this as the 25 VU load command:

```bash
k6 run k6/exam-flow.test.js \
  -e BASE_URL=https://staging-srv.smafi.my.id/ \
  -e TOKEN_UJIAN=ABC123 \
  -e TEST_TYPE=load_25 \
  --summary-export results/load-25.json
```

Document this as the 50 VU load command:

```bash
k6 run k6/exam-flow.test.js \
  -e BASE_URL=https://staging-srv.smafi.my.id/ \
  -e TOKEN_UJIAN=ABC123 \
  -e TEST_TYPE=load_50 \
  --summary-export results/load-50.json
```

Document this as the primary load command:

```bash
k6 run k6/exam-flow.test.js \
  -e BASE_URL=https://staging-srv.smafi.my.id/ \
  -e TOKEN_UJIAN=ABC123 \
  -e TEST_TYPE=load_100 \
  --summary-export results/load-100.json
```

- [ ] **Step 4: Add account capacity note**

Add:

```markdown
Jumlah akun siswa harus lebih besar atau sama dengan total iterasi yang terjadi selama durasi skenario.
Contoh: jika `load_100` menyelesaikan 2.500 iterasi total selama run, maka minimal dibutuhkan 2.500 akun unik.
Script tidak akan mengulang akun saat akun habis.
```

- [ ] **Step 5: Commit**

```bash
git add performance-test/README.md
git commit -m "docs(k6): document duration based rotating user load test"
```

### Task 4: Final Verification

**Files:**
- Read: `performance-test/k6/config.js`
- Read: `performance-test/k6/exam-flow.test.js`
- Read: `performance-test/k6/exam-flow.contract.test.mjs`
- Read: `performance-test/README.md`

- [ ] **Step 1: Run contract tests**

Run:

```bash
cd performance-test/k6
node --test exam-flow.contract.test.mjs
```

Expected: PASS.

- [ ] **Step 2: Dry-check k6 script syntax**

Run:

```bash
cd performance-test
k6 inspect k6/exam-flow.test.js
```

Expected: k6 prints resolved options with `load_100`, `executor: constant-vus`, `vus: 100`, and `duration: 10m`.

- [ ] **Step 3: Verify student account capacity and uniqueness**

Run:

```bash
cd performance-test
node -e "const fs=require('fs'); const students=JSON.parse(fs.readFileSync('k6/data/students.json','utf8')); const usernames=new Set(students.map((s)=>s.username)); if (students.length !== usernames.size) throw new Error('students.json contains duplicate usernames'); if (students.length < 100) throw new Error('students.json must contain at least 100 accounts for load_100 startup'); console.log('students=' + students.length + ', unique=' + usernames.size);"
```

Expected: command prints equal `students` and `unique` counts. This only validates fixture uniqueness and minimum startup capacity; runtime exhaustion is still guarded by `exec.test.abort()` because duration-based total iterations depends on backend response time.

- [ ] **Step 4: Run the requested 100 VU test**

Run:

```bash
cd performance-test
k6 run k6/exam-flow.test.js \
  -e BASE_URL=https://staging-srv.smafi.my.id/ \
  -e TOKEN_UJIAN=ABC123 \
  -e TEST_TYPE=load_100 \
  --summary-export results/load-100.json
```

Expected: 100 VU run for 10 minutes. If total completed iterations exceeds available accounts, run aborts with:

```text
Not enough student accounts for this run: iteration <n> requires a new account, but students.json only contains <total> accounts
```

## Self-Review

- Spec coverage: baseline 5 VU selama 5 menit, load 25/50/100 VU selama 10 menit, unchanged exam flow, unique account per iteration, and clear error on account exhaustion are mapped to tasks.
- Placeholder scan: no unfinished markers or vague implementation placeholders remain.
- Type consistency: scenario names are consistently `baseline_5`, `load_25`, `load_50`, and `load_100`; account rotation consistently uses `exec.scenario.iterationInTest`.
