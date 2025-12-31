// // hooks/useUser.ts
// import { useEffect, useState } from "react";
// import { ApiError } from "../services/api";
// import {
//   getUserById,
//   updateUser,
//   type User,
//   type UpdateUserPayload,
// } from "../services/users.service";

// export function useUser(userId: string, token?: string | null) {
//   const [user, setUser] = useState<User | null>(null);
//   const [loading, setLoading] = useState(false);
//   const [error, setError] = useState<string | null>(null);

//   useEffect(() => {
//     let mounted = true;
//     (async () => {
//       setLoading(true);
//       setError(null);
//       try {
//         const u = await getUserById(userId, token);
//         if (mounted) setUser(u);
//       } catch (e) {
//         if (!mounted) return;
//         if (e instanceof ApiError) setError(e.message);
//         else setError("Unknown error");
//       } finally {
//         if (mounted) setLoading(false);
//       }
//     })();
//     return () => {
//       mounted = false;
//     };
//   }, [userId, token]);

//   async function save(payload: UpdateUserPayload) {
//     const updated = await updateUser(userId, payload, token);
//     setUser(updated);
//     return updated;
//   }

//   return { user, loading, error, save };
// }
