import { createMemoryRouter, RouterProvider } from "react-router-dom";
import { expect, test } from "vitest";
import { routes } from "./routes";
import { render, screen } from "@testing-library/react";

test("display not found page", () => {
  const router = createMemoryRouter(routes, {
    initialEntries: ["/some-wrong-page-path"],
  });
  render(<RouterProvider router={router} />);
  expect(screen.getByText("404", { exact: false })).toBeTruthy();
});
