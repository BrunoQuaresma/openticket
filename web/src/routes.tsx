import { RouteObject } from "react-router-dom";
import { Dashboard } from "./layouts/dashboard";
import { LoginPage } from "./pages/login";

export const routes: RouteObject[] = [
  {
    path: "/login",
    element: <LoginPage />,
  },
  {
    path: "/",
    element: <Dashboard />,
    children: [
      {
        index: true,
        element: <div>app page</div>,
      },
    ],
  },
  {
    path: "*",
    element: <div>404 page</div>,
  },
];
