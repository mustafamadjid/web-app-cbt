import { env } from "../config/env.js";

export const accounts = env.accounts;

export function accountFor(role) {
  const account = accounts[role];
  if (!account) {
    throw new Error(`Role akun tidak dikenal: ${role}`);
  }
  return account;
}
