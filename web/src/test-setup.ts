import { afterAll, beforeAll, afterEach } from "vitest";
import { server } from "./test-utils";

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
