import { expect, test } from "vitest";
import { server } from "../test-utils";
import { http, HttpResponse } from "msw";
import { createMemoryRouter, RouterProvider } from "react-router-dom";
import { Dashboard } from "./dashboard";
import { render, waitFor } from "@testing-library/react";
import { StatusProvider } from "../status";
import { QueryClient, QueryClientProvider } from "@tanstack/react-query";
import { StatusResponse } from "@/sdk/types.gen";

const user = {
  id: 1,
  name: "Test User",
  username: "testuser",
  email: "",
  role: "admin",
};

test("redirects to /login if not authenticated", async () => {
  server.use(
    http.get("/api/status", () => {
      return HttpResponse.json<StatusResponse>({
        data: { setup: true, user: undefined },
      });
    })
  );
  const router = createMemoryRouter(
    [
      {
        path: "/login",
        element: <>login page</>,
      },
      {
        path: "/app",
        element: <Dashboard />,
        children: [
          {
            index: true,
            element: <>app page</>,
          },
        ],
      },
    ],
    {
      initialEntries: ["/app"],
    }
  );
  render(
    <QueryClientProvider client={new QueryClient()}>
      <StatusProvider fallback={<>setup page</>}>
        <RouterProvider router={router} />
      </StatusProvider>
    </QueryClientProvider>
  );
  await waitFor(() => {
    expect(router.state.location.pathname).toEqual("/login");
  });
});

test("stay on the page if authenticated", async () => {
  server.use(
    http.get("/api/status", () => {
      return HttpResponse.json<StatusResponse>({
        data: { setup: true, user },
      });
    })
  );
  const router = createMemoryRouter(
    [
      {
        path: "/login",
        element: <>login page</>,
      },
      {
        path: "/app",
        element: <Dashboard />,
        children: [
          {
            index: true,
            element: <>app page</>,
          },
        ],
      },
    ],
    {
      initialEntries: ["/app"],
    }
  );
  render(
    <QueryClientProvider client={new QueryClient()}>
      <StatusProvider fallback={<>setup page</>}>
        <RouterProvider router={router} />
      </StatusProvider>
    </QueryClientProvider>
  );
  await waitFor(() => {
    expect(router.state.location.pathname).toEqual("/app");
  });
});
