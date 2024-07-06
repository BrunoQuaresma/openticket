import { RouterProvider, createBrowserRouter } from "react-router-dom";
import { SetupPage } from "./setup";
import { LoginPage } from "./login";

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
  return <RouterProvider router={router} />;
}

export default App;
