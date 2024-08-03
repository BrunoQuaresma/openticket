import { Navigate, RouteObject } from "react-router-dom";
import { Dashboard } from "./pages/dashboard/layout";
import { LoginPage } from "./pages/login";
import { IndexDashboardPage } from "./pages/dashboard";
import { TicketPage } from "./pages/dashboard/tickets/ticket";

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
      {
        path: "tickets",
        children: [
          { index: true, element: <Navigate to="/" replace /> },
          {
            path: ":ticketId",
            element: <TicketPage />,
          },
        ],
      },
    ],
  },
  {
    path: "*",
    element: <div>404 page</div>,
  },
];
