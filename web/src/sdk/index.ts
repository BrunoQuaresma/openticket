import {
  LoginRequest,
  LoginResponse,
  SetupRequest,
  SetupResponse,
  StatusResponse,
} from "./types";
import axios, { AxiosInstance } from "axios";

type ErrorResponse = {
  message: string;
  errors?: { field: string; validator: string }[];
};

export class OpenticketSdk {
  private client: AxiosInstance;

  constructor(token?: string) {
    this.client = axios.create({
      baseURL: "/api",
      timeout: 5_000,
      headers: { "OPENTICKET-TOKEN": token },
    });
  }

  async setup(req: SetupRequest) {
    return this.post<SetupResponse>("/setup", req);
  }

  async login(req: LoginRequest) {
    return this.post<LoginResponse>("/login", req);
  }

  async status() {
    return this.client.get<StatusResponse>("/status").then((res) => res.data);
  }

  static isErrorResponse(res: unknown): res is ErrorResponse {
    return typeof res === "object" && res !== null && "message" in res;
  }

  private async post<TResponse, TRequest = unknown>(
    path: string,
    req: TRequest
  ) {
    return this.client
      .post<TResponse>(path, req)
      .then((res) => res.data)
      .catch((err) => {
        if (axios.isAxiosError(err) && err.response?.data) {
          return err.response.data as ErrorResponse;
        }
        throw err;
      });
  }
}
