import { afterAll, beforeAll, afterEach } from "vitest";
import { server } from "./test-utils";
import "@testing-library/jest-dom/vitest";

beforeAll(() => {
  server.listen({
    onUnhandledRequest: "error",
  });
});

afterEach(() => {
  server.resetHandlers();
});

afterAll(() => {
  server.close();
});
