import fs from "node:fs";
import path from "node:path";
import { fileURLToPath } from "node:url";

const __filename = fileURLToPath(import.meta.url);
const rootDir = path.resolve(path.dirname(__filename), "..");
const specsDir = path.join(rootDir, "specs");

const expected = [];
const scenarioCounts = {
  "F-01": 3, "F-02": 3, "F-03": 3, "F-04": 3, "F-05": 3, "F-06": 3,
  "F-07": 3, "F-08": 3, "F-09": 3, "F-10": 3, "F-11": 3, "F-12": 3,
  "F-13": 3, "F-14": 3, "F-15": 3, "F-16": 3, "F-17": 3, "F-18": 3,
  "F-19": 3, "F-20": 3, "F-21": 3, "F-22": 3, "F-23": 3, "F-24": 3,
  "F-25": 3, "F-26": 3, "F-27": 3, "F-28": 3, "F-29": 3, "F-30": 3,
  "F-31": 3, "F-32": 3, "F-33": 3, "F-34": 3, "F-35": 3, "F-36": 3,
  "F-37": 3, "F-38": 3, "F-39": 3, "F-40": 3, "F-41": 3, "F-42": 3,
  "F-43": 3, "F-44": 3, "F-45": 3, "F-46": 3, "F-47": 3, "F-48": 3,
};

for (const [featureId, count] of Object.entries(scenarioCounts)) {
  for (let index = 1; index <= count; index += 1) {
    expected.push(`${featureId}-S${String(index).padStart(2, "0")}`);
  }
}

const content = fs
  .readdirSync(specsDir)
  .filter((file) => file.endsWith(".spec.js"))
  .map((file) => fs.readFileSync(path.join(specsDir, file), "utf8"))
  .join("\n");

const found = new Set([...content.matchAll(/\[(F-\d{2}-S\d{2})\]/g)].map((match) => match[1]));
const missing = expected.filter((id) => !found.has(id));

if (missing.length > 0) {
  console.error(`Scenario ID belum lengkap. Missing: ${missing.join(", ")}`);
  process.exit(1);
}

console.log(`OK: ${expected.length} scenario ID ditemukan di specs.`);
