import { RouterProvider, createBrowserRouter } from "react-router-dom";
import { SetupPage } from "./setup";

const router = createBrowserRouter([
  {
    path: "/",
    element: <SetupPage />,
  },
]);

export function App() {
  return <RouterProvider router={router} />;
}

export default App;
