import dotenv from "dotenv";
import path from "node:path";
import { fileURLToPath } from "node:url";

const __filename = fileURLToPath(import.meta.url);
const __dirname = path.dirname(__filename);
const rootDir = path.resolve(__dirname, "..");

dotenv.config({ path: path.join(rootDir, ".env.e2e") });

const bool = (value, fallback = false) => {
  if (value === undefined) return fallback;
  return String(value).toLowerCase() === "true";
};

const int = (value, fallback) => {
  const parsed = Number.parseInt(value, 10);
  return Number.isFinite(parsed) ? parsed : fallback;
};

export const env = {
  rootDir,
  baseUrl: process.env.E2E_BASE_URL || "http://localhost:5173",
  apiUrl: process.env.E2E_API_URL || "http://localhost:8080",
  browser: process.env.E2E_BROWSER || "chrome",
  headless: bool(process.env.E2E_HEADLESS, true),
  timeout: int(process.env.E2E_TIMEOUT, 15000),
  allowDestructive: bool(process.env.E2E_ALLOW_DESTRUCTIVE, false),
  accounts: {
    admin: {
      username: process.env.E2E_ADMIN_USERNAME || "admin_e2e",
      password: process.env.E2E_ADMIN_PASSWORD || "Password123!",
      dashboardPath: "/dashboard/administrator",
    },
    guru: {
      username: process.env.E2E_GURU_USERNAME || "guru_e2e",
      password: process.env.E2E_GURU_PASSWORD || "Password123!",
      dashboardPath: "/dashboard/guru",
    },
    siswa: {
      username: process.env.E2E_SISWA_USERNAME || "siswa_e2e",
      password: process.env.E2E_SISWA_PASSWORD || "Password123!",
      dashboardPath: "/dashboard/siswa",
    },
    siswa2: {
      username: process.env.E2E_SISWA_2_USERNAME || "siswa2_e2e",
      password: process.env.E2E_SISWA_2_PASSWORD || "Password123!",
      dashboardPath: "/dashboard/siswa",
    },
  },
};
