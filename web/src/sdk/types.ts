type Response<T> = {
  data: T;
};

export type SetupRequest = {
  name: string;
  username: string;
  email: string;
  password: string;
};

export type SetupResponse = Response<{
  id: string;
}>;

export type LoginRequest = {
  email: string;
  password: string;
};

export type LoginResponse = Response<{
  token: string;
}>;

export type HealthResponse = Response<{
  setup: true;
  health: true;
}>;
