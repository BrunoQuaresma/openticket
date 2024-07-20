import { expect, test } from "vitest";
import { StatusProvider } from "./status";
import { render, waitFor } from "@testing-library/react";
import { QueryClient, QueryClientProvider } from "@tanstack/react-query";
import { http, HttpResponse } from "msw";
import { server } from "./test-utils";
import { StatusResponse } from "./sdk/types.gen";

const user = {
  id: 1,
  name: "Test User",
  username: "testuser",
  email: "",
  role: "admin",
};

test("shows setup page when setup is not complete", async () => {
  server.use(
    http.get("/api/status", () => {
      return HttpResponse.json<StatusResponse>({
        data: {
          setup: false,
        },
      });
    })
  );

  const { getByText } = render(
    <QueryClientProvider client={new QueryClient()}>
      <StatusProvider fallback={<div>setup component</div>}>
        <div>children</div>
      </StatusProvider>
    </QueryClientProvider>
  );

  await waitFor(() => {
    expect(getByText("setup component")).toBeTruthy();
  });
});

test("shows the children when setup is complete", async () => {
  server.use(
    http.get("/api/status", () => {
      return HttpResponse.json<StatusResponse>({
        data: {
          setup: true,
          user,
        },
      });
    })
  );

  const { getByText } = render(
    <QueryClientProvider client={new QueryClient()}>
      <StatusProvider fallback={<div>setup component</div>}>
        <div>children</div>
      </StatusProvider>
    </QueryClientProvider>
  );

  await waitFor(() => {
    expect(getByText("children")).toBeTruthy();
  });
});
