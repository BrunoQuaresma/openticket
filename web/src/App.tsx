import { RouterProvider, createBrowserRouter } from "react-router-dom";
import { SetupPage } from "./setup";
import { LoginPage } from "./login";
import { Toaster } from "./ui/toaster";

const router = createBrowserRouter([
  {
    path: "/",
    element: <SetupPage />,
  },
  {
    path: "/login",
    element: <LoginPage />,
  },
]);

export function App() {
  return (
    <>
      <RouterProvider router={router} />
      <Toaster />
    </>
  );
}

export default App;
