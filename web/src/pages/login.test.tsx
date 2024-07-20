import { createMemoryRouter, RouterProvider } from "react-router-dom";
import { LoginPage } from "./login";
import { server } from "@/test-utils";
import { http, HttpResponse } from "msw";
import { LoginResponse, StatusResponse } from "@/sdk/types.gen";
import { QueryClientProvider, QueryClient } from "@tanstack/react-query";
import { render, screen, waitFor } from "@testing-library/react";
import userEvent from "@testing-library/user-event";
import { expect, test } from "vitest";
import { StatusProvider } from "@/status";

test("goes to the app page when the user logs in", async () => {
  server.use(
    http.get("/api/status", () => {
      return HttpResponse.json<StatusResponse>({
        data: { setup: true, user: undefined },
      });
    })
  );
  server.use(
    http.post("/api/login", () => {
      return HttpResponse.json<LoginResponse>({
        data: { session_token: "12345678" },
      });
    })
  );
  const router = createMemoryRouter(
    [
      {
        path: "/login",
        element: <LoginPage />,
      },
      {
        path: "/",
        element: <>app page</>,
      },
    ],
    {
      initialEntries: ["/login"],
    }
  );
  render(
    <QueryClientProvider client={new QueryClient()}>
      <StatusProvider fallback={<>setup page</>}>
        <RouterProvider router={router} />
      </StatusProvider>
    </QueryClientProvider>
  );
  const user = userEvent.setup();

  const emailField = await screen.findByLabelText(/email/i);
  await user.type(emailField, "user@openticket.com");
  await user.type(screen.getByLabelText(/password/i), "s3cur3p@ssw0rd");
  await user.click(screen.getByRole("button", { name: /login/i }));

  await waitFor(() => {
    expect(router.state.location.pathname).toEqual("/");
  });
});

const userData = {
  id: 1,
  name: "Test User",
  username: "testuser",
  email: "",
  role: "admin",
};

test("redirects to the app page when the user is already logged in", async () => {
  server.use(
    http.get("/api/status", () => {
      return HttpResponse.json<StatusResponse>({
        data: { setup: true, user: userData },
      });
    })
  );
  const router = createMemoryRouter(
    [
      {
        path: "/login",
        element: <LoginPage />,
      },
      {
        path: "/",
        element: <>app page</>,
      },
    ],
    {
      initialEntries: ["/login"],
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
    expect(router.state.location.pathname).toEqual("/");
  });
});
