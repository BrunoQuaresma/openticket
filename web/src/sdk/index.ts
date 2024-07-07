import {
  LoginRequest,
  LoginResponse,
  SetupRequest,
  SetupResponse,
} from "./types";
import axios, { AxiosInstance } from "axios";

export class OpenticketSdk {
  private client: AxiosInstance;

  constructor(token?: string) {
    this.client = axios.create({
      baseURL: "http://localhost:8080/api",
      timeout: 1_000,
      headers: { "OPENTICKET-TOKEN": token },
    });
  }

  async setup(req: SetupRequest) {
    return this.client.post<SetupResponse>("/setup", req);
  }

  async login(req: LoginRequest) {
    return this.client.post<LoginResponse>("/login", req);
  }
}
