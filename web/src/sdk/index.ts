import axios, { AxiosInstance } from "axios";
import {
  LoginRequest,
  LoginResponse,
  SetupRequest,
  SetupResponse,
  StatusResponse,
  Response,
  CreateTicketRequest,
  CreateTicketResponse,
  TicketsResponse,
  TicketResponse,
  CreateCommentResponse,
  CreateCommentRequest,
  CommentsResponse,
  LabelsResponse,
  CreateLabelRequest,
  CreateLabelResponse,
  PatchTicketRequest,
} from "./types.gen";

type ErrorResponse = {
  message: string;
  errors?: { field: string; validator: string }[];
};

type SuccessResponse<T> = T extends Response<infer U> ? { data: U } : never;

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

  login = async (req: LoginRequest) => {
    return this.post<LoginResponse>("/login", req);
  };

  async status() {
    return this.client
      .get<SuccessResponse<StatusResponse>>("/status")
      .then((res) => res.data);
  }

  createTicket = async (req: CreateTicketRequest) => {
    return this.post<CreateTicketResponse>("/tickets", req);
  };

  tickets = async () => {
    return this.client
      .get<SuccessResponse<TicketsResponse>>("/tickets")
      .then((res) => res.data);
  };

  ticket = async (id: number) => {
    return this.client
      .get<SuccessResponse<TicketResponse>>(`/tickets/${id}`)
      .then((res) => res.data.data);
  };

  patchTicket = async (id: number, req: PatchTicketRequest) => {
    return this.client
      .patch<SuccessResponse<TicketResponse>>(`/tickets/${id}`, req)
      .then((res) => res.data.data);
  };

  createComment = async (ticketId: number, req: CreateCommentRequest) => {
    return this.client
      .post<CreateCommentResponse>(`/tickets/${ticketId}/comments`, req)
      .then((res) => res.data.data);
  };

  comments = async (ticketId: number) => {
    return this.client
      .get<SuccessResponse<CommentsResponse>>(`/tickets/${ticketId}/comments`)
      .then((res) => res.data.data);
  };

  labels = async (name?: string) => {
    const querySearch = new URLSearchParams();
    if (name) {
      querySearch.append("name", name);
    }
    return this.client
      .get<SuccessResponse<LabelsResponse>>(`/labels?${querySearch.toString()}`)
      .then((res) => res.data.data);
  };

  createLabel = async (req: CreateLabelRequest) => {
    return this.client
      .post<CreateLabelResponse>(`/labels`, req)
      .then((res) => res.data.data);
  };

  static isErrorResponse(res: unknown): res is ErrorResponse {
    return typeof res === "object" && res !== null && "message" in res;
  }

  private async post<TResponse, TRequest = unknown>(
    path: string,
    req: TRequest
  ) {
    return this.client
      .post<TResponse>(path, req)
      .then((res) => res.data as SuccessResponse<TResponse>)
      .catch((err) => {
        if (axios.isAxiosError(err) && err.response?.data) {
          return err.response.data as ErrorResponse;
        }
        throw err;
      });
  }
}

export function isSuccess<T>(res: ErrorResponse | T): res is T {
  return !(res as ErrorResponse).message;
}
