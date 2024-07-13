import { RouterProvider, createBrowserRouter } from "react-router-dom";
import { SetupPage } from "./pages/setup";
import { LoginPage } from "./pages/login";
import { Toaster } from "./ui/toaster";
import { QueryClient, QueryClientProvider } from "@tanstack/react-query";
import { StatusProvider } from "./status";

const queryClient = new QueryClient();

const router = createBrowserRouter([
  {
    path: "/login",
    element: <LoginPage />,
  },
]);

export function App() {
  return (
    <QueryClientProvider client={queryClient}>
      <StatusProvider setup={<SetupPage />}>
        <RouterProvider router={router} />
      </StatusProvider>
      <Toaster />
    </QueryClientProvider>
  );
}

export default App;
