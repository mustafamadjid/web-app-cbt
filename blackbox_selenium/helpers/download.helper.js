import fs from "node:fs/promises";
import path from "node:path";
import { expect } from "chai";
import { downloadsDir } from "../config/browser.js";
import { env } from "../config/env.js";
import { sleep } from "./wait.helper.js";

export async function waitForDownloadedFile(extension, timeout = env.timeout) {
  const startedAt = Date.now();
  while (Date.now() - startedAt < timeout) {
    const files = await fs.readdir(downloadsDir());
    const match = files.find((file) => file.endsWith(extension) && !file.endsWith(".crdownload"));
    if (match) {
      const filePath = path.join(downloadsDir(), match);
      const stat = await fs.stat(filePath);
      if (stat.size > 0) return filePath;
    }
    await sleep(500);
  }
  throw new Error(`File download ${extension} tidak muncul di ${downloadsDir()}.`);
}

export async function expectDownloadedFileNotEmpty(filePath) {
  const stat = await fs.stat(filePath);
  expect(stat.size).to.be.greaterThan(0);
}
