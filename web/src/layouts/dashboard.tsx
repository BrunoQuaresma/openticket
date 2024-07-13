import { Navigate, Outlet } from "react-router-dom";
import { useStatus } from "../status";

export function Dashboard() {
  const { data } = useStatus();

  if (!data.user) {
    return <Navigate to="/login" replace />;
  }

  return (
    <>
      <h1>Dashboard</h1>
      <Outlet />
    </>
  );
}
