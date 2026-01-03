/** ===== Access token store (in-memory) ===== */
let accessToken: string | null = null;

export const authToken = {
  get: () => accessToken,
  set: (t: string | null) => {
    accessToken = t;
  },
  clear: () => {
    accessToken = null;
  },
};
