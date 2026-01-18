export type User = {
  id?: number;
  username: string;
  role: string;
};

export type AuthStatus = "loading" | "authenticated" | "guest";

export type LoginPayload = {
  username: string;
  password: string;
};

export type AuthContextValue = {
  user: User | null;
  status: AuthStatus;
  login: (payload: LoginPayload) => Promise<void>;
  logout: () => Promise<void>;
  refetchMe: () => Promise<void>;
};
