import fs from "node:fs/promises";
import path from "node:path";
import { env } from "../config/env.js";

function sanitizeFileName(value) {
  return value.replace(/[^a-z0-9-_]+/gi, "_").replace(/^_+|_+$/g, "").slice(0, 120);
}

export async function takeScreenshot(driver, testTitle) {
  if (!driver) return null;
  const dir = path.join(env.rootDir, "artifacts", "screenshots");
  await fs.mkdir(dir, { recursive: true });
  const fileName = `${new Date().toISOString().replace(/[:.]/g, "-")}_${sanitizeFileName(testTitle)}.png`;
  const filePath = path.join(dir, fileName);
  const image = await driver.takeScreenshot();
  await fs.writeFile(filePath, image, "base64");
  return filePath;
}
