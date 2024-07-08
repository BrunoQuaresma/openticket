import { expect, test } from "vitest";
import { StatusProvider } from "./status";
import { render } from "@testing-library/react";
import { QueryClient, QueryClientProvider } from "@tanstack/react-query";

test("shows the children when setup is complete", () => {
  const { getByText } = render(
    <QueryClientProvider client={new QueryClient()}>
      <StatusProvider setup={<div>Setup component here</div>}>
        <div>Children component here</div>
      </StatusProvider>
    </QueryClientProvider>
  );

  expect(getByText("Setup complete")).toBeTruthy();
});
