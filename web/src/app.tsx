import { RouterProvider, createBrowserRouter } from "react-router-dom";
import { SetupPage } from "./pages/setup";
import { Toaster } from "./ui/toaster";
import { QueryClient, QueryClientProvider } from "@tanstack/react-query";
import { StatusProvider } from "./status";
import { routes } from "./routes";

const queryClient = new QueryClient();
const router = createBrowserRouter(routes);

export function App() {
  return (
    <QueryClientProvider client={queryClient}>
      <StatusProvider fallback={<SetupPage />}>
        <RouterProvider router={router} />
      </StatusProvider>
      <Toaster />
    </QueryClientProvider>
  );
}

export default App;
