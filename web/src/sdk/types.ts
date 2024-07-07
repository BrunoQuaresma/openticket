export type SetupRequest = {
  name: string;
  username: string;
  email: string;
  password: string;
};

export type SetupResponse = {
  id: string;
};

export type LoginRequest = {
  email: string;
  password: string;
};

export type LoginResponse = {
  token: string;
};
