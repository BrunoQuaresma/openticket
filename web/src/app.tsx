import { RouterProvider, createBrowserRouter } from "react-router-dom";
import { SetupPage } from "./pages/setup";
import { Toaster } from "./ui/toaster";
import { QueryClient, QueryClientProvider } from "@tanstack/react-query";
import { StatusProvider } from "./status";
import { routes } from "./routes";
import { HelmetProvider } from "react-helmet-async";

type Router = ReturnType<typeof createBrowserRouter>;

const defaultQueryClient = new QueryClient();
const defaultRouter = createBrowserRouter(routes);

type AppProps = {
  queryClient?: QueryClient;
  router?: Router;
};

export function App({
  queryClient = defaultQueryClient,
  router = defaultRouter,
}: AppProps) {
  return (
    <HelmetProvider>
      <QueryClientProvider client={queryClient}>
        <StatusProvider fallback={<SetupPage />}>
          <RouterProvider router={router} />
        </StatusProvider>
        <Toaster />
      </QueryClientProvider>
    </HelmetProvider>
  );
}

export default App;
