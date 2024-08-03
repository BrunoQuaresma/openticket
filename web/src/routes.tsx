import { RouteObject } from "react-router-dom";
import { Dashboard } from "./pages/dashboard/layout";
import { LoginPage } from "./pages/login";
import { IndexDashboardPage } from "./pages/dashboard";

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
        element: <IndexDashboardPage />,
      },
    ],
  },
  {
    path: "*",
    element: <div>404 page</div>,
  },
];
