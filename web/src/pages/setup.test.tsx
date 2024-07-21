import { createMemoryRouter } from "react-router-dom";
import { server } from "@/test-utils";
import { http, HttpResponse } from "msw";
import { SetupResponse, StatusResponse } from "@/sdk/types.gen";
import { render, screen, waitFor } from "@testing-library/react";
import userEvent from "@testing-library/user-event";
import { expect, test } from "vitest";
import App from "@/app";
import { routes } from "@/routes";
import { QueryClient } from "@tanstack/react-query";

test("sucessful setup", async () => {
  const user = userEvent.setup();
  const router = createMemoryRouter(routes, { initialEntries: ["/"] });

  server.use(
    http.get("/api/status", () => {
      return HttpResponse.json<StatusResponse>({
        data: { setup: false, user: undefined },
      });
    })
  );
  server.use(
    http.post("/api/setup", () => {
      return HttpResponse.json<SetupResponse>({
        data: { id: 1, role: "admin" },
      });
    })
  );

  const { container } = render(
    <App router={router} queryClient={new QueryClient()} />
  );

  const nameField = await screen.findByLabelText("Name");
  await user.type(nameField, "User");
  await user.type(screen.getByLabelText(/username/i), "user");
  await user.type(screen.getByLabelText(/email/i), "user@openticket.com");
  await user.type(screen.getByLabelText("Password"), "s3cur3p@ssw0rd");
  await user.type(screen.getByLabelText(/confirm password/i), "s3cur3p@ssw0rd");
  await user.click(screen.getByRole("button", { name: /setup/i }));
  await waitFor(() => {
    expect(container).toHaveTextContent(/login/i);
  });
});
