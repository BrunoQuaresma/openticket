import { createMemoryRouter } from "react-router-dom";
import { server } from "@/test-utils";
import { http, HttpResponse } from "msw";
import { LoginResponse, Response, StatusResponse } from "@/sdk/types.gen";
import { render, screen, waitFor } from "@testing-library/react";
import userEvent from "@testing-library/user-event";
import { expect, test } from "vitest";
import App from "@/app";
import { routes } from "@/routes";
import { QueryClient } from "@tanstack/react-query";

const userData = {
  id: 1,
  name: "Test User",
  username: "testuser",
  email: "",
  role: "admin",
};

test("goes to the app page when the user logs in", async () => {
  const user = userEvent.setup();
  const router = createMemoryRouter(routes, { initialEntries: ["/login"] });

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
        data: {
          token: "secure-token",
          user: userData,
        },
      });
    })
  );

  render(<App router={router} queryClient={new QueryClient()} />);

  const emailField = await screen.findByLabelText(/email/i);
  await user.type(emailField, "user@openticket.com");
  await user.type(screen.getByLabelText(/password/i), "s3cur3p@ssw0rd");
  await user.click(screen.getByRole("button", { name: /login/i }));
  await waitFor(() => {
    expect(router.state.location.pathname).toEqual("/");
  });
});

test("redirects to the app page when the user is already logged in", async () => {
  const router = createMemoryRouter(routes, { initialEntries: ["/login"] });

  server.use(
    http.get("/api/status", () => {
      return HttpResponse.json<StatusResponse>({
        data: { setup: true, user: userData },
      });
    })
  );

  render(<App router={router} queryClient={new QueryClient()} />);

  await waitFor(() => {
    expect(router.state.location.pathname).toEqual("/");
  });
});

test.only("display errors from the server", async () => {
  const user = userEvent.setup();
  const router = createMemoryRouter(routes, { initialEntries: ["/login"] });
  const errorMessage = "Invalid email or password";

  server.use(
    http.get("/api/status", () => {
      return HttpResponse.json<StatusResponse>({
        data: { setup: true, user: undefined },
      });
    })
  );
  server.use(
    http.post("/api/login", () => {
      return HttpResponse.json<Response<undefined>>(
        { message: errorMessage },
        { status: 401 }
      );
    })
  );

  render(<App router={router} queryClient={new QueryClient()} />);

  const emailField = await screen.findByLabelText(/email/i);
  await user.type(emailField, "user@openticket.com");
  await user.type(screen.getByLabelText(/password/i), "s3cur3p@ssw0rd");
  await user.click(screen.getByRole("button", { name: /login/i }));
  await waitFor(() => {
    expect(screen.getByText(errorMessage)).toBeInTheDocument();
  });
});
